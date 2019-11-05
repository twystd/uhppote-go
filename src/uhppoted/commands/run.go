package commands

import (
	"errors"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"math"
	"net"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"
	"uhppote"
	"uhppote/types"
	"uhppoted/config"
	"uhppoted/rest"
)

type status struct {
	touched time.Time
	status  types.Status
}

type alerts struct {
	missing      bool
	touched      bool
	synchronized bool
}

type state struct {
	started time.Time

	healthcheck struct {
		touched *time.Time
		alerted bool
	}

	devices struct {
		status sync.Map
		errors sync.Map
	}
}

const (
	IDLE  = time.Duration(60 * time.Second)
	DELTA = 60
	DELAY = 30
)

func (c *Run) Parse(args []string) error {
	flagset := c.FlagSet()
	if flagset == nil {
		panic(fmt.Sprintf("'run' command implementation without a flagset: %#v", c))
	}

	return flagset.Parse(args)
}

func (c *Run) Description() string {
	return "Runs the uhppoted daemon/service until terminated by the system service manager"
}

func (c *Run) Usage() string {
	return "uhppoted [--debug] [--config <file>] [--logfile <file>] [--logfilesize <bytes>] [--pid <file>]"
}

func (c *Run) Help() {
	fmt.Println()
	fmt.Println("  Usage: uhppoted <options>")
	fmt.Println()
	fmt.Println("  Options:")
	fmt.Println()
	runCmd.FlagSet().VisitAll(func(f *flag.Flag) {
		fmt.Printf("    --%-12s %s\n", f.Name, f.Usage)
	})
	fmt.Println()
}

func (r *Run) execute(ctx Context, f func(*config.Config) error) error {
	conf := config.NewConfig()
	if err := conf.Load(r.configuration); err != nil {
		log.Printf("\n   WARN:  Could not load configuration (%v)\n\n", err)
	}

	if err := os.MkdirAll(r.dir, os.ModeDir|os.ModePerm); err != nil {
		return fmt.Errorf("Unable to create working directory '%v': %v", r.dir, err)
	}

	pid := fmt.Sprintf("%d\n", os.Getpid())

	if err := ioutil.WriteFile(r.pidFile, []byte(pid), 0644); err != nil {
		return fmt.Errorf("Unable to create pid file: %v\n", err)
	}

	defer func() {
		os.Remove(r.pidFile)
	}()

	return f(conf)
}

func (r *Run) run(c *config.Config, logger *log.Logger) {
	logger.Printf("START")

	// ... syscall SIG handlers

	interrupt := make(chan os.Signal, 1)

	signal.Notify(interrupt, syscall.SIGINT, syscall.SIGTERM)

	// ... listen forever

	for {
		err := r.listen(c, logger, interrupt)
		if err != nil {
			log.Printf("ERROR: %v", err)
			continue
		}

		log.Printf("exit\n")
		break
	}

	logger.Printf("STOP")
}

func (r *Run) listen(c *config.Config, logger *log.Logger, interrupt chan os.Signal) error {
	s := state{
		started: time.Now(),

		healthcheck: struct {
			touched *time.Time
			alerted bool
		}{
			touched: nil,
			alerted: false,
		},
		devices: struct {
			status sync.Map
			errors sync.Map
		}{
			status: sync.Map{},
			errors: sync.Map{},
		},
	}

	u := uhppote.UHPPOTE{
		BindAddress:      c.BindAddress,
		BroadcastAddress: c.BroadcastAddress,
		Devices:          make(map[uint32]*net.UDPAddr),
		Debug:            r.debug,
	}

	for id, d := range c.Devices {
		if d.Address != nil {
			u.Devices[id] = d.Address
		}
	}

	// ... REST task

	restd := rest.RestD{
		HttpEnabled:        c.REST.HttpEnabled,
		HttpPort:           c.REST.HttpPort,
		HttpsEnabled:       c.REST.HttpsEnabled,
		HttpsPort:          c.REST.HttpsPort,
		TLSKeyFile:         c.REST.TLSKeyFile,
		TLSCertificateFile: c.REST.TLSCertificateFile,
		CACertificateFile:  c.REST.CACertificateFile,
		CORSEnabled:        c.REST.CORSEnabled,
		OpenApi: rest.OpenApi{
			Enabled:   c.OpenApi.Enabled,
			Directory: c.OpenApi.Directory,
		},
	}

	go func() {
		restd.Run(&u, logger)
	}()

	defer rest.Close()

	// ... health-check task

	k := time.NewTicker(15 * time.Second)

	defer k.Stop()

	go func() {
		for {
			<-k.C
			healthcheck(&u, &s, logger)
		}
	}()

	// ... wait until interrupted/closed

	closed := make(chan struct{})
	w := time.NewTicker(5 * time.Second)

	defer w.Stop()

	for {
		select {
		case <-w.C:
			if err := watchdog(&u, &s, logger); err != nil {
				return err
			}

		case <-interrupt:
			logger.Printf("... interrupt")
			return nil

		case <-closed:
			logger.Printf("... closed")
			return errors.New("Server error")
		}
	}

	logger.Printf("... exit")
	return nil
}

func healthcheck(u *uhppote.UHPPOTE, st *state, l *log.Logger) {
	l.Printf("health-check")

	now := time.Now()

	for id, _ := range u.Devices {
		s, err := u.GetStatus(id)
		if err == nil {
			st.devices.status.Store(id, status{
				touched: now,
				status:  *s,
			})
		}
	}

	st.healthcheck.touched = &now
}

func watchdog(u *uhppote.UHPPOTE, st *state, l *log.Logger) error {
	ok := true
	healthCheckRunning := false
	now := time.Now()
	seconds, _ := time.ParseDuration("1s")

	// ... verify health-check

	dt := time.Since(st.started).Round(seconds)
	if st.healthcheck.touched != nil {
		dt = time.Since(*st.healthcheck.touched)
		if int64(math.Abs(dt.Seconds())) < DELAY {
			healthCheckRunning = true
		}
	}

	if int64(math.Abs(dt.Seconds())) > DELAY {
		ok = false
		if !st.healthcheck.alerted {
			l.Printf("ERROR 'health-check' subsystem has not run since %s (%s)", types.DateTime(st.started), dt)
			st.healthcheck.alerted = true
		}
	} else {
		if st.healthcheck.alerted {
			l.Printf("INFO  'health-check' subsystem is running")
			st.healthcheck.alerted = false
		}
	}

	// verify devices

	if healthCheckRunning {
		for id, _ := range u.Devices {
			alerted := alerts{
				missing:      false,
				touched:      false,
				synchronized: false,
			}

			if v, found := st.devices.errors.Load(id); found {
				alerted.missing = v.(alerts).missing
				alerted.touched = v.(alerts).touched
				alerted.synchronized = v.(alerts).synchronized
			}

			if _, found := st.devices.status.Load(id); !found {
				ok = false
				if !alerted.missing {
					l.Printf("ERROR UTC0311-L0x %s missing", types.SerialNumber(id))
					alerted.missing = true
				}
			}

			if v, found := st.devices.status.Load(id); found {
				touched := v.(status).touched
				t := time.Time(v.(status).status.SystemDateTime)
				dt := time.Since(t).Round(seconds)
				dtt := int64(math.Abs(time.Since(touched).Seconds()))

				if alerted.missing {
					l.Printf("ERROR UTC0311-L0x %s present", types.SerialNumber(id))
					alerted.missing = false
				}

				if now.After(touched.Add(IDLE)) {
					ok = false
					if !alerted.touched {
						l.Printf("ERROR UTC0311-L0x %s no response for %s", types.SerialNumber(id), time.Since(touched).Round(seconds))
						alerted.touched = true
					}
				} else {
					if alerted.touched {
						l.Printf("INFO  UTC0311-L0x %s connected", types.SerialNumber(id))
						alerted.touched = false
					}
				}

				if dtt < DELTA/2 && int64(math.Abs(dt.Seconds())) > DELTA {
					ok = false
					if !alerted.synchronized {
						l.Printf("ERROR UTC0311-L0x %s system time not synchronized: %s (%s)", types.SerialNumber(id), types.DateTime(t), dt)
						alerted.synchronized = true
					}
				} else {
					if alerted.synchronized {
						l.Printf("INFO   UTC0311-L0x %s system time synchronized: %s (%s)", types.SerialNumber(id), types.DateTime(t), dt)
						alerted.synchronized = false
					}
				}
			}

			st.devices.errors.Store(id, alerted)
		}
	}

	if ok {
		l.Printf("watchdog: OK")
	} else {
		l.Printf("watchdog: ERROR")
	}

	return nil
}

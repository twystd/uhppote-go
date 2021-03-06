package main

import (
	"flag"
	"fmt"
	"net"
	"os"
	"uhppote"
	"uhppote-cli/commands"
	"uhppote-cli/config"
)

type addr struct {
	address *net.UDPAddr
}

var cli = []commands.Command{
	&commands.VersionCommand{},
	&commands.GetDevicesCommand{},
	&commands.GetDeviceCommand{},
	&commands.SetAddressCommand{},
	&commands.GetStatusCommand{},
	&commands.GetTimeCommand{},
	&commands.SetTimeCommand{},
	&commands.GetDoorDelayCommand{},
	&commands.SetDoorDelayCommand{},
	&commands.GetDoorControlCommand{},
	&commands.SetDoorControlCommand{},
	&commands.GetListenerCommand{},
	&commands.SetListenerCommand{},
	&commands.GetCardsCommand{},
	&commands.GetCardCommand{},
	&commands.GrantCommand{},
	&commands.RevokeCommand{},
	&commands.RevokeAllCommand{},
	&commands.Load{},
	&commands.GetEventsCommand{},
	&commands.GetEventIndexCommand{},
	&commands.SetEventIndexCommand{},
	&commands.OpenDoorCommand{},
	&commands.ListenCommand{},
}

var options = struct {
	config    string
	bind      addr
	broadcast addr
	listen    addr
	debug     bool
}{
	config:    ".config",
	bind:      addr{nil},
	broadcast: addr{nil},
	listen:    addr{nil},
	debug:     false,
}

func main() {
	flag.StringVar(&options.config, "config", "", "Specifies the path for the config file")
	flag.Var(&options.bind, "bind", "Sets the local IP address and port to which to bind (e.g. 192.168.0.100:60001)")
	flag.Var(&options.broadcast, "broadcast", "Sets the IP address and port for UDP broadcast (e.g. 192.168.0.255:60000)")
	flag.Var(&options.listen, "listen", "Sets the local IP address and port to which to bind for events (e.g. 192.168.0.100:60001)")
	flag.BoolVar(&options.debug, "debug", false, "Displays vaguely useful information while processing a command")
	flag.Parse()

	u := uhppote.UHPPOTE{
		Devices: make(map[uint32]*uhppote.Device),
		Debug:   options.debug,
	}

	conf := config.NewConfig()
	if err := conf.Load(options.config); err != nil {
		fmt.Printf("\n   ERROR: %v\n\n", err)
		os.Exit(1)
	}

	u.BindAddress = conf.BindAddress
	u.BroadcastAddress = conf.BroadcastAddress
	u.ListenAddress = conf.ListenAddress

	for s, d := range conf.Devices {
		u.Devices[s] = uhppote.NewDevice(s, d.Address, d.Rollover)
	}

	if options.bind.address != nil {
		u.BindAddress = options.bind.address
		u.ListenAddress = options.bind.address
	}

	if options.broadcast.address != nil {
		u.BroadcastAddress = options.broadcast.address
	}

	if options.listen.address != nil {
		u.ListenAddress = options.listen.address
	}

	cmd, err := parse()
	if err != nil {
		fmt.Printf("\n   ERROR: %v\n\n", err)
		os.Exit(1)
	}

	if cmd == nil {
		help()
		return
	}

	ctx := commands.NewContext(&u, conf)

	err = cmd.Execute(ctx)
	if err != nil {
		fmt.Printf("\n   ERROR: %v\n\n", err)
		os.Exit(1)
	}
}

func parse() (commands.Command, error) {
	var cmd commands.Command = nil
	var err error = nil

	if len(os.Args) > 1 {
		for _, c := range cli {
			if c.CLI() == flag.Arg(0) {
				cmd = c
			}
		}
	}

	return cmd, err
}

func (b *addr) String() string {
	return b.address.String()
}

func (b *addr) Set(s string) error {
	address, err := net.ResolveUDPAddr("udp", s)
	if err != nil {
		return err
	}

	b.address = address

	return nil
}

func help() {
	if len(flag.Args()) > 0 && flag.Arg(0) == "help" {
		if len(flag.Args()) > 1 {

			if flag.Arg(1) == "commands" {
				helpCommands()
				return
			}

			for _, c := range cli {
				if c.CLI() == flag.Arg(1) {
					c.Help()
					return
				}
			}

			fmt.Printf("Invalid command: %v. Type 'help commands' to get a list of supported commands\n", flag.Arg(1))
			return
		}
	}

	usage()
}

func usage() {
	fmt.Println()
	fmt.Println("  Usage: uhppote-cli [options] <command>")
	fmt.Println()
	fmt.Println("  Commands:")
	fmt.Println()
	fmt.Println("    help             Displays this message")
	fmt.Println("                     For help on a specific command use 'uhppote-cli help <command>'")

	for _, c := range cli {
		fmt.Printf("    %-16s %s\n", c.CLI(), c.Description())
	}

	fmt.Println()
	fmt.Println("  Options:")
	fmt.Println()
	fmt.Println("    --config    Sets the configuration file")
	fmt.Println("    --bind      Sets the local IP address and port to use")
	fmt.Println("    --broadcast Sets the IP address and port to use for UDP broadcast")
	fmt.Println("    --listen    Sets the local IP address and port to use for receiving device events")
	fmt.Println("    --debug     Displays vaguely useful internal information")
	fmt.Println()
}

func helpCommands() {
	fmt.Println("Supported commands:")
	fmt.Println()

	for _, c := range cli {
		fmt.Printf(" %-16s %s\n", c.CLI(), c.Usage())
	}

	fmt.Println()
}

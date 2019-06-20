CLI = ./bin/uhppote-cli
SIMULATOR = ./bin/uhppote-simulator
DEBUG = --debug
LOCAL = 192.168.1.100:51234
CARD = 6154412
SERIALNO = 423187757
DOOR = 3

all: test      \
	 benchmark \
     coverage

format: 
	gofmt -w=true src/uhppote/*.go
	gofmt -w=true src/uhppote/types/*.go
	gofmt -w=true src/uhppote/encoding/bcd/*.go
	gofmt -w=true src/uhppote/encoding/UTO311-L0x/*.go
	gofmt -w=true src/uhppote-cli/*.go
	gofmt -w=true src/uhppote-cli/commands/*.go
	gofmt -w=true src/uhppote-cli/config/*.go
	gofmt -w=true src/uhppote-cli/parsers/*.go
	gofmt -w=true src/uhppote-simulator/*.go
	gofmt -w=true src/uhppote-simulator/commands/*.go
	gofmt -w=true src/uhppote-simulator/simulator/*.go
	gofmt -w=true src/uhppote-simulator/simulator/entities/*.go
	gofmt -w=true src/integration-tests/*.go

release: format
	mkdir -p dist/windows
	mkdir -p dist/macosx
	mkdir -p dist/linux
	env GOOS=windows GOARCH=amd64  go build uhppote-cli;       mv uhppote-cli.exe dist/windows
	env GOOS=darwin  GOARCH=amd64  go build uhppote-cli;       mv uhppote-cli dist/macosx
	env GOOS=linux   GOARCH=amd64  go build uhppote-cli;       mv uhppote-cli dist/linux
	env GOOS=windows GOARCH=amd64  go build uhppote-simulator; mv uhppote-simulator.exe dist/windows
	env GOOS=darwin  GOARCH=amd64  go build uhppote-simulator; mv uhppote-simulator dist/macosx
	env GOOS=linux   GOARCH=amd64  go build uhppote-simulator; mv uhppote-simulator dist/linux

build: format
	go install uhppote-cli
	go install uhppote-simulator

test: build
	go clean -testcache
	go test -count=1 src/uhppote/*.go
	go test -count=1 src/uhppote/encoding/bcd/*.go
	go test -count=1 src/uhppote/encoding/UTO311-L0x/*.go

integration-tests: build
	go clean -testcache
	go test -count=1 src/integration-tests/*.go

benchmark: build
	go test src/encoding/bcd/*.go -bench .

coverage: build
	go test -cover .

clean:
	go clean
	rm -rf bin

usage: build
	$(CLI)

debug: build
#	$(CLI) $(DEBUG) --bind 192.168.0.14:12345 --broadcast 192.168.0.255:60000 get-devices
#	$(CLI) $(DEBUG) --bind 0.0.0.0:12345                                      get-devices
#	$(CLI) $(DEBUG) --bind 0.0.0.0:12345 get-card $(SERIALNO) $(CARD)
#	$(CLI) $(DEBUG) get-devices
#	$(CLI) $(DEBUG) grant 1234567890 $(CARD) 2019-01-01 2019-12-31 1
	$(CLI) $(DEBUG) --config '.simulation' get-cards 305419896 
	$(CLI) $(DEBUG) --config '.simulation' get-card  305419896 65537
#	$(CLI) $(DEBUG) get-card  1234567890 $(CARD)
#	$(CLI) $(DEBUG) get-card  1234567890 9154419
#	$(CLI) --bind $(LOCAL) --broadcast "192.168.1.255:60000" $(DEBUG) get-card  $(SERIALNO) 9154419
#	go test -count=1 src/uhppote/encoding/UTO311-L0x/*.go -run TestUnmarshalInterface

help: build
	$(CLI)       help
	$(CLI)       help get-devices
	$(SIMULATOR) help
	$(SIMULATOR) help new-device

version: build
	$(CLI)       version
	$(SIMULATOR) version

run: build
	$(CLI) --bind $(LOCAL) $(DEBUG) get-devices
#	$(CLI) --bind $(LOCAL) $(DEBUG) set-address    $(SERIALNO) '192.168.1.125' '255.255.255.0' '0.0.0.0'
	$(CLI) --bind $(LOCAL) $(DEBUG) get-cards      $(SERIALNO)
	$(CLI) --bind $(LOCAL) $(DEBUG) get-door-delay $(SERIALNO) $(DOOR)
	$(CLI) --bind $(LOCAL) $(DEBUG) set-time       $(SERIALNO)
	$(CLI) --bind $(LOCAL) $(DEBUG) revoke         $(SERIALNO) $(CARD)

get-devices: build
#	$(CLI) --bind 0.0.0.0:0 $(DEBUG) get-devices
#	$(CLI) --bind $(LOCAL) --broadcast "192.168.1.255:60000" $(DEBUG) get-devices
	$(CLI) --bind $(LOCAL) --broadcast "255.255.255.255:60000" $(DEBUG) get-devices

set-address: build
	$(CLI) -bind $(LOCAL) $(DEBUG) set-address $(SERIALNO) '192.168.1.125' '255.255.255.0' '0.0.0.0'

get-status: build
	$(CLI) --bind $(LOCAL) $(DEBUG) get-status $(SERIALNO)

get-time: build
	$(CLI) --bind $(LOCAL) $(DEBUG) get-time $(SERIALNO)

set-time: build
	# $(CLI) -debug set-time 423187757 '2019-01-08 12:34:56'
	$(CLI) --bind $(LOCAL) $(DEBUG) set-time $(SERIALNO)

get-door-delay: build
	$(CLI) --bind $(LOCAL) $(DEBUG) get-door-delay $(SERIALNO) $(DOOR)

set-door-delay: build
	$(CLI) --bind $(LOCAL) $(DEBUG) set-door-delay $(SERIALNO) $(DOOR) 5

get-listener: build
	$(CLI) --bind  $(DEBUG) get-listener $(SERIALNO)

set-listener: build
	$(CLI) --bind $(LOCAL) $(DEBUG) set-listener $(SERIALNO) 192.168.1.100:40000

get-cards: build
	$(CLI) --config ".local" $(DEBUG) get-cards $(SERIALNO)

get-card: build
	$(CLI) --bind $(LOCAL) $(DEBUG) get-card $(SERIALNO) $(CARD)

grant: build
	$(CLI) --config ".local" $(DEBUG) grant $(SERIALNO) $(CARD) 2019-01-01 2019-12-31 1

revoke: build
	$(CLI) --config ".simulation" $(DEBUG) revoke 305419896 65537

revoke-all: build
	$(CLI) --bind $(LOCAL) $(DEBUG) revoke-all $(SERIALNO)

load: build
	$(CLI) --config ".simulation" $(DEBUG) load example.tsv

get-events: build
	$(CLI) --bind $(LOCAL) $(DEBUG) get-events $(SERIALNO)

get-event-index: build
	$(CLI) --bind $(LOCAL) $(DEBUG) get-event-index $(SERIALNO)

set-events-index: build
	$(CLI) --bind $(LOCAL) $(DEBUG) set-events-index $(SERIALNO) 23

open: build
	$(CLI) --bind $(LOCAL) $(DEBUG) open $(SERIALNO) 1

listen: build
	$(CLI) --bind 192.168.1.100:40000 $(DEBUG) listen 

simulator: build
	./bin/uhppote-simulator --debug --devices "./runtime/simulation/devices"

simulator-device: build
	./bin/uhppote-simulator --debug --devices "runtime/simulation/devices" new-device 666 --gzip




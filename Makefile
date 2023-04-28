GOCMD=go
GOBUILD=$(GOCMD) build
GORUN=$(GOCMD) run
BINARY_FOLDER=./build

Version := $(shell git describe --tags --dirty)
GitCommit := $(shell git rev-parse HEAD)
LDFLAGS := "-s -w -X main.Version=$(Version) -X main.GitCommit=$(GitCommit)"

.PHONY: all clean-all clean-binaries build-all build-% run-% run-tests run-tests-cover run-tests-cover-profile build-all-platforms

default:
	${MAKE} run-dumper

all:
		$(MAKE) run

clean-all:
		$(MAKE) clean-binaries
clean-binaries:
		-rm -f $(BINARY_FOLDER)/*

build-all:
		$(MAKE) build-dumper
build-%:
		$(GOBUILD) -ldflags="-s -w" -o $(BINARY_FOLDER)/$* -v ./examples/$*
		chmod +x $(BINARY_FOLDER)/$*

build-all-platforms:
	mkdir -p build/
	CGO_ENABLED=0 GOOS=darwin GOARCH=arm64 go build -a -ldflags $(LDFLAGS) -o build/dumper-darwin-m1 ./examples/dumper
	CGO_ENABLED=0 GOOS=linux go build -a -ldflags $(LDFLAGS) -o build/dumper-amd64 ./examples/dumper
	CGO_ENABLED=0 GOOS=darwin go build -a -ldflags $(LDFLAGS) -o build/dumper-darwin -a ./examples/dumper ./examples/dumper
	GOARM=7 GOARCH=arm CGO_ENABLED=0 GOOS=linux go build -a -ldflags $(LDFLAGS) -o build/dumper-arm ./examples/dumper
	GOARCH=arm64 CGO_ENABLED=0 GOOS=linux go build -a -ldflags $(LDFLAGS) -o build/dumper-arm64 ./examples/dumper
	GOOS=windows CGO_ENABLED=0 go build -a -ldflags $(LDFLAGS) -o build/dumper.exe ./examples/dumper

run-%:
		$(MAKE) build-$*
		$(BINARY_FOLDER)/$* --help

run-tests:
	$(GOCMD) test -count=1 ./... -v
run-tests-cover:
	$(GOCMD) test -count=1 ./... -v -cover
run-tests-cover-profile:
	$(GOCMD) test -count=1 ./... -v -coverprofile cover.out
	$(GOCMD) tool cover -html=cover.out
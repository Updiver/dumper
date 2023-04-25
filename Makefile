GOCMD=go
GOBUILD=$(GOCMD) build
GORUN=$(GOCMD) run
BINARY_FOLDER=./build

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

GOCMD=go
GOBUILD=$(GOCMD) build
GORUN=$(GOCMD) run
BINARY_FOLDER=./build

all:
		$(MAKE) run

clean-all:
		$(MAKE) clean-binaries
clean-binaries:
		-rm -f $(BINARY_FOLDER)/*

build-all:
		$(MAKE) build-backup
build-%:
		$(GOBUILD) -ldflags="-s -w" -o $(BINARY_FOLDER)/$* -v ./cmd/$*
		chmod +x $(BINARY_FOLDER)/$*

run-%:
		$(MAKE) build-$*
		$(BINARY_FOLDER)/$* --help

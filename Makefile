GOCMD=go
GOBUILD=$(GOCMD) build
GOTEST=$(GOCMD) test
GODEPS=$(GOCMD) get

FILENAME=docker-machine-driver-proxmoxve

BUILDDIR=./bin


all: deps test build-all

deps:
	$(GODEPS) ./...

test:
	$(GOTEST) -v ./test

clean:
	rm -Rf $(BUILDDIR)

build-all: $(BUILDDIR)/$(FILENAME).linux.amd64 \
		   $(BUILDDIR)/$(FILENAME).linux.arm64 \
		   $(BUILDDIR)/$(FILENAME).linux.arm \
		   $(BUILDDIR)/$(FILENAME).darwin.amd64 \
		   $(BUILDDIR)/$(FILENAME).windows.amd64.exe

$(BUILDDIR)/$(FILENAME).linux.amd64:
	GOOS=linux GOARCH=amd64 $(GOBUILD) -o $(BUILDDIR)/$(FILENAME).linux.amd64 ./cmd/$(FILENAME)

$(BUILDDIR)/$(FILENAME).linux.arm64:
	GOOS=linux GOARCH=arm64 $(GOBUILD) -o $(BUILDDIR)/$(FILENAME).linux.arm64 ./cmd/$(FILENAME)

$(BUILDDIR)/$(FILENAME).linux.arm:
	GOOS=linux GOARCH=arm $(GOBUILD) -o $(BUILDDIR)/$(FILENAME).linux.arm ./cmd/$(FILENAME)

$(BUILDDIR)/$(FILENAME).darwin.amd64:
	GOOS=darwin GOARCH=amd64 $(GOBUILD) -o $(BUILDDIR)/$(FILENAME).darwin.amd64 ./cmd/$(FILENAME)

$(BUILDDIR)/$(FILENAME).windows.amd64.exe:
	GOOS=windows GOARCH=amd64 $(GOBUILD) -o $(BUILDDIR)/$(FILENAME).windows.amd64.exe ./cmd/$(FILENAME)

.PHONY: all deps test build-all clean
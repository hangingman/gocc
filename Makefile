# Go parameters
GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get

BIN=gocc

.PHONY: all

all: dep test build

build:
	$(GOBUILD) -o ${BIN} -v

test:
	$(GOTEST) -v ./...

clean:
	$(GOCLEAN)

fmt:
	for go_file in `find . -name \*.go`; do \
		go fmt $${go_file}; \
	done

dep:
	$(GOGET) github.com/stretchr/testify
	$(GOGET) github.com/comail/colog

emacs:
	$(GOGET) github.com/rogpeppe/godef
	$(GOGET) -u github.com/nsf/gocode
	$(GOGET) github.com/golang/lint/golint
	$(GOGET) github.com/kisielk/errcheck
	$(GOGET) -u github.com/derekparker/delve/cmd/dlv

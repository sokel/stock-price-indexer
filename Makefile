#!/usr/bin/env make

.PHONY: fmt vet test version

# Version of the entire package. Do not forget to update this when it's time
# to bump the version.
VERSION = v0.0.1

# Build tag. Useful to distinguish between same-version builds, but from
# different commits.
BUILD = $(shell git rev-parse --short HEAD)

# Full version includes both semantic version and git ref if present.
ifeq (${BUILD},)
	FULL_VERSION = $(VERSION)
else
	FULL_VERSION = $(VERSION)-$(BUILD)
endif

# NOTE: variables defined with := in GNU make are expanded when they are
# defined rather than when they are used.
GOCMD := ./cmd

# NOTE: variables defined with ?= sets the default value, which can be
# overriden using env.
GO ?= go

TARGETDIR := target
INSTALLDIR := ${GOPATH}/bin/

HOSTOS := $(shell uname -s | tr '[:upper:]' '[:lower:]')
HOSTARCH := $(shell uname -m)

GOOS ?= ${HOSTOS}
GOARCH ?= ${HOSTARCH}

# Set the execution extension for Windows.
ifeq (${GOOS},windows)
    EXE := .exe
endif

OS_ARCH := $(GOOS)_$(GOARCH)$(EXE)

LDFLAGS = -X stock-price-indexer/version.Version=$(FULL_VERSION)

TAGS = nocgo


INDEXER := ${TARGETDIR}/indexer_$(OS_ARCH)


build: build/indexer

build/indexer:
	@echo "+ $@"
	${GO} build -tags "$(TAGS)" -ldflags "$(LDFLAGS)" -o ${INDEXER} ${GOCMD}/indexer

vet:
	@echo "+ $@"
	@go vet ./...

fmt:
	@echo "+ $@"
	@test -z "$$(gofmt -s -l . 2>&1 | grep -v ^vendor/ | tee /dev/stderr)" || \
		(echo >&2 "+ please format Go code with 'gofmt -s'" && false)

lint:
	@echo "+ $@"
	@golangci-lint run --config=./.golangci.yml --timeout 300s

test:
	@echo "+ $@"
	@go test ./... -cover

test-race:
	@echo "+ $@"
	@go test ./... -cover --race

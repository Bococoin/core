VERSION := $(shell echo $(shell git describe --tags) | sed 's/^v//')
COMMIT := $(shell git log -1 --format='%H')
BINDIR ?= $(GOPATH)/bin

ldflags = -X github.com/cosmos/cosmos-sdk/version.Name=BocoCoin \
	-X github.com/cosmos/cosmos-sdk/version.ServerName=bocod \
	-X github.com/cosmos/cosmos-sdk/version.ClientName=bococli \
	-X github.com/cosmos/cosmos-sdk/version.Version=$(VERSION) \
	-X github.com/cosmos/cosmos-sdk/version.Commit=$(COMMIT)

BUILD_FLAGS := -ldflags '$(ldflags)'

export GO111MODULE=on

# The below include contains the tools and runsim targets.
include contrib/devtools/Makefile

all: up-swagger go.sum wnd linux

wnd:
	go install  $(BUILD_FLAGS) ./cmd/bocod
	go install  $(BUILD_FLAGS) ./cmd/bococli

go.sum: go.mod
	@echo "--> Ensure dependencies have not been modified"
	go mod verify

up-swagger: statik
	$(BINDIR)/statik -src=./client/lcd/swagger-ui -dest=./client/lcd -f -m
	@if [ -n "$(git status --porcelain)" ]; then \
        echo "Swagger docs are out of sync!!!";\
        exit 1;\
    else \
    	echo "Swagger docs are in sync";\
    fi
.PHONY: update-swagger-docs


linux:
	env GOOS=linux GOARCH=amd64 go install  $(BUILD_FLAGS) ./cmd/bocod
	env GOOS=linux GOARCH=amd64 go install  $(BUILD_FLAGS) ./cmd/bococli

install_debug: go.sum
	go install -mod=readonly $(BUILD_FLAGS) ./cmd/bocobug

test:
	@go test -mod=readonly $(PACKAGES)

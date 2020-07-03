
VERSION := $(shell echo $(shell git describe --tags) | sed 's/^v//')
COMMIT := $(shell git log -1 --format='%H')

ldflags = -X github.com/cosmos/cosmos-sdk/version.Name=BocoCoin \
	-X github.com/cosmos/cosmos-sdk/version.ServerName=bocod \
	-X github.com/cosmos/cosmos-sdk/version.ClientName=bococli \
	-X github.com/cosmos/cosmos-sdk/version.Version=$(VERSION) \
	-X github.com/cosmos/cosmos-sdk/version.Commit=$(COMMIT)

BUILD_FLAGS := -ldflags '$(ldflags)'

export GO111MODULE=on

all: install

install: go.sum
		go install  $(BUILD_FLAGS) ./cmd/bocod
		go install  $(BUILD_FLAGS) ./cmd/bococli

go.sum: go.mod
		@echo "--> Ensure dependencies have not been modified"
		go mod verify

linux: go.sum

	env GOOS=linux GOARCH=amd64 go install  $(BUILD_FLAGS) ./cmd/bocod
	env GOOS=linux GOARCH=amd64 go install  $(BUILD_FLAGS) ./cmd/bococli

install_debug: go.sum
	go install -mod=readonly $(BUILD_FLAGS) ./cmd/bocobug

test:
	@go test -mod=readonly $(PACKAGES)

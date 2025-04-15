# Makefile
.PHONY: all lint build
all: build

# Project name
PROJECT_NAME := cobravsviper
GO_MODULE_NAME := "github.com/nicop311/$(PROJECT_NAME)"

# Useful variables for build metadata
VERSION ?= $(shell git describe --tags --always --dirty)
COMMIT_LONG ?= $(shell git rev-parse HEAD)
COMMIT_SHORT ?= $(shell git rev-parse --short=8 HEAD)
COMMIT_TIMESTAMP := $(shell git show -s --format=%cI HEAD)
GO_VERSION ?= $(shell go version)
BUILD_PLATFORM  ?= $(shell uname -m)
BUILD_DATE ?= $(shell date -u --iso-8601=seconds)
LDFLAGS = "-X '$(GO_MODULE_NAME)/pkg/version.RawGitDescribe=$(VERSION)' -X '$(GO_MODULE_NAME)/pkg/version.GitCommitIdLong=$(COMMIT_LONG)' -X '$(GO_MODULE_NAME)/pkg/version.GitCommitIdShort=$(COMMIT_SHORT)' -X '$(GO_MODULE_NAME)/pkg/version.GoVersion=$(GO_VERSION)' -X '$(GO_MODULE_NAME)/pkg/version.BuildPlatform=$(BUILD_PLATFORM)' -X '$(GO_MODULE_NAME)/pkg/version.BuildDate=$(BUILD_DATE)' -X '$(GO_MODULE_NAME)/pkg/version.GitCommitTimestamp=$(COMMIT_TIMESTAMP)'"
GO_LDFLAGS = -ldflags=$(LDFLAGS)
BINARY_NAME = $(PROJECT_NAME)

lint:
	@golangci-lint run

build:
	@go version
	@go build $(GO_LDFLAGS) -o cobravsviper main.go

build-debug:
	@go version
	@go build -gcflags="all=-N -l" $(GO_LDFLAGS) -o cobravsviper main.go

get-ldflags:
	@echo $(LDFLAGS)
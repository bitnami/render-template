SHELL := /bin/bash
GOPATH ?= $(shell go env GOPATH)
PATH := $(GOPATH)/bin:$(PATH)

BUILD_DIR := $(abspath ./out)
TOOL_NAME ?= $(shell basename $(CURDIR))
TOOL_PATH ?= $(BUILD_DIR)/$(TOOL_NAME)

BUILD_DATE := $(shell date -u '+%Y-%m-%d %I:%M:%S UTC' 2> /dev/null)
GIT_HASH := $(shell git rev-parse HEAD 2> /dev/null)
LDFLAGS="-X 'main.buildDate=$(BUILD_DATE)' -X main.commit=$(GIT_HASH) -s -w"

DEBUG ?= 0

ifeq ($(DEBUG),1)
GO_TEST := @go test -v
else
GO_TEST := @go test
endif

GO_MOD := @go mod

# Do not do goimport of the vendor dir
go_files=$$(find $(1) -type f -name '*.go' -not -path "./vendor/*")
fmtcheck = @if goimports -l $(go_files) | read var; then echo "goimports check failed for $(1):\n `goimports -d $(go_files)`"; exit 1; fi

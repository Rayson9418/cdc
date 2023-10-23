APP := examples
VERSION := $(shell date +'%Y%m%d')

MKFILE_PATH := $(abspath $(lastword $(MAKEFILE_LIST)))
MKFILE_DIR := $(dir $(MKFILE_PATH))

BASEDIR := $(shell cd ${MKFILE_DIR})

.PHONY: default
default:run

.PHONY: build
build:export GO111MODULE=on
build:export CGO_ENABLED=0
build:export GOSUMDB=off
build:
	@go mod tidy -go=1.21
	@echo "building $(APP):$(VERSION)..."
	@go build -o $(APP)

.PHONY: run
run:
	@go run ./
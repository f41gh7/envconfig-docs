GOCMD=GO111MODULE=on go
GOOS ?= $(shell go env GOOS)
GOARCH ?= $(shell go env GOARCH)
GOBUILD=CGO_ENABLED=0 GOOS=${GOOS} GOARCH=${GOARCH}  $(GOCMD) build -trimpath
GOCLEAN=$(GOCMD) clean
BINARY_NAME=envconfig-docs
REPO=github.com/f41gh7/envconfig-docs


build:
	$(GOBUILD) -o ./bin/$(BINARY_NAME) $(REPO)

GOCMD=GO111MODULE=on go
GOOS ?= linux
GOARCH ?= amd64
GOBUILD=CGO_ENABLED=0 GOOS=${GOOS} GOARCH=${GOARCH}  $(GOCMD) build -trimpath
GOCLEAN=$(GOCMD) clean
BINARY_NAME=envconfig-docs
REPO=github.com/f41gh7/envconfig-docs


build:
	$(GOBUILD) -o $(BINARY_NAME) $(REPO)
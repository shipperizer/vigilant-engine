.PHONY=build proto proto-deps

GO111MODULE?=on
GOPRIVATE?=github.com/shipperizer/*
CGO_ENABLED?=0
GOOS?=linux
GOARCH?=amd64
GO?=go
APP_NAME?=app

.EXPORT_ALL_VARIABLES:

build:
	$(MAKE) -C cmd/$(APP_NAME) build

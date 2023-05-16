PROJECT ?= $(shell basename $(CURDIR))
MODULE  ?= $(shell go list -m)

GO      ?= GO111MODULE=on go
VERSION ?= $(shell git describe --tags 2>/dev/null || git rev-parse --short HEAD)
BIDTIME ?= $(shell date +%FT%T%z)

BITTAGS :=
LDFLAGS := -s -w
LDFLAGS += -X "github.com/starudream/go-lib/constant.VERSION=$(VERSION)"
LDFLAGS += -X "github.com/starudream/go-lib/constant.BIDTIME=$(BIDTIME)"
LDFLAGS += -X "github.com/starudream/go-lib/constant.NAME="
LDFLAGS += -X "github.com/starudream/go-lib/constant.PREFIX="

.PHONY: bin

bin:
	@$(MAKE) tidy
	CGO_ENABLED=0 $(GO) build -tags '$(BITTAGS)' -ldflags '$(LDFLAGS)' -o bin/app $(MODULE)/cmd

run:
	@$(MAKE) tidy
	CGO_ENABLED=1 $(GO) run -race -tags '$(BITTAGS)' -ldflags '$(LDFLAGS)' $(MODULE)/cmd

tidy:
	$(GO) mod tidy

clean:
	rm -rf bin/*

upx:
	upx bin/*

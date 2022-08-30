PROJECT ?= $(shell basename $(CURDIR))
MODULE  ?= $(shell go list -m)

GO      ?= GO111MODULE=on go
VERSION ?= $(shell git describe --tags 2>/dev/null || echo "dev")
BIDTIME ?= $(shell date +%FT%T%z)

BITTAGS := viper_logger
LDFLAGS := -s -w
LDFLAGS += -X "github.com/starudream/go-lib/constant.VERSION=$(VERSION)"
LDFLAGS += -X "github.com/starudream/go-lib/constant.BIDTIME=$(BIDTIME)"
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

lint:
	golangci-lint run --skip-dirs-use-default

PLATFORMS = linux-amd64 linux-arm64 darwin-amd64 darwin-arm64 windows-amd64 windows-arm64

zip_releases = $(addsuffix .zip, $(PLATFORMS))

$(zip_releases): %.zip: %
	@if test $(findstring windows, $@); then \
		zip -m -j bin/$(PROJECT)-$(basename $@)-$(VERSION).zip bin/$(PROJECT)-$(basename $@).exe; \
	else \
		chmod +x bin/$(PROJECT)-$(basename $@); \
		zip -m -j bin/$(PROJECT)-$(basename $@)-$(VERSION).zip bin/$(PROJECT)-$(basename $@); \
	fi

releases: clean $(zip_releases)

linux-amd64:
	GOOS=linux GOARCH=amd64 CGO_ENABLED=0 $(GO) build -tags '$(BITTAGS)' -ldflags '$(LDFLAGS)' -o bin/$(PROJECT)-$@ $(MODULE)/cmd

linux-arm64:
	GOOS=linux GOARCH=arm64 CGO_ENABLED=0 $(GO) build -tags '$(BITTAGS)' -ldflags '$(LDFLAGS)' -o bin/$(PROJECT)-$@ $(MODULE)/cmd

darwin-amd64:
	GOOS=darwin GOARCH=amd64 CGO_ENABLED=0 $(GO) build -tags '$(BITTAGS)' -ldflags '$(LDFLAGS)' -o bin/$(PROJECT)-$@ $(MODULE)/cmd

darwin-arm64:
	GOOS=darwin GOARCH=arm64 CGO_ENABLED=0 $(GO) build -tags '$(BITTAGS)' -ldflags '$(LDFLAGS)' -o bin/$(PROJECT)-$@ $(MODULE)/cmd

windows-amd64:
	GOOS=windows GOARCH=amd64 CGO_ENABLED=0 $(GO) build -tags '$(BITTAGS)' -ldflags '$(LDFLAGS)' -o bin/$(PROJECT)-$@.exe $(MODULE)/cmd

windows-arm64:
	GOOS=windows GOARCH=arm64 CGO_ENABLED=0 $(GO) build -tags '$(BITTAGS)' -ldflags '$(LDFLAGS)' -o bin/$(PROJECT)-$@.exe $(MODULE)/cmd

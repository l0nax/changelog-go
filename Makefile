VERSION := $(shell git describe --long --dirty --tags)
DATE := $(date)
COMMIT := $(shell git rev-parse HEAD)
LDFLAGS := -ldflags "-X=internal.Version=$(VERSION) -X=internal.BuildTime=$(DATE) -X=internal.Hash=$(COMMIT)"

.PHONY: all
all: build-dev

.PHONY: build-dev
build-dev:
	go build $(LDFLAGS) -o changelog-go

.PHONY: release
release:
	@echo "TO BE DONE!!!"

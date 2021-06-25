VERSION := $(shell git describe --long --dirty --tags)
DATE := $(date)
COMMIT := $(shell git rev-parse HEAD)
LDFLAGS := -ldflags "-X 'gitlab.com/l0nax/changelog-go/pkg/version.Version=$(VERSION)' -X 'gitlab.com/l0nax/changelog-go/pkg/version.BuildTime=$(DATE)' -X 'gitlab.com/l0nax/changelog-go/pkg/version.Hash=$(COMMIT)'"

.PHONY: all
all: build-dev

.PHONY: build-dev
build-dev:
	go build $(LDFLAGS) -o changelog-go

.PHONY: release
release:
	goreleaser release
	godownloader --repo=l0nax/changelog-go > install.sh

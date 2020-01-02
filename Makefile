VERSION := $(shell git describe --long --dirty --tags)
DATE := $(date)
COMMIT := $(shell git rev-parse HEAD)
LDFLAGS := -ldflags "-X 'gitlab.com/l0nax/changelog-go/pkg/version.Version=$(VERSION)' -X 'gitlab.com/l0nax/changelog-go/pkg/version.BuildTime=$(DATE)' -X 'gitlab.com/l0nax/changelog-go/pkg/version.Hash=$(COMMIT)'"

.PHONY: all
all: build-dev

.PHONY: build-dev
build-dev: assets
	go build $(LDFLAGS) -o changelog-go

.PHONY: release
release:
	goreleaser release
	godownloader --repo=l0nax/changelog-go > install.sh

.PHONY: assets
assets:
	resources -declare -fmt -output internal/assets.go -package internal example/.changelog-go.yaml

.PHONY: depends
depends:
	go get github.com/omeid/go-resources/cmd/resources

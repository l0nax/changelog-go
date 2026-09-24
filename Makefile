# The build carries no -ldflags stamps: the Go toolchain records the version,
# commit and dirtiness on its own, and internal/version reads them back.

.PHONY: all
all: build-dev

.PHONY: build-dev
build-dev:
	go build -o changelog-go

.PHONY: test
test:
	go test ./...

.PHONY: lint
lint:
	golangci-lint run ./...

.PHONY: release
release:
	goreleaser release

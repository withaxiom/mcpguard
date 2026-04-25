.PHONY: build test lint clean install snapshot release

VERSION ?= dev
COMMIT  := $(shell git rev-parse --short HEAD 2>/dev/null || echo "none")
DATE    := $(shell date -u +%Y-%m-%dT%H:%M:%SZ)
LDFLAGS := -s -w \
	-X github.com/withaxiom/mcpguard/cmd.Version=$(VERSION) \
	-X github.com/withaxiom/mcpguard/cmd.Commit=$(COMMIT) \
	-X github.com/withaxiom/mcpguard/cmd.Date=$(DATE) \
	-X github.com/withaxiom/mcpguard/cmd.BuiltBy=make

build:
	go build -ldflags "$(LDFLAGS)" -o bin/mcpguard .

install:
	go install -ldflags "$(LDFLAGS)" .

test:
	go test ./... -v -race

lint:
	golangci-lint run ./...

clean:
	rm -rf bin/ dist/

# GoReleaser local snapshot (no publish)
snapshot:
	goreleaser release --snapshot --clean

# Full release (requires GITHUB_TOKEN)
release:
	goreleaser release --clean

# Quick dev build + run scan
dev: build
	./bin/mcpguard scan -v

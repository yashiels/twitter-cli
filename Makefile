BINARY := twt
CMD    := ./cmd/twt

.PHONY: build clean install lint fmt vet test

build:
	go build -ldflags "-s -w" -o $(BINARY) $(CMD)

clean:
	rm -f $(BINARY)

install:
	go install $(CMD)

lint:
	golangci-lint run --timeout=5m

fmt:
	golangci-lint run --fix --timeout=5m

vet:
	go vet ./...

test:
	go test ./...

.DEFAULT_GOAL := build

hooks:
	git config core.hooksPath .githooks
	chmod +x .githooks/pre-commit
	@echo "✓ Git hooks installed"

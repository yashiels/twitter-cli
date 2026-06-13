BINARY := twt
CMD    := ./cmd/twt

.PHONY: build clean install lint vet

build:
	go build -ldflags "-s -w" -o $(BINARY) $(CMD)

clean:
	rm -f $(BINARY)

install:
	go install $(CMD)

lint:
	go vet ./...

vet:
	go vet ./...

test:
	go test ./...

.DEFAULT_GOAL := build

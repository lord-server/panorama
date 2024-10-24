GOEXE=$(shell go env GOEXE)

.PHONY: all
all: panorama

.PHONY: panorama
panorama:
	go build -o bin/panorama${GOEXE} ./cmd/panorama

.PHONY: lint
lint:
	golangci-lint run --fix=false --color=always

.PHONY: clean
clean:
	rm -rf bin

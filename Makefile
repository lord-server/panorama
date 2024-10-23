GOEXE=$(shell go env GOEXE)
PANORAMA_BIN=bin/panorama${GOEXE}
SRCS:=$(shell find . -wholename 'internal/**/*.go')

.PHONY: all
all: panorama

.PHONY: panorama
panorama:
	go build -o ${PANORAMA_BIN} ./cmd/panorama

.PHONY: lint
lint:
	golangci-lint run --fix=false --color=always

.PHONY: clean
clean:
	rm -rf bin

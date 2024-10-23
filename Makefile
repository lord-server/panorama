GOEXE=$(shell go env GOEXE)
PANORAMA_BIN=bin/panorama${GOEXE}

.PHONY: all
all: ${PANORAMA_BIN}

bin/panorama: \
		cmd/panorama/* \
		internal/**/* \
		static/**/*
	go build -o ${PANORAMA_BIN} ./cmd/panorama

.PHONY: lint
lint:
	golangci-lint run --fix=false --color=always

.PHONY: clean
clean:
	rm -rf bin

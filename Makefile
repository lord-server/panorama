.PHONY:bin
bin:
	go build ./cmd/panorama

.PHONY:all
all: bin
	$(MAKE) -C ui all

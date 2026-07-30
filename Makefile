BINARY := goelfcheck

.PHONY: build test clean

build:
	go build -o $(BINARY) ./cmd/goelfcheck

test:
	go test ./...

clean:
	rm -f $(BINARY) testdata/sample_default testdata/sample_hardened

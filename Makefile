.PHONY: build
build:
	go build -o bin/ ./cmd/...

.PHONY: clean
clean:
	rm -rf bin/

.PHONY: build test

build:
	echo "Building..."
	go build -o ./bin/hackerone-target-retrieval ./main.go; \
	echo "Built to ./bin/hackerone-target-retrieval"; \
	echo "\n\n";

test:
	go fmt ./...; \
	go build ./...; \
	go test -v ./...; \
	echo "\n\n";

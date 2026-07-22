.PHONY: build run test clean lint

BINARY=es-meili-growi-bridge

build:
	go build -o $(BINARY) ./cmd/$(BINARY)/

run: build
	./$(BINARY)

test:
	go test ./... -v -count=1

test-short:
	go test ./... -count=1

lint:
	go vet ./...

clean:
	rm -f $(BINARY)

docker-build:
	docker build -t $(BINARY) .

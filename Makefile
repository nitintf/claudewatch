BINARY_NAME=claudewatch

.PHONY: all build test lint fmt vet clean tidy install

all: lint test build

build:
	go build -o $(BINARY_NAME) .

test:
	go test -race -coverprofile=coverage.txt -covermode=atomic ./...

lint:
	golangci-lint run ./...

fmt:
	gofmt -s -w .

vet:
	go vet ./...

clean:
	rm -f $(BINARY_NAME) coverage.txt

tidy:
	go mod tidy

install: build
	cp $(BINARY_NAME) $(GOPATH)/bin/$(BINARY_NAME)
	./$(BINARY_NAME) install

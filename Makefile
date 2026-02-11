.PHONY: all default fmt lint clean build install

default: all

fmt:
	go fmt ./...

lint:
	golint ./...
	go vet ./...

clean:
	go clean -i ./...
	rm -fv ./bin/taskusama || true

# builds binaries into ./bin/
build:
	mkdir -p bin
	go build -o bin/taskusama  ./cmd/taskusama

# installs binaries into $GOBIN
install:
	go install ./cmd/taskusama

# all
all: fmt lint clean install build


.PHONY: all default fmt lint clean build install image docker

# project info
IMAGE   := fakelassian/taskusama
VERSION := 1.0.0

# sets reproducible build flags
GOFLAGS := -trimpath -buildvcs=false
LDFLAGS := -s -w

default: all

fmt:
	go fmt ./...

lint:
	@command -v golint >/dev/null 2>&1 || { \
		echo "golint not installed. Run: go install golang.org/x/lint/golint@latest"; \
		exit 1; \
	}
	golint ./...
	go vet ./...

clean:
	go clean ./...
	rm -f ./bin/taskusama

# builds binaries into ./bin/
build:
	mkdir -p bin
	GOFLAGS="$(GOFLAGS)" \
	go build \
		-ldflags "$(LDFLAGS)" \
		-o bin/taskusama \
		./cmd/server

# installs binary into $GOBIN or $GOPATH/bin
install:
	GOFLAGS="$(GOFLAGS)" \
	go install \
		-ldflags "$(LDFLAGS)" \
		./cmd/server

# builds Docker image
image:
	docker build -t $(IMAGE):$(VERSION) .
	docker image tag $(IMAGE):$(VERSION) $(IMAGE):latest

# alias
docker: image

# all
all: fmt lint clean build install


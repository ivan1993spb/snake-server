
EXECUTABLES=git go find docker tar

_=$(foreach exec,$(EXECUTABLES), \
	$(if $(shell which $(exec)), ok, $(error "No $(exec) in PATH")))

IMAGE=ivan1993spb/snake-server

IMAGE_GOLANG=golang:1.16.4-alpine3.13
IMAGE_ALPINE=alpine:3.13

REPO=github.com/ivan1993spb/snake-server

DEFAULT_GOOS=linux
DEFAULT_GOARCH=amd64

BINARY_NAME=snake-server
VERSION=$(shell git describe --tags --abbrev=0)
BUILD=$(shell git rev-parse --short HEAD)

PLATFORMS=darwin linux windows
ARCHITECTURES=386 amd64

LDFLAGS=-ldflags "-s -w -X main.Version=$(VERSION) -X main.Build=$(BUILD)"
DOCKER_BUILD_ARGS=\
 --build-arg VERSION=$(VERSION) \
 --build-arg BUILD=$(BUILD) \
 --build-arg IMAGE_GOLANG=$(IMAGE_GOLANG) \
 --build-arg IMAGE_ALPINE=$(IMAGE_ALPINE)

default: build

docker/build:
	@docker build $(DOCKER_BUILD_ARGS) -t $(IMAGE):$(VERSION) .
	@docker tag $(IMAGE):$(VERSION) $(IMAGE):latest
	@echo "Build $(BUILD) tagged $(IMAGE):$(VERSION)"
	@echo "Build $(BUILD) tagged $(IMAGE):latest"

docker/push:
	@echo "Push build $(BUILD) with tag $(IMAGE):$(VERSION)"
	@docker push $(IMAGE):$(VERSION)
	@echo "Push build $(BUILD) with tag $(IMAGE):latest"
	@docker push $(IMAGE):latest

go/vet:
	@docker run --rm -v $(PWD):/go/src/$(REPO) -w /go/src/$(REPO) \
		-e CGO_ENABLED=0 $(IMAGE_GOLANG) go vet ./...

go/test:
	@docker run --rm -v $(PWD):/go/src/$(REPO) -w /go/src/$(REPO) \
		-e CGO_ENABLED=0 $(IMAGE_GOLANG) \
		go test -v -cover ./...

go/test/benchmarks:
	@docker run --rm -v $(PWD):/go/src/$(REPO) -w /go/src/$(REPO) \
		-e CGO_ENABLED=0 $(IMAGE_GOLANG) \
		go test -bench . -timeout 1h ./...

go/build:
	@docker run --rm -v $(PWD):/go/src/$(REPO) -w /go/src/$(REPO) \
		-e GOOS=$(DEFAULT_GOOS) -e GOARCH=$(DEFAULT_GOARCH) \
		-e CGO_ENABLED=0 $(IMAGE_GOLANG) \
		go build $(LDFLAGS) -v -o $(BINARY_NAME)

go/crosscompile:
	@_=$(foreach GOOS, $(PLATFORMS), \
		$(foreach GOARCH, $(ARCHITECTURES), \
			$(shell docker run --rm \
				-v $(PWD):/go/src/$(REPO) \
				-w /go/src/$(REPO) \
				-e GOOS=$(GOOS) \
				-e GOARCH=$(GOARCH) \
				-e CGO_ENABLED=0 \
				$(IMAGE_GOLANG) go build $(LDFLAGS) -o $(BINARY_NAME)-$(VERSION)-$(GOOS)-$(GOARCH)) \
		) \
	)
	@_=$(foreach GOOS, $(PLATFORMS), \
		$(foreach GOARCH, $(ARCHITECTURES), \
			$(shell tar -zcf \
				$(BINARY_NAME)-$(VERSION)-$(GOOS)-$(GOARCH).tar.gz \
				--transform="flags=r;s|-$(VERSION)-$(GOOS)-$(GOARCH)||" \
				$(BINARY_NAME)-$(VERSION)-$(GOOS)-$(GOARCH)) \
		) \
	)
	@echo -n

build:
	@go build $(LDFLAGS) -v -o $(BINARY_NAME)

install:
	@go install $(LDFLAGS) -v

clean:
	@find -maxdepth 1 -type f -name '${BINARY_NAME}*' -print -delete

coverprofile:
	@go test -coverprofile=coverage.out ./...
	@go tool cover -func=coverage.out
	@go tool cover -html=coverage.out

go/generate:
	@go list ./... | grep -v vendor | xargs go generate -v

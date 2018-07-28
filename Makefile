
EXECUTABLES=git go find docker

_=$(foreach exec,$(EXECUTABLES), \
	$(if $(shell which $(exec)), ok, $(error "No $(exec) in PATH")))

IMAGE=ivan1993spb/snake-server
IMAGE_GOLANG=golang:1.10-alpine3.7

REPO=github.com/ivan1993spb/snake-server

BINARY_NAME=snake-server
VERSION=$(shell git describe --tags --abbrev=0)
BUILD=$(shell git rev-parse --short HEAD)

PLATFORMS=darwin linux windows
ARCHITECTURES=386 amd64

LDFLAGS=-ldflags "-X main.Version=$(VERSION) -X main.Build=$(BUILD)"
BUILD_ARGS=--build-arg VERSION=$(VERSION) --build-arg BUILD=$(BUILD)

default: build

docker/build:
	@docker build $(BUILD_ARGS) -t $(IMAGE):$(VERSION) .
	@docker tag $(IMAGE):$(VERSION) $(IMAGE):latest
	@echo "Build $(BUILD) tagged $(IMAGE):$(VERSION)"
	@echo "Build $(BUILD) tagged $(IMAGE):latest"

docker/push:
	@echo "Push build $(BUILD) with tag $(IMAGE):$(VERSION)"
	@docker push $(IMAGE):$(VERSION)
	@echo "Push build $(BUILD) with tag $(IMAGE):latest"
	@docker push $(IMAGE):latest

go/vet:
	@docker run --rm -v $(PWD):/go/src/$(REPO) -w /go/src/$(REPO) $(IMAGE_GOLANG) \
		sh -c "go list ./... | grep -v vendor | xargs go vet"

go/test:
	@docker run --rm -v $(PWD):/go/src/$(REPO) -w /go/src/$(REPO) $(IMAGE_GOLANG) \
		sh -c "go list ./... | grep -v vendor | xargs go test -v"

go/crosscompile:
	$(foreach GOOS, $(PLATFORMS),\
		$(foreach GOARCH, $(ARCHITECTURES), $(shell docker run --rm \
			-v $(PWD):/go/src/$(REPO) \
			-w /go/src/$(REPO) \
			-e GOOS=$(GOOS) \
			-e GOARCH=$(GOARCH) \
			$(IMAGE_GOLANG) go build $(LDFLAGS) -o $(BINARY_NAME)-$(VERSION)-$(GOOS)-$(GOARCH))))
	@echo -n

build:
	@go build $(LDFLAGS) -v -o $(BINARY_NAME)

install:
	@go install $(LDFLAGS) -v

clean:
	@find -maxdepth 1 -type f -name '${BINARY_NAME}*' -print -delete

go/generate:
	@go list ./... | grep -v vendor | xargs go generate -v

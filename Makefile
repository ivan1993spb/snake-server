IMAGE=ivan1993spb/snake-server:latest
IMAGE_GOLANG=golang:1.9.5-alpine3.7
REPO=github.com/ivan1993spb/snake-server

docker/build:
	@docker build -t $(IMAGE) .

docker/push:
	@docker push $(IMAGE)

go/vet:
	@docker run --rm -v $(PWD):/go/src/$(REPO) -w /go/src/$(REPO) $(IMAGE_GOLANG) sh -c "go list ./... | grep -v vendor | xargs go vet"

go/test:
	@docker run --rm -v $(PWD):/go/src/$(REPO) -w /go/src/$(REPO) $(IMAGE_GOLANG) sh -c "go list ./... | grep -v vendor | xargs go test -v"

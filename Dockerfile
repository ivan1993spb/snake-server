
# See Makefile
ARG IMAGE_GOLANG
ARG IMAGE_ALPINE

FROM $IMAGE_GOLANG AS intermediate

ARG VERSION=docker
ARG BUILD=docker

WORKDIR /go/src/github.com/ivan1993spb/snake-server

COPY . .

RUN CGO_ENABLED=0 go build -ldflags "-X main.Version=$VERSION -X main.Build=$BUILD" \
    -v -x -o /snake-server github.com/ivan1993spb/snake-server

FROM $IMAGE_ALPINE

COPY --from=intermediate /snake-server /usr/local/bin/snake-server

ENTRYPOINT ["snake-server"]

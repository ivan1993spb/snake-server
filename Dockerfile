FROM golang:1.9.5-alpine3.7 AS intermediate

ARG VERSION=docker
ARG BUILD=docker

WORKDIR /go/src/github.com/ivan1993spb/snake-server

COPY . .

RUN go build -ldflags "-X main.Version=$VERSION -X main.Build=$BUILD" -v -x -o /snake-server github.com/ivan1993spb/snake-server

FROM alpine:3.7

COPY --from=intermediate /snake-server /usr/local/bin/snake-server

ENTRYPOINT ["snake-server"]

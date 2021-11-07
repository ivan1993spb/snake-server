
# See Makefile
ARG IMAGE_GOLANG=golang:1.16.4-alpine3.13
ARG IMAGE_ALPINE=alpine:3.13

FROM $IMAGE_ALPINE AS helper

RUN adduser -u 10001 -h /dev/null -H -D -s /sbin/nologin snake

RUN sed -i '/^snake/!d' /etc/passwd

FROM $IMAGE_GOLANG AS builder

ARG VERSION=unknown
ARG BUILD=unknown

WORKDIR /go/src/github.com/ivan1993spb/snake-server

COPY . .

ENV CGO_ENABLED=0

RUN go build -ldflags "-s -w -X main.Version=$VERSION -X main.Build=$BUILD" \
    -v -x -o /snake-server github.com/ivan1993spb/snake-server

FROM scratch

COPY --from=helper /etc/passwd /etc/passwd

USER snake

COPY --from=builder /snake-server /usr/local/bin/snake-server

ENTRYPOINT ["snake-server"]

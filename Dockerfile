
# See Makefile
ARG IMAGE_GOLANG=golang:1.16.4-alpine3.13
ARG IMAGE_ALPINE=alpine:3.13
ARG IMAGE_CLIENT=ivan1993spb/snake-lightweight-client:1.4.0

FROM $IMAGE_CLIENT AS client

FROM $IMAGE_ALPINE AS helper

RUN adduser -u 10001 -h /dev/null -H -D -s /sbin/nologin snake

RUN sed -i '/^snake/!d' /etc/passwd

FROM $IMAGE_GOLANG AS builder

ARG VERSION=unknown
ARG BUILD=unknown

WORKDIR /go/src/github.com/ivan1993spb/snake-server

COPY . .

COPY --from=client \
    /usr/local/share/snake-lightweight-client \
    client/public/dist

ENV CGO_ENABLED=0

RUN go build \
    -ldflags "-s -w -X main.Version=${VERSION} -X main.Build=${BUILD::7}" \
    -v -x -o /snake-server github.com/ivan1993spb/snake-server

FROM scratch

COPY --from=helper /etc/passwd /etc/passwd

USER snake

COPY --from=builder /snake-server /usr/local/bin/snake-server

ENTRYPOINT ["snake-server"]


# See Makefile
ARG IMAGE_GOLANG=golang:1.14-alpine3.11

FROM $IMAGE_GOLANG AS intermediate

ARG VERSION=unknown
ARG BUILD=unknown

WORKDIR /go/src/github.com/ivan1993spb/snake-server

COPY . .

ENV CGO_ENABLED=0

RUN go build -ldflags "-X main.Version=$VERSION -X main.Build=$BUILD" \
    -v -x -o /snake-server github.com/ivan1993spb/snake-server

FROM scratch

COPY --from=intermediate /snake-server /usr/local/bin/snake-server

ENTRYPOINT ["snake-server"]

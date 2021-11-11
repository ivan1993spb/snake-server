
# Snake-Server

![Build Status](https://github.com/ivan1993spb/snake-server/actions/workflows/go.yml/badge.svg?branch=master)
[![GitHub release](https://img.shields.io/github/release/ivan1993spb/snake-server.svg)](https://github.com/ivan1993spb/snake-server/releases/latest)
[![Go Report Card](https://goreportcard.com/badge/github.com/ivan1993spb/snake-server)](https://goreportcard.com/report/github.com/ivan1993spb/snake-server)
[![Docker Pulls](https://img.shields.io/docker/pulls/ivan1993spb/snake-server)](https://hub.docker.com/r/ivan1993spb/snake-server)

Snake-Server is a server for multiplayer snake game. You can play with your
friends! The special feature is that you can eat small snakes!

Take a look at a working instance here - https://snakeonline.xyz

[![Game demo](demo.gif)](https://snakeonline.xyz)

## Usage

1. `docker run --rm -p 8080:8080 ivan1993spb/snake-server --enable-web`
2. Open in the browser http://localhost:8080/.

## How to play?

* You control a snake
* You need to grow the biggest snake
* You can eat apples, mice, watermelons, small and dead snakes
* If the snake dies, you will have to start over

## Installation

* **Go get**

  ```
  go get github.com/ivan1993spb/snake-server@latest
  snake-server -h
  ```

* **Docker**

  Check out [**the repo**](https://hub.docker.com/r/ivan1993spb/snake-server).

  ```bash
  docker pull ivan1993spb/snake-server

  docker run --rm -p 8080:8080 ivan1993spb/snake-server --enable-web

  docker run --rm ivan1993spb/snake-server -h
  ```

* **Download and install the binary**

  Take a look at the [**release page**](https://github.com/ivan1993spb/snake-server/releases/latest)

  Curl:

  + Set *VERSION*, *PLATFORM* and *ARCHITECTURE*:
    ```bash
    VERSION=v4.3.0
    # darwin or linux or windows
    PLATFORM=linux
    # amd64 or 386
    ARCHITECTURE=amd64
    ```
  + Download and install the binary to `/usr/local/bin/`:
    ```bash
    curl -sL "https://github.com/ivan1993spb/snake-server/releases/download/${VERSION}/snake-server-${VERSION}-${PLATFORM}-${ARCHITECTURE}.tar.gz" |\
      tar xvz -C /usr/local/bin/
    ```

* **Deploy the server using the ansible playbook**

  https://github.com/ivan1993spb/snake-ansible.

## CLI options

Use `snake-server -h` for more information.

Options:

* `--address` - **string** - sets an address to listen and serve (default: *:8080*). For example: *:8080*, *localhost:7070*
* `--conns-limit` - **integer** - to limit the number of opened web-socket connections (default: *1000*)
* `--groups-limit` - **integer** - to limit the number of games for a server instance (default: *100*)
* `--enable-web` - **bool** - to enable the embedded web client (default: *false*)
* ~~`--enable-broadcast` - **bool** - to enable the broadcasting API method (default: *false*)~~
* `--forbid-cors` - **bool** - to forbid cross-origin resource sharing (default: *false*)
* `--log-json` - **bool** - to enable JSON log output format (default: *false*)
* `--log-level` - **string** - to set the log level: *panic*, *fatal*, *error*, *warning* (*warn*), *info* or *debug* (default: *info*)
* `--seed` - **integer** - to specify a random seed (default: *the number of nanoseconds elapsed since January 1, 1970 UTC*)
* `--sentry-enable` - **bool** - to enable sending logs to sentry (default: *false*)
* `--sentry-dsn` - **string** - sentry's DSN (default: ""). For example: `https://public@sentry.example.com/44`
* `--tls-cert` - **string** - to specify a path to a certificate file
* `--tls-enable` - **bool** - to enable TLS
* `--tls-key` - **string** - to specify a path to a key file
* `--debug` - **bool** - to enable profiling routes

## Clients

There is an embedded JavaScript web client compiled into the server.
You can enable it with CLI flag `--enable-web`.

You are always welcome to create your own client!

You can find examples here:

* VueJS client repo: https://github.com/ivan1993spb/snake-lightweight-client

  *This is the embedded web client*

* Python backend repo: https://github.com/ivan1993spb/snake-backend

  *Development is in progress*

See documentation [docs/api.md](docs/api.md) and [docs/websocket.md](docs/websocket.md).

REST API specification: [openapi.yaml](openapi.yaml).

## License

See [LICENSE](LICENSE).

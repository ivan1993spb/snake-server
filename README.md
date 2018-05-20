
# Snake server [![Build Status](https://travis-ci.org/ivan1993spb/snake-server.svg?branch=master)](https://travis-ci.org/ivan1993spb/snake-server) [![Go Report Card](https://goreportcard.com/badge/github.com/ivan1993spb/snake-server)](https://goreportcard.com/report/github.com/ivan1993spb/snake-server)

Server for online arcade game - snake.

// TODO: Create screen shot.

## Game rules

// TODO: Write game rules.

## Client

Known clients:

* Python client repo: https://github.com/ivan1993spb/snake-client

## Installation

### Install from source

* `go get -u github.com/ivan1993spb/snake-server` to load source code
* `go install github.com/ivan1993spb/snake-server` to install server
* `snake-server` to start server
* Use `snake-server -h` to see usage information

### Install from docker-hub

See docker-hub repo: https://hub.docker.com/r/ivan1993spb/snake-server

* Install docker: [use fast installation script](https://get.docker.com/)
* Choose image tag: https://hub.docker.com/r/ivan1993spb/snake-server/tags/
* Use `docker pull ivan1993spb/snake-server` to pull server image from docker hub
* `docker run --rm --net host --name snake-server ivan1993spb/snake-server` to start server
* `docker run --rm ivan1993spb/snake-server -h` for usage information

## CLI arguments

Use `snake-server -help` for help info.

Arguments:

* `--address` - **string** - address to serve (default: *:8080*). For example: *:8080*, *localhost:7070*
* `--groups-limit` - **int** - groups limit for server (default: *100*)
* `--conns-limit` - **int** - open web-socket connections limit (default: *1000*)
* `--seed` - **int** - random seed (default: the number of nanoseconds elapsed since January 1, 1970 UTC)
* `--log-json` - **bool** - set this flag to use JSON log format (default: *false*)
* `--log-level` - **string** - set log level: *panic*, *fatal*, *error*, *warning* (*warn*), *info* or *debug* (default: *info*)

## API Description

API methods provide JSON format.

### Request `POST /games`

Creates game and returns JSON details.

```
curl -X POST -d limit=3 -d width=100 -d height=100 http://localhost:8080/games
{"id":0,"limit":3,"width":100,"height":100}
```

### Request `GET /games`

Returns info about all games on server.

```
curl -X GET http://localhost:8080/games
{"games":[{"id":0,"limit":3,"count":0}]}
```

### Request `GET /games/{id}`

Returns game information.

```
curl -X GET http://localhost:8080/games/0
{"id":0,"limit":3,"count":0}
```

### Request `DELETE /games/{id}`

Deletes game if there is not players.

```
curl -X DELETE http://localhost:8080/games/0
{"id":0}
```

### Request `GET /games/{id}/ws`

Connects to game Web-Socket.

* Returns playground size: width and height `{width: w, height: h}`
* Initialize gamer's objects and session
* Returns snake id
* Returns all objects on playground
* Pushes events and objects from game

## Game objects

Primitives

* Area: `[width, height]`
* Direction: `"n"`, `"w"`, `"s"`, `"e"`
* Dot: `[x, y]`
* Dot list: `[[x, y], [x, y], [x, y], [x, y], [x, y], [x, y]]`
* Location: `[[x, y], [x, y], [x, y], [x, y], [x, y], [x, y]]`
* Rect: `[x, y, width, height]`

Game objects:

* Apple: `{"type": "apple", "id": 1, "dot": [x, y]}`
* Corpse: `{"type": "corpse", "id": 2, "dots": [[x, y], [x, y], [x, y]]}`
* Mouse: `{"type": "mouse", "id": 3, dot: [x, y], "dir": "n"}`
* Snake: `{"type": "snake", "id": 4, "dots": [[x, y], [x, y], [x, y]]}`
* Wall: `{"type": "wall", "id": 5, "dots": [[x, y], [x, y], [x, y]]}`
* Watermelon: `{"type": "watermelon", "id": 6, "dots": [[x, y], [x, y], [x, y]]}`

Message types:

* Object: `{"type": "object", "object": {}}` - delete, update or create
* Error: `{"type": "error", "message": "text"}`
* Notice: `{"type": "notice", "message": "text"}`

## License

See [LICENSE](LICENSE).

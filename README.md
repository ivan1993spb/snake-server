
# Snake-Server [![Build Status](https://travis-ci.org/ivan1993spb/snake-server.svg?branch=master)](https://travis-ci.org/ivan1993spb/snake-server) [![Go Report Card](https://goreportcard.com/badge/github.com/ivan1993spb/snake-server)](https://goreportcard.com/report/github.com/ivan1993spb/snake-server) [![Swagger Validator](https://img.shields.io/swagger/valid/2.0/https/raw.githubusercontent.com/ivan1993spb/snake-server/master/swagger.yml.svg)](https://raw.githubusercontent.com/ivan1993spb/snake-server/master/swagger.yml) [![GitHub release](https://img.shields.io/github/release/ivan1993spb/snake-server/all.svg)](https://github.com/ivan1993spb/snake-server/releases/latest) [![license](https://img.shields.io/github/license/ivan1993spb/snake-server.svg)](LICENSE)

Snake-Server is the server for online arcade game - snake.

## Table of contents

- [Game rules](#game-rules)
- [Installation](#installation)
    * [Download and install binary](#download-and-install-binary)
    * [Build and install latest server](#build-and-install-latest-server)
    * [Pull server image from docker-hub](#pull-server-image-from-docker-hub)
    * [Build and install server of specific version from source code](#build-and-install-server-of-specific-version-from-source-code)
    * [Deploy the server with ansible playbook](#deploy-the-server-with-ansible-playbook)
- [CLI arguments](#cli-arguments)
- [Basic usage](#basic-usage)
- [Clients](#clients)
- [API description](#api-description)
    * [API requests](#api-requests)
    * [API errors](#api-errors)
- [Game Web-Socket messages description](#game-web-socket-messages-description)
    * [Game primitives](#game-primitives)
    * [Game objects](#game-objects)
    * [Output messages](#output-messages)
    * [Input messages](#input-messages)
- [License](#license)

## Game rules

A player controls a snake. The task of the game is to grow the biggest snake. In order to do that players may eat apples, watermelons, smallest snakes and remains of dead snakes of other players. If a snake hits a wall, that snake will die, and the player will start again with new small snake. Once a snake has grown it may eat the smallest snakes.

## Installation

There are many options to get Snake-Server. You can download the server binary, build the server from the source or pull the server docker image. See below.

### Download and install binary

You can download binary from releases page: https://github.com/ivan1993spb/snake-server/releases/latest

Or using curl:

* Setup variables *VERSION*, *PLATFORM* (darwin, linux or windows) and *ARCHITECTURE* (386 or amd64)
* Use curl to download snake-server binary: `curl -sL https://github.com/ivan1993spb/snake-server/releases/download/${VERSION}/snake-server-${VERSION}-${PLATFORM}-${ARCHITECTURE}.tar.gz | tar xvz`

Then:

* Rename binary to `snake-server`: `mv snake-server-${VERSION}-${PLATFORM}-${ARCHITECTURE} snake-server`
* Make binary file executable with `chmod +x snake-server`
* Move snake-server to `/usr/local/bin/`: `mv snake-server /usr/local/bin/`
* Use `snake-server -h` to see the usage information

### Build and install latest server

In order to build Snake-Server you need installed [Go compiler](https://golang.org/) (version 1.6+ is required).

The simpliest way to install the latest Snake-Server is to run go-get command to load the source code from the master branch, build and install snake-server into `${GOPATH}/bin` directory:

```bash
go get -u github.com/ivan1993spb/snake-server
snake-server -h
```

### Pull server image from docker-hub

Firstly, you need installed docker: [use fast installation script](https://get.docker.com/).

See snake-server docker-hub repository: https://hub.docker.com/r/ivan1993spb/snake-server.

Choose image tag from the [tags list](https://hub.docker.com/r/ivan1993spb/snake-server/tags/).

* Use `docker pull ivan1993spb/snake-server` to pull server image from docker-hub
* `docker run --rm --net host --name snake-server ivan1993spb/snake-server` to start server
* `docker run --rm ivan1993spb/snake-server -h` for usage information

Add alias for running snake-server container:

* `alias snake-server="docker run --rm -it --net host --name snake-server ivan1993spb/snake-server:latest"`
* `snake-server --help`

### Build and install server of specific version from source code

You may get Snake-Server source code of specific version by curl command:

* `mkdir -p ${GOPATH}/src/github.com/ivan1993spb/snake-server`
* Download and extract source code `curl -sL https://github.com/ivan1993spb/snake-server/archive/${VERSION}.tar.gz | tar xvz --strip 1 -C ${GOPATH}/src/github.com/ivan1993spb/snake-server`
* `cd ${GOPATH}/src/github.com/ivan1993spb/snake-server`

Then there are three options how to build the server from sources:

* Build the server binary using Go compiler
* Build the server binary using Docker
* Build the Docker image

#### Build the server binary using Go compiler

In order to build Snake-Server you need installed [Go compiler](https://golang.org/) (version 1.6+ is required).

If you have appropriate Go compiler, use:

* `make VERSION=${VERSION} BUILD=custom`
* `make install VERSION=${VERSION} BUILD=custom`

Then:

* Use `snake-server -h` to see the usage information
* `snake-server` to start the server

#### Build the server binary using Docker

If you have Docker, use:

* `make go/build VERSION=${VERSION} BUILD=custom`
* Move snake-server to `/usr/local/bin/`: `mv snake-server /usr/local/bin/`

Then:

* Use `snake-server -h` to see the usage information
* `snake-server` to start the server

#### Build Docker image

In order to build Snake-Server Docker image you need installed [Docker](https://www.docker.com/) (version 17.05+ is required)

Use command `make docker/build` to build image.

Then:

* Use `docker run --rm ivan1993spb/snake-server -h` to see the usage information
* `docker run --rm ivan1993spb/snake-server` to start the server

### Deploy the server with ansible playbook

Snake's playbook repository - [here](https://github.com/ivan1993spb/snake-ansible).

```bash
git clone https://github.com/ivan1993spb/snake-ansible.git
```

## CLI arguments

Use `snake-server --help` for help info.

Arguments:

* `--address` - **string** - address to serve (default: *:8080*). For example: *:8080*, *localhost:7070*
* `--conns-limit` - **int** - opened web-socket connections limit (default: *1000*)
* `--groups-limit` - **int** - game groups limit for server (default: *100*)
* `--enable-broadcast` - **bool** - enable broadcasting API method (default: *false*)
* `--log-json` - **bool** - set this flag to use JSON log format (default: *false*)
* `--log-level` - **string** - set log level: *panic*, *fatal*, *error*, *warning* (*warn*), *info* or *debug* (default: *info*)
* `--seed` - **int** - random seed (default: the number of nanoseconds elapsed since January 1, 1970 UTC)
* `--tls-cert` - **string** - path to certificate file
* `--tls-enable` - **bool** - flag: enable TLS
* `--tls-key` - **string** - path to key file

## Basic usage

Start snake-server:

```bash
snake-server
```

Add a game for 5 players with map width 40 dots and height 30 dots:

```bash
curl -s -X POST -d limit=5 -d width=40 -d height=30 http://localhost:8080/api/games
```

Result:

```json
{
  "id": 1,
  "limit": 5,
  "count": 0,
  "width": 40,
  "height": 30,
  "rate": 0
}
```

Now the web-socket connection handler ready to serve players on url `ws://localhost:8080/ws/games/1`

## Clients

You are welcome to create your own client using described API.

Some samples you can see here:

* Python client repo: https://github.com/ivan1993spb/snake-client
* JS client repo: https://github.com/ivan1993spb/ivan1993spb.github.io

## API description

All API methods provide JSON format. If errors are occurred methods return HTTP statuses and JSON formatted error objects. See [swagger.yml](swagger.yml) for details. Also, see API curl examples below.

### API requests

#### Request `POST /api/games`

Request creates a game and returns JSON game object.

```
curl -s -X POST -d limit=3 -d width=100 -d height=100 http://localhost:8080/api/games | jq
{
  "id": 1,
  "limit": 3,
  "count": 0,
  "width": 100,
  "height": 100,
  "rate": 0
}
```

#### Request `GET /api/games`

Request returns an information about all games on server.

```
curl -s -X GET http://localhost:8080/api/games | jq
{
  "games": [
    {
      "id": 1,
      "limit": 10,
      "count": 0,
      "width": 100,
      "height": 100,
      "rate": 0
    },
    {
      "id": 2,
      "limit": 10,
      "count": 0,
      "width": 100,
      "height": 100,
      "rate": 0
    }
  ],
  "limit": 100,
  "count": 2
}
```

#### Request `GET /api/games/{id}`

Request returns an information about a game by id.

```
curl -s -X GET http://localhost:8080/api/games/1 | jq
{
  "id": 1,
  "limit": 10,
  "count": 0,
  "width": 100,
  "height": 100,
  "rate": 0
}
```

#### Request `DELETE /api/games/{id}`

Request deletes a game by id if there is not players in the game.

```
curl -s -X DELETE http://localhost:8080/api/games/1 | jq
{
  "id": 1
}
```

#### Request `POST /api/games/{id}/broadcast`

Request sends a message to all players in selected game. Returns `true` on success. **Request body size is limited: maximum 128 bytes**

```
curl -s -X POST -d message=text http://localhost:8080/api/games/1/broadcast | jq
{
  "success": true
}
```

If request method is disabled, you will get 404 error. See CLI arguments.

#### Request `GET /api/games/{id}/objects`

Request returns all objects on map of a game with passed identifier.

```
curl -s -X GET http://localhost:8080/api/games/1/objects | jq
{
  "objects": [
    {
      "uuid": "066167c0-38eb-424e-82fc-942ded486a84",
      "dots": [
        [0, 2],
        [1, 2],
        [0, 0],
        [1, 0],
        [1, 1],
        [2, 1]
      ],
      "type": "wall"
    },
    {
      "uuid": "e91944bc-f31f-4b43-8a6c-2189db3734e5",
      "dot": [18, 16],
      "type": "apple"
    },
    {
      "uuid": "680575ca-5ec0-4071-a495-be107b0fd255",
      "dots": [
        [9, 17],
        [10, 17],
        [9, 18],
        [10, 18]
      ],
      "type": "watermelon"
    }
  ]
}
```

#### Request `GET /api/capacity`

Request returns the server instance capacity metric. Capacity is the number of opened web-socket connections divided by the number of allowed connections for the server instance.

```
curl -s -X GET http://localhost:8080/api/capacity | jq
{
  "capacity": 0.02
}
```

#### Request `GET /api/info`

Request returns the common information about the server: author, license, version, build.

```
curl -s -X GET http://localhost:8080/api/info | jq
{
  "author": "Ivan Pushkin",
  "license": "MIT",
  "version": "v4.0.0",
  "build": "85b6b0e"
}
```

#### Request `GET /api/ping`

Request returns pong response from the server.

```
curl -s -X GET localhost:8080/api/ping | jq
{
  "pong": 1
}
```

### API errors

API methods return error status codes (400, 404, 500, etc.) with error description in JSON format: `{"code": error_code , "text": error_text }`. JSON error structure can contains additional fields.

Example:

```
curl -s -X GET http://localhost:8080/api/games/1 -v | jq
*   Trying 127.0.0.1...
* Connected to localhost (127.0.0.1) port 8080 (#0)
> GET /api/games/0 HTTP/1.1
> Host: localhost:8080
> User-Agent: curl/7.47.0
> Accept: */*
>
< HTTP/1.1 404 Not Found
< Server: Snake-Server/v3.1.1-rc (build 85b6b0e)
< Vary: Origin
< Date: Wed, 20 Jun 2018 12:24:44 GMT
< Content-Length: 44
< Content-Type: application/json; charset=utf-8
<
{ [44 bytes data]
* Connection #0 to host localhost left intact
{
  "code": 404,
  "text": "game not found",
  "id": 1
}
```

## Game Web-Socket messages description

The request `ws://localhost:8080/ws/games/1` connects to the game Web-Socket JSON stream by game identificator.

When connection has established, handler:

* Initializes the game session
* Returns the playground size
* Returns all objects on playground
* Creates the snake
* Returns the snake uuid
* Pushes game events and objects to web-socket stream

There are input and output web-socket messages.

### Game primitives

To explain game objects there is JSON primitives:

* Direction: `"north"`, `"west"`, `"south"`, `"east"`
* Dot: `[x, y]`
* Dot list: `[[x, y], [x, y], [x, y], [x, y], [x, y], [x, y]]`
* Rectangle: `[x, y, width, height]`

### Game objects

Game objects:

* Snake:
  ```json
  {
    "type": "snake",
    "uuid": "a065eabe-101f-48ba-8b23-d8d5ded7957c",
    "dots": [[4, 3], [3, 3], [2, 3]]
  }
  ```
* Apple:
  ```json
  {
    "type": "apple",
    "uuid": "a065eabe-101f-48ba-8b23-d8d5ded7957c",
    "dot": [3, 2]
  }
  ```
* Corpse:
  ```json
  {
    "type": "corpse",
    "uuid": "a065eabe-101f-48ba-8b23-d8d5ded7957c",
    "dots": [[3, 2], [3, 1], [3, 0]]
  }
  ```
* Watermelon:
  ```json
  {
    "type": "watermelon",
    "uuid": "a065eabe-101f-48ba-8b23-d8d5ded7957c",
    "dots": [[4, 4], [4, 5], [5, 4], [5, 5]]
  }
  ```
* Wall:
  ```json
  {
    "type": "wall",
    "uuid": "a065eabe-101f-48ba-8b23-d8d5ded7957c",
    "dots": [[4, 2], [2, 1], [2, 3]]
  }
  ```

### Output messages

Output messages - when server sends to the client a data.

Output message structure:

```
{
  "type": output_message_type,
  "payload": output_message_payload
}
```

Output message can be type of:

* *game* - message payload contains a game events. Game events has type and payload: `{"type": game_event_type, "payload": game_event_payload}`. Game events contains information about creation, updation, deletion of game objects on playground
* *player* - message payload contains a player specified info. Player messages has type and payload: `{"type": player_message_type, "payload": player_message_payload}`. Player messages contains user specific game information: user notifications, errors, snake uuid, etc.
* *broadcast* - message payload contains a group broadcast messages. Payload of output message of type *broadcast* contains **string** message

Examples:

```
{
  "type": "player",
  "payload": ...
}
{
  "type": "game",
  "payload": ...
}
{
  "type": "broadcast",
  "payload": ...
}
```

#### Game events

Output message type: *game*

Game events types:

* *error* - a payload contains **string**: error description
* *create* - a payload contains game object that was created
* *delete* - a payload contains game object that was deleted
* *update* - a payload contains game object that was updated
* *checked* - a payload contains game object that was checked by another game object

Examples:

```json
{
  "type": "game",
  "payload": {
    "type": "create",
    "payload": {
      "uuid": "b065eade-101f-48ba-8b23-d8d5ded7957c",
      "dots": [[9, 9], [9, 8], [9, 7]],
      "type": "snake"
    }
  }
}
{
  "type": "game",
  "payload": {
    "type": "update",
    "payload": {
      "uuid": "a4a82fbe-a3d6-4cfa-9e2e-7d7ac1f949b1",
      "dots": [[19, 6], [19, 7], [19, 8]],
      "type": "snake"
    }
  }
}
{
  "type": "game",
  "payload": {
    "type": "checked",
    "payload": {
      "uuid": "110fd923-8167-4475-a9d5-b8cd41a60f9e",
      "dots": [[6, 17], [6, 18], [6, 19], [7, 19], [8, 19], [8, 20], [8, 21]],
      "type": "corpse"
    }
  }
}
{
  "type": "game",
  "payload": {
    "type": "update",
    "payload": {
      "uuid": "110fd923-8167-4475-a9d5-b8cd41a60f9e",
      "dots": [[6, 17], [6, 18], [6, 19], [7, 19], [8, 19], [8, 20], [8, 21]],
      "type": "corpse"
    }
  }
}
```

#### Player messages

Output message type: *player*

Player messages types:

* *size* - a payload contains playground size **object**: `{"width":10,"height":10}`
* *snake* - a payload contains **string**: snake identifier
* *notice* - a payload contains **string**: a notification
* *error* - a payload contains **string**: error description
* *countdown* - a payload contains **int**: number of seconds for countdown
* *objects* - a payload contains a list of all objects on the playground. The message contained objects is necessary to initialize the map on client side

Examples:

```json
{
  "type": "player",
  "payload": {
    "type": "notice",
    "payload": "welcome to snake-server!"
  }
}
{
  "type": "player",
  "payload": {
    "type": "size",
    "payload": {
      "width":255,
      "height":255
    }
  }
}
{
  "type": "player",
  "payload": {
    "type": "objects",
    "payload": [
      {
        "uuid": "e0d5c710-cdc7-43d5-9c4f-5e1e171c5207",
        "dot": [17, 18],
        "type": "apple"
      },
      {
        "uuid": "db7b856c-6f8e-4229-aee6-b90cdc575e0e",
        "dots": [[24, 24], [25, 24], [26, 24]],
        "type": "corpse"
      }
    ]
  }
}
{
  "type": "player",
  "payload": {
    "type": "countdown",
    "payload": 5
  }
}
```

#### Broadcast messages

Output message type: *broadcast*

Payload of output message of type *broadcast* contains **string** - a group notice that sends to all players in game group.

Example:

```json
{
  "type": "broadcast",
  "payload": "hello world!"
}
```

### Input messages

Input messages - when client sends to server an information.

Input message structure:

```
{
  "type": input_message_type,
  "payload": input_message_payload
}
```

Input message types:

* *snake* - when player sends a game command in message payload to control snake
* *broadcast* - when player sends a short phrase or emoji to broadcast it for players in game

**Input message size is limited: maximum 128 bytes**

#### Snake input message

Snake input message contains game command. Game command sets snake direction if it possible.

Accepted commands:

* *north* - sets snake direction to north
* *east* - sets snake direction to east
* *south* - sets snake direction to south
* *west* - sets snake direction to west

Examples:

```json
{
  "type": "snake",
  "payload": "north"
}
{
  "type": "snake",
  "payload": "east"
}
{
  "type": "snake",
  "payload": "south"
}
{
  "type": "snake",
  "payload": "west"
}
```

#### Broadcast input message

A broadcast input message contains a short message to be sent to all players in a game.

Examples:

```json
{
  "type": "broadcast",
  "payload": "hello!"
}
{
  "type": "broadcast",
  "payload": ";)"
}
```

## TODO

* Create more tests
* Create an object for mouse:
  ```json
  {
    "type": "mouse",
    "uuid": "b065eade-101f-48ba-8b23-d8d5ded7957c",
    "dot": [3, 2],
    "dir": "north"
  }
  ```
* Create a core layer to invoke methods from API handlers.
* Create ffjson to API handlers.
* Create CLI flag to set up max limit value of gamers in a game.
* Create a queue of the commands to the snake.
  Create queue in Snake object:
  ```golang
  commandQueue chan snake.Command
  ```
  Collect and implement commands:
  ```json
  ["north", "east", "north", "north", "north", "east"]
  ```
* Create embedded lightweight javascript client
* Create log message with list of addresses to listen and serve:

## License

See [LICENSE](LICENSE).


# Snake-Server [![Build Status](https://travis-ci.org/ivan1993spb/snake-server.svg?branch=master)](https://travis-ci.org/ivan1993spb/snake-server) [![Go Report Card](https://goreportcard.com/badge/github.com/ivan1993spb/snake-server)](https://goreportcard.com/report/github.com/ivan1993spb/snake-server) [![Swagger Validator](https://img.shields.io/swagger/valid/2.0/https/raw.githubusercontent.com/ivan1993spb/snake-server/master/swagger.yml.svg)](https://raw.githubusercontent.com/ivan1993spb/snake-server/master/swagger.yml)

Server for online arcade game - snake.

## Game rules

The player controls snake. The task of the game is to grow the biggest snake. In order to do this players can eat apples, watermelons and the remains of dead snakes of other players. If the snake hits a wall, the snake will die, and the player will start again with small snake.

## Client

* Python client repo: https://github.com/ivan1993spb/snake-client

## Installation

You can download server binary, build server from source or pull server docker image.

### Download and install binary

You can download binary from releases page: https://github.com/ivan1993spb/snake-server/releases

Or using curl:

* Setup variables *VERSION*, *PLATFORM* (darwin, linux or windows) and *ARCHITECTURE* (386 or amd64)
* Use curl to download snake-server binary: `curl -sL https://github.com/ivan1993spb/snake-server/releases/download/${VERSION}/snake-server-${VERSION}-${PLATFORM}-${ARCHITECTURE} -o snake-server`

Then:

* Make binary file executable with `chmod +x snake-server`
* Move snake-server to `/usr/local/bin/`: `mv snake-server /usr/local/bin/`
* Use `snake-server -h` to see usage information

### Build and install from source

In order to build snake-server you need installed [Go compiler](https://golang.org/) (version 1.6+ is required).

Steps to get snake-server:

* `go get -u github.com/ivan1993spb/snake-server` to load source code
* `cd ${GOPATH}/src/github.com/ivan1993spb/snake-server`
* `make build`
* `make install`
* `snake-server` to start server
* Use `snake-server -h` to see usage information

### Pull server image from docker-hub

Firstly you need installed docker: [use fast installation script](https://get.docker.com/)

See snake-server docker-hub repository: https://hub.docker.com/r/ivan1993spb/snake-server

Choose image tag from [tags list](https://hub.docker.com/r/ivan1993spb/snake-server/tags/)

* Use `docker pull ivan1993spb/snake-server` to pull server image from docker-hub
* `docker run --rm --net host --name snake-server ivan1993spb/snake-server` to start server
* `docker run --rm ivan1993spb/snake-server -h` for usage information

## CLI arguments

Use `snake-server --help` for help info.

Arguments:

* `--address` - **string** - address to serve (default: *:8080*). For example: *:8080*, *localhost:7070*
* `--conns-limit` - **int** - opened web-socket connections limit (default: *1000*)
* `--groups-limit` - **int** - game groups limit for server (default: *100*)
* `--log-json` - **bool** - set this flag to use JSON log format (default: *false*)
* `--log-level` - **string** - set log level: *panic*, *fatal*, *error*, *warning* (*warn*), *info* or *debug* (default: *info*)
* `--seed` - **int** - random seed (default: the number of nanoseconds elapsed since January 1, 1970 UTC)
* `--tls-cert` - **string** - path to certificate file
* `--tls-enable` - **bool** - flag: enable TLS
* `--tls-key` - **string** - path to key file

## API description

All API methods provide JSON format. If errors occurred methods return HTTP statuses and JSON formatted error objects. See [swagger.yml](swagger.yml) for details. Also, see API curl examples below.

### Request `POST /api/games`

Request creates game and returns JSON game object.

```
curl -s -X POST -d limit=3 -d width=100 -d height=100 http://localhost:8080/api/games | jq
{
    "id": 0,
    "limit": 3,
    "count": 0,
    "width": 100,
    "height": 100
}
```

### Request `GET /api/games`

Request returns info about all games on server.

```
curl -s -X GET http://localhost:8080/api/games | jq
{
    "games": [
        {
            "id": 1,
            "limit": 10,
            "count": 0,
            "width": 100,
            "height": 100
        },
        {
            "id": 0,
            "limit": 10,
            "count": 0,
            "width": 100,
            "height": 100
        }
    ],
    "limit": 100,
    "count": 2
}
```

### Request `GET /api/games/{id}`

Request returns information about game by id.

```
curl -s -X GET http://localhost:8080/api/games/0 | jq
{
    "id": 0,
    "limit": 10,
    "count": 0,
    "width": 100,
    "height": 100
}
```

### Request `DELETE /api/games/{id}`

Request deletes game by id if there is not players in the game.

```
curl -s -X DELETE http://localhost:8080/api/games/0 | jq
{
    "id": 0
}
```

### Request `GET /api/capacity`

Request returns server capacity. Capacity is the number of opened web-socket connections divided by the number of allowed connections for server instance.

```
curl -s -X GET http://localhost:8080/api/capacity | jq
{
    "capacity": 0.02
}
```

### Request `GET /api/info`

Request returns common info about server: author, license, version, build.

```
curl -s -X GET http://localhost:8080/api/info | jq
{
  "author": "Ivan Pushkin",
  "license": "MIT",
  "version": "v3.1.1-rc",
  "build": "85b6b0e"
}
```

### Request `POST /api/games/{id}/broadcast`

Request sends message to all players in selected game. Returns `true` on success.

```
curl -s -X POST -d message=text http://localhost:8080/api/games/0/broadcast | jq
{
  "success": true
}
```

**Request body size is limited: maximum 128 bytes**

### API errors

API methods returns status codes (400, 404, 500, etc.) with errors in JSON format: `{"code": error_code , "text": error_text }`. JSON error structure can contains additional fields.

Example:

```
curl -s -X GET http://localhost:8080/api/games/0 -v | jq
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
  "id": 0
}
```

## Game Web-Socket messages description

Request `GET /ws/games/{id}` connects to game Web-Socket JSON stream by game identificator.

On connection establishing handler:

* Initializes game session
* Returns playground size
* Returns all objects on playground
* Creates snake
* Returns snake uuid
* Pushes game events and objects in output messages

There are input and output web-socket messages.

### Output messages

Output messages - when server sends to client a data.

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

* *error* - payload contains **string**: error description
* *create* - payload contains game object that was created
* *delete* - payload contains game object that was deleted
* *update* - payload contains game object that was updated
* *checked* - payload contains game object that was checked by another game object

Examples:

```
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

* *size* - payload contains playground size **object**: `{"width":10,"height":10}`
* *snake* - payload contains **string**: snake identifier
* *notice* - payload contains **string**: a notification
* *error* - payload contains **string**: error description
* *countdown* - payload contains **int**: number of seconds for countdown
* *objects* - payload contains list of all objects on playground. The message contained objects is needed to initialize the map on client side

Examples:

```
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

```
{
    "type": "broadcast",
    "payload": "hello world!"
}
```

### Game primitives

Primitives that used to explain game objects:

* Direction: `"north"`, `"west"`, `"south"`, `"east"`
* Dot: `[x, y]`
* Dot list: `[[x, y], [x, y], [x, y], [x, y], [x, y], [x, y]]`

### Game objects

Game objects:

* Apple: `{"type": "apple", "uuid": ... , "dot": [x, y]}`
* Corpse: `{"type": "corpse", "uuid": ... , "dots": [[x, y], [x, y], [x, y]]}`
* Snake: `{"type": "snake", "uuid": ... , "dots": [[x, y], [x, y], [x, y]]}`
* Wall: `{"type": "wall", "uuid": ... , "dots": [[x, y], [x, y], [x, y]]}`
* Watermelon: `{"type": "watermelon", "uuid": ... , "dots": [[x, y], [x, y], [x, y], [x, y]]}`

Objects TODO:

* Mouse: `{"type": "mouse", "uuid": ... , dot: [x, y], "dir": "north"}`
* ...

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

```
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

Broadcast input message contains short message to send to all players in game.

Examples:

```
{
    "type": "broadcast",
    "payload": "hello!"
}
{
    "type": "broadcast",
    "payload": ";)"
}
```

## License

See [LICENSE](LICENSE).

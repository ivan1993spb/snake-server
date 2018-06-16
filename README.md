
# Snake server [![Build Status](https://travis-ci.org/ivan1993spb/snake-server.svg?branch=master)](https://travis-ci.org/ivan1993spb/snake-server) [![Go Report Card](https://goreportcard.com/badge/github.com/ivan1993spb/snake-server)](https://goreportcard.com/report/github.com/ivan1993spb/snake-server)

Server for online arcade game - snake.

## Game rules

The player controls the snake. The task of the game is to grow the biggest snake. In order to do this gamers can eat apples and the remains of dead snakes of other games. If the snake hits the wall, the snake will die, and the player will start again with small snake.

## Client

* Python client repo: https://github.com/ivan1993spb/snake-client

## Installation

You can download server binary, build server from source and pull server docker image.

### Download binary

You can download binary from releases page: https://github.com/ivan1993spb/snake-server/releases

* Setup variables *VERSION*, *PLATFORM* (darwin, linux or windows) and *ARCHITECTURE* (386 or amd64)
* Use curl to download snake-server binary: `curl -sL https://github.com/ivan1993spb/snake-server/releases/download/${VERSION}/snake-server-${VERSION}-${PLATFORM}-${ARCHITECTURE} -o snake-server`
* Make binary executable with `chmod +x snake-server`
* Move snake-server to `/usr/local/bin/`: `mv snake-server /usr/local/bin/`
* Use `snake-server -h` to see usage information

### Install from source

With make install:

* `go get -u github.com/ivan1993spb/snake-server` to load source code
* `cd ${GOPATH}/src/github.com/ivan1993spb/snake-server`
* `make build`
* `make install`

Then:

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
* `--conns-limit` - **int** - open web-socket connections limit (default: *1000*)
* `--groups-limit` - **int** - groups limit for server (default: *100*)
* `--log-json` - **bool** - set this flag to use JSON log format (default: *false*)
* `--log-level` - **string** - set log level: *panic*, *fatal*, *error*, *warning* (*warn*), *info* or *debug* (default: *info*)
* `--seed` - **int** - random seed (default: the number of nanoseconds elapsed since January 1, 1970 UTC)
* `--tls-cert` - **string** - path to certificate file
* `--tls-enable` - **bool** - flag: enable TLS
* `--tls-key` - **string** - path to key file

## API description

API methods provide JSON format.

### Request `POST /games`

Creates game and returns JSON details.

```
curl -s -X POST -d limit=3 -d width=100 -d height=100 http://localhost:8080/games | jq
{
    "id": 0,
    "limit": 3,
    "width": 100,
    "height": 100
}
```

### Request `GET /games`

Returns info about all games on server.

```
curl -s -X GET http://localhost:8080/games | jq
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

### Request `GET /games/{id}`

Returns game information.

```
curl -s -X GET http://localhost:8080/games/0 | jq
{
    "id": 0,
    "limit": 10,
    "count": 0,
    "width": 100,
    "height": 100
}
```

### Request `DELETE /games/{id}`

Deletes game if there is not players.

```
curl -s -X DELETE http://localhost:8080/games/0 | jq
{
    "id": 0
}
```

### Request `GET /capacity`

Returns server capacity. Capacity is the number of opened connections divided by the number of allowed connections for server instance.

```
curl -s -X GET http://localhost:8080/capacity | jq
{
    "capacity": 0.02
}
```

### Request `GET /games/{id}/ws`

Connects to game Web-Socket.

* Returns playground size
* Initialize game session
* Returns all objects on playground
* Creates snake
* Returns snake uuid
* Pushes game events and objects

## Game Web-Socket messages description

There are input and output messages.

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

* *game* - message payload contains a game events. Game events has type and payload: `{"type": game_event_type, "payload": game_event_payload}`. Game events contains information about creation, updation, deletion of objects on playground
* *player* - message payload contains a player info. Player messages has type and payload: `{"type": player_message_type, "payload": player_message_payload}`
* *broadcast* - message payload contains a group broadcast messages. Output message of type *broadcast* is **string**

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
* *objects* - payload contains list of all objects on playground

Examples:

```
{
    "type": "player",
    "payload": {
        "type": "notice",
        "payload": "welcome to snake server!"
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

Objects TODO:

* Watermelon: `{"type": "watermelon", "uuid": ... , "dots": [[x, y], [x, y], [x, y]]}`
* Mouse: `{"type": "mouse", "uuid": ... , dot: [x, y], "dir": "north"}`

### Input messages

Input messages - when client sends to server a game commands.

Input message structure:

```
{
    "type": input_message_type,
    "payload": input_message_payload
}
```

Input message types:

* *snake* - when player sends a game command in message payload to control snake
* *broadcast* - when player sends a short phrase or emoji to broadcast for game group

Accepted game commands:

* *north*
* *east*
* *south*
* *west*

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
{
    "type": "broadcast",
    "payload": "xD"
}
{
    "type": "broadcast",
    "payload": "ok!"
}
{
    "type": "broadcast",
    "payload": "hello!"
}
{
    "type": "broadcast",
    "payload": ";)"
}
```

**Input message size is limited: maximum 128 bytes**

## Game on client side

// TODO: Describe topic

## License

See [LICENSE](LICENSE).

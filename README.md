
# Snake server [![Build Status](https://travis-ci.org/ivan1993spb/snake-server.svg?branch=master)](https://travis-ci.org/ivan1993spb/snake-server) [![Go Report Card](https://goreportcard.com/badge/github.com/ivan1993spb/snake-server)](https://goreportcard.com/report/github.com/ivan1993spb/snake-server)

Server for online arcade game - snake.

## Game rules

// TODO: Write game rules.

## Client

* Python client repo: https://github.com/ivan1993spb/snake-client

## Installation

Download binary from releases page: https://github.com/ivan1993spb/snake-server/releases

Also you can build server from source or pull docker image.

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
* `--groups-limit` - **int** - groups limit for server (default: *100*)
* `--conns-limit` - **int** - open web-socket connections limit (default: *1000*)
* `--seed` - **int** - random seed (default: the number of nanoseconds elapsed since January 1, 1970 UTC)
* `--log-json` - **bool** - set this flag to use JSON log format (default: *false*)
* `--log-level` - **string** - set log level: *panic*, *fatal*, *error*, *warning* (*warn*), *info* or *debug* (default: *info*)

## API description

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

### Request `GET /capacity`

Returns server capacity.

```
curl -X GET http://localhost:8080/capacity
{"capacity":0.02}
```

### Request `GET /games/{id}/ws`

Connects to game Web-Socket.

* Returns playground size
* Initialize game session
* Returns all objects on playground
* Creates snake
* Returns snake id
* Pushes game events and objects

## Game Web-Socket messages description

There are input and output messages.

### Output messages

Output messages - when server sends to client a data.

Output message structure:

```
{"type": output_message_type, "payload": output_message_payload}
```

Output message can be type of:

* *game* - message payload contains a game events. Game events has type and payload: `{"type": game_event_type, "payload": game_event_payload}`. Game events contains information about creation, updation, deletion of objects on playground
* *player* - message payload contains a player info. Player messages has type and payload: `{"type": player_message_type, "payload": player_message_payload}`
* *broadcast* - message payload contains a group broadcast messages. Output message of type *broadcast* is **string**

Examples:

```
{"type":"player","payload": ... }
{"type":"game","payload": ... }
{"type":"broadcast","payload": ... }
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
{"type":"game","payload":{"type":"update","payload":{"id":"c4203c7b00","dots":[[233,236],[233,235],[233,234]]}}}
{"type":"game","payload":{"type":"update","payload":{"id":"c420425240","dots":[[120,130],[120,129],[120,128]]}}}
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
{"type":"player","payload":{"type":"notice","payload":"welcome to snake server!"}}
{"type":"player","payload":{"type":"size","payload":{"width":255,"height":255}}}
{"type":"player","payload":{"type":"objects","payload":[{"id":"c4203c7b00","dots":[[233,235],[233,234],[233,233]]},{"id":"c420425240","dots":[[120,129],[120,128],[120,127]]},{"id":"c420311800","dots":[[60,166],[61,166],[62,166]]},{"id":"c4203c7c00","dots":[[40,46],[41,46],[42,46]]}]}}
{"type":"player","payload":{"type":"countdown","payload":5}}
```

#### Broadcast messages

Output message type: *broadcast*

Payload of output message of type *broadcast* contains **string** - a group notice that sends to all players in game group. 

Examples:

```
{"type":"broadcast","payload":"hello world!"}
```

### Game primitives

Primitives that used to explain game objects:

* Direction: `"north"`, `"west"`, `"south"`, `"east"`
* Dot: `[x, y]`
* Dot list: `[[x, y], [x, y], [x, y], [x, y], [x, y], [x, y]]`

### Game objects

Game objects:

* Apple: `{"type": "apple", "id": 1, "dot": [x, y]}`
* Corpse: `{"type": "corpse", "id": 2, "dots": [[x, y], [x, y], [x, y]]}`
* Mouse: `{"type": "mouse", "id": 3, dot: [x, y], "dir": "n"}`
* Snake: `{"type": "snake", "id": 4, "dots": [[x, y], [x, y], [x, y]]}`
* Wall: `{"type": "wall", "id": 5, "dots": [[x, y], [x, y], [x, y]]}`
* Watermelon: `{"type": "watermelon", "id": 6, "dots": [[x, y], [x, y], [x, y]]}`

// TODO: Fix this topic

### Input messages

Input messages - when client sends to server a game commands.

Input message structure:

```
{"type": input_message_type, "payload": input_message_payload}
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
{"type":"snake","payload":"north"}
{"type":"snake","payload":"east"}
{"type":"snake","payload":"south"}
{"type":"snake","payload":"west"}
{"type":"broadcast","payload":"xD"}
{"type":"broadcast","payload":"ok!"}
{"type":"broadcast","payload":"hello!"}
{"type":"broadcast","payload":";)"}
```

**Input message size is limited: maximum 128 bytes**

## Game on client side

// TODO: Describe topic

## License

See [LICENSE](LICENSE).

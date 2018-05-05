
# Snake server

Server for online arcade game - snake.

// TODO: Create screen shot

## CLI arguments

// TODO: Create arguments description

## API Description

API methods provide JSON format.

### `POST /game/`

Creates game.

### `DELETE /game/{id}`

Deletes game if there is not players.

### `GET /game/{id}`

Connects to game WebSocket.

* return playground size : width and height
* return room_id and player_id
* initialize gamer objects and session
* return all objects on playground
* push events and objects from game

Primitives:

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

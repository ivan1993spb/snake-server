
# Snake server

Server for online arcade game snake.

```
+--------------------------------------------------------------+
|                                                              |
|     0                                                        |
|                   o                                          |
|                   x                                          |
|            xxxxx  x                                          |
|            x   xxxx                     0                    |
|          xxx                                                 |
|          x                0             0                    |
|          x                                                   |
|                                                              |
|                                        xo   0                |
|                                     xxxx                     |
|               0                     x                        |
|  xxo                          xxxxxxx                        |
|  x                     xxxxxxxx                              |
|  x   0                 x                                     |
|xxx                     x                                 xxxx|
|                                                              |
|                                                              |
+--------------------------------------------------------------+
```

## API Description

API methods provide JSON format.

### `GET /info`

* returns information about rooms
* gamer limit in room
* count of gamer in room

Response:

```json
[
    {
        "id": 1,
        "available": 4,
        "players": 11
    },
    {
        "id": 2,
        "available": 0,
        "players": 15
    }
]
```

### `/game`

game WebSocket handler - json stream

* verify token
* return playground size : width and height
* return room_id and player_id
* initialize gamer objects and session
* return all objects on playground
* push events and objects from game

Objects:

* Dot: `[x, y]`
* Location: `[[x, y], [x, y], [x, y], [x, y], [x, y], [x, y]]`
* Direction: `"n"`, `"w"`, `"s"`, `"e"`

Game objects:

* Apple: `{"type": "apple", "id": 1, "dot": [x, y]}`
* Corpse: `{"type": "corpse", "id": 2, "dots": [[x, y], [x, y], [x, y]]}`
* Mouse: `{"type": "mouse", "id": 3, dot: [x, y], "dir": "n"}`
* Snake: `{"type": "snake", "id": 3, "dots": [[x, y], [x, y], [x, y]]}`
* Wall: `{"type": "wall", "id": 4, "dots": [[x, y], [x, y], [x, y]]}`
* Watermelon: `{"type": "watermelon", "id": 4, "dots": [[x, y], [x, y], [x, y]]}`

Message types:

* Object: `{"type": "object", "object": {}}`
* Error: `{"type": "error", "message": "text"}`
* Notice: `{}`
* Event: `{}`


# Game Web-Socket messages

`ws://localhost:8080/ws/games/1` connects a client to the game's web-socket JSON stream.

When connection has been established, the server:

* Initializes a game session
* Returns the map size
* Returns all objects in the game
* Creates a snake
* Returns an identifier of the snake
* Starts pushing updates into the stream

## Game primitives

There are a few game primitives:

* Direction: `"north"`, `"west"`, `"south"`, `"east"`
* Dot: `[x, y]`
* Dot list: `[[x, y], [x, y], [x, y], [x, y], [x, y], [x, y]]`
* Rectangle: `[x, y, width, height]`

## Game objects

Game objects:

* Snake:
  ```json
  {
    "type": "snake",
    "id": 12,
    "dots": [[4, 3], [3, 3], [2, 3]]
  }
  ```
* Apple:
  ```json
  {
    "type": "apple",
    "id": 123,
    "dot": [3, 2]
  }
  ```
* Corpse:
  ```json
  {
    "type": "corpse",
    "id": 142,
    "dots": [[3, 2], [3, 1], [3, 0]]
  }
  ```
* Mouse:
  ```json
  {
    "type": "mouse",
    "id": 123,
    "dot": [3, 2],
    "direction": "south"
  }
  ```
* Watermelon:
  ```json
  {
    "type": "watermelon",
    "id": 123,
    "dots": [[4, 4], [4, 5], [5, 4], [5, 5]]
  }
  ```
* Wall:
  ```json
  {
    "type": "wall",
    "id": 351,
    "dots": [[4, 2], [2, 1], [2, 3]]
  }
  ```

## Game messages 

There are *input* and *output* game messages.

### Output messages

Output messages are sent by server to client.

Example:

```
{
  "type": <output_message_type>,
  "payload": <output_message_payload>
}
```

Types:

* *game* - contains game events. Game events have a type and a payload:
  
  ```
  {
    "type": "game",
    "payload": {
      "type": <game_event_type>,
      "payload": <game_event_payload>
    }
  }
  ```
  
  Game events contain information about creating, updating, deleting of game objects on the map.

* *player* - contains player specific information. Player messages have a type and a payload:

  ```
  {
    "type": "player",
    "payload": {
      "type": <player_message_type>,
      "payload": <player_message_payload>
    }
  }
  ```

  Player messages contain user specific information such as user notifications, errors, snake identifiers, etc.

* *broadcast* - contains broadcasted messages (type **string**):

  ```json
  {
    "type": "broadcast",
    "payload": "Surprise!"
  }
  ```

#### Game events

Output message type: *game*

Game event types:

* *error* - contains description (**string**) of the error

* *create* - contains the created object:

  ```json
  {
    "type": "game",
    "payload": {
      "type": "create",
      "payload": {
        "id": 41,
        "dots": [[9, 9], [9, 8], [9, 7]],
        "type": "snake"
      }
    }
  }
  ```

* *delete* - contains the deleted object

* *update* - contains the updated object
  + When the snake moves:
    ```json
    {
      "type": "game",
      "payload": {
        "type": "update",
        "payload": {
          "id": 123,
          "dots": [[19, 6], [19, 7], [19, 8]],
          "type": "snake"
        }
      }
    }
    ```
  + An update of a corpse:
    ```json
    {
      "type": "game",
      "payload": {
        "type": "update",
        "payload": {
          "id": 142,
          "dots": [[6, 17], [6, 18], [6, 19], [7, 19], [8, 19], [8, 20], [8, 21]],
          "type": "corpse"
        }
      }
    }
    ```

* ~~*checked* - contains an object which was checked by another game object (**deprecated**)~~

#### Player messages

Output message type: *player*

Player's messages types:

* *size* - contains size of the map (**object**):
  ```json
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
  ```

* *snake* - contains an **integer**: a snake identifier

* *notice* - contains a notification (**string**)
  ```json
  {
    "type": "player",
    "payload": {
      "type": "notice",
      "payload": "welcome to snake-server!"
    }
  }
  ```

* *error* - contains a **string**: error description
  ```json
  {
    "type": "player",
    "payload": {
      "type": "error",
      "payload": "something went wrong!"
    }
  }
  ```

* *countdown* - contains an **integer**: the number of seconds to wait
  ```json
  {
    "type": "player",
    "payload": {
      "type": "countdown",
      "payload": 5
    }
  }
  ```

* *objects* - contains a list of all objects in the game to initialize the map on the client side
  ```json
  {
    "type": "player",
    "payload": {
      "type": "objects",
      "payload": [
        {
          "id": 21,
          "dot": [17, 18],
          "type": "apple"
        },
        {
          "id": 63,
          "dots": [[24, 24], [25, 24], [26, 24]],
          "type": "corpse"
        }
      ]
    }
  }
  ```

#### Broadcast messages

Output message type: *broadcast*

Contains a **string** - a group notification to all players in the game.

Example:

```json
{
  "type": "broadcast",
  "payload": "hello world!"
}
```

### Input messages

Input messages are sent by client to server.

Input message structure:

```
{
  "type": <input_message_type>,
  "payload": <input_message_payload>
}
```

**Input message size is limited: maximum 128 bytes**

Input message types:

* *snake* - snake commands
* *broadcast* - short phrases or emojis to be broadcasted in the game

#### Snake input message

A *snake* input message contains a command which sets snake's movement direction.

Accepted commands:

* *north* - to the north
  ```json
  {
    "type": "snake",
    "payload": "north"
  }
  ```
* *east* - to the east
  ```json
  {
    "type": "snake",
    "payload": "east"
  }
  ```
* *south* - to the south
  ```json
  {
    "type": "snake",
    "payload": "south"
  }
  ```
* *west* - to the west
  ```json
  {
    "type": "snake",
    "payload": "west"
  }
  ```

#### Broadcast input message

A broadcast input message contains a short message to all players in the game.

Examples:

* Text: *hello!*
  ```json
  {
    "type": "broadcast",
    "payload": "hello!"
  }
  ```
* Smile: *;)*
  ```json
  {
    "type": "broadcast",
    "payload": ";)"
  }
  ```

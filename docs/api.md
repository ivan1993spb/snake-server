
# API description

All API methods provide JSON formated responses.

If errors occurre, methods return HTTP statuses and JSON formatted objects.

OpenAPI specification `openapi.yaml` is served at `/openapi.yaml`.

## API requests

It is recommended to use the header `X-Snake-Client` to specify the client's name, version and build hash. For example:

```
X-Snake-Client: SnakeLightweightClient/v0.3.2 (build 8554f6b)
```

API methods:

* **`POST /api/games`**

  Creates a game and returns a JSON game object.

  ```
  curl -s -X POST -d limit=3 -d width=100 -d height=100 -d enable_walls=true http://localhost:8080/api/games | jq
  {
    "id": 1,
    "limit": 3,
    "count": 0,
    "width": 100,
    "height": 100,
    "rate": 0
  }
  ```

  `enable_walls` is an optional parameter, the default value is `true`

* **`GET /api/games`**

  Returns information about all games on the server.

  Optional **query string** params:

  + `limit` - **integer** - limit the number of games in response
  + `sorting` - **string** - set a sorting rule. Could be either `smart` or `random`. The default value is `random`

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

* **`GET /api/games/{id}`**

  Returns information about a game by id.

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

* **`DELETE /api/games/{id}`**

  Deletes a game by id if there are no players in the game.

  ```
  curl -s -X DELETE http://localhost:8080/api/games/1 | jq
  {
    "id": 1
  }
  ```

* ~~**`POST /api/games/{id}/broadcast`**~~

  ***DEPRECATED***

  Sends a message to all players in a selected game. Returns `true` on success.
  
  **Request body size is limited: maximum 128 bytes**

  ```
  curl -s -X POST -d message=text http://localhost:8080/api/games/1/broadcast | jq
  {
    "success": true
  }
  ```

  If request method is disabled, you will receive error 404.

  It is a good idea to never ever enable this method. By default it is disabled.

* **`GET /api/games/{id}/objects`**

  Returns all objects and map properties of a game.

  ```
  curl -s -X GET http://localhost:8080/api/games/1/objects | jq
  {
    "objects": [
      {
        "id": 99,
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
        "id": 124,
        "dot": [18, 16],
        "type": "apple"
      },
      {
        "id": 312,
        "dots": [
          [9, 17],
          [10, 17],
          [9, 18],
          [10, 18]
        ],
        "type": "watermelon"
      }
    ],
    "map": {
      "width": 120,
      "height": 75
    }
  }
  ```

* **`GET /api/capacity`**

  Returns capacity of the server. Capacity is the number of opened web-socket
  connections divided by the number of allowed connections for the instance.

  ```
  curl -s -X GET http://localhost:8080/api/capacity | jq
  {
    "capacity": 0.02
  }
  ```

* **`GET /api/info`**

  Returns information about the server: author, license, version, build.

  ```
  curl -s -X GET http://localhost:8080/api/info | jq
  {
    "author": "Ivan Pushkin",
    "license": "MIT",
    "version": "v4.0.0",
    "build": "85b6b0e"
  }
  ```

* **`GET /api/ping`**

  Returns a pong response from a server.

  ```
  curl -s -X GET http://localhost:8080/api/ping | jq
  {
    "pong": 1
  }
  ```

## API errors

API methods return error status codes (400, 404, 500, etc.) with descriptions in JSON format:

```
{
  "code": <error_code>,
  "text": <error_text>
}
```

JSON error structure may contain additional fields.

Example:

```
curl -s -X GET http://localhost:8080/api/games/1 -v | jq
{
  "code": 404,
  "text": "game not found",
  "id": 1
}
```

package main

// TODO: Move main.go to cmd/snake-server/main.go
//       Rename the package of embed.go to snakeserver.

import _ "embed"

//go:embed openapi.yaml
var OpenAPISpec []byte

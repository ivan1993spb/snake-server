package server

import (
	"fmt"
	"net/http"

	"github.com/go-chi/render"

	"github.com/ivan1993spb/snake-server/pkg/openapi/models"
	"github.com/ivan1993spb/snake-server/pkg/openapi/server"
)

type Handler struct {
}

// Server capacity
// (GET /capacity)
func (h *Handler) GetCapacity(w http.ResponseWriter, r *http.Request) {
	panic("not implemented") // TODO: Implement
}

// Get a list of games
// (GET /games)
func (h *Handler) GetGames(w http.ResponseWriter, r *http.Request, params models.GetGamesParams) {
	panic("not implemented") // TODO: Implement
}

// Create a new game
// (POST /games)
func (h *Handler) PostGames(w http.ResponseWriter, r *http.Request) {
	input := &models.PostGamesFormdataBody{}
	if err := render.DecodeForm(r.Body, input); err != nil {
		fmt.Println(err)
	}

	fmt.Println(input)

	// TODO: Implement
}

// Delete a game
// (DELETE /games/{id})
func (h *Handler) DeleteGamesId(w http.ResponseWriter, r *http.Request, id models.GameID) {
	panic("not implemented") // TODO: Implement
}

// Get information about a game
// (GET /games/{id})
func (h *Handler) GetGamesId(w http.ResponseWriter, r *http.Request, id models.GameID) {
	panic("not implemented") // TODO: Implement
}

// Broadcast a message
// (POST /games/{id}/broadcast)
func (h *Handler) PostGamesIdBroadcast(w http.ResponseWriter, r *http.Request, id models.GameID) {
	panic("not implemented") // TODO: Implement
}

// A list of objects on the map
// (GET /games/{id}/objects)
func (h *Handler) GetGamesIdObjects(w http.ResponseWriter, r *http.Request, id models.GameID) {
	panic("not implemented") // TODO: Implement
}

// Information about the server
// (GET /info)
func (h *Handler) GetInfo(w http.ResponseWriter, r *http.Request) {
	panic("not implemented") // TODO: Implement
}

// Ping-pong requesting
// (GET /ping)
func (h *Handler) GetPing(w http.ResponseWriter, r *http.Request) {
	panic("not implemented") // TODO: Implement
}

var _ server.ServerInterface = (*Handler)(nil)

package ws

import (
	"github.com/fykyby/chat-app-backend/database"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/jwtauth/v5"
)

const MESSAGE_ERROR_GENERIC = "An unknown error has occurred"

type WebSocketHandler struct {
	DB        *database.Queries
	TokenAuth *jwtauth.JWTAuth
}

var db *database.Queries
var tokenAuth *jwtauth.JWTAuth

func (h *WebSocketHandler) Handler(r chi.Router) {
	db = h.DB
	tokenAuth = h.TokenAuth

	r.Get("/chats/{id}", getChatWs)
}

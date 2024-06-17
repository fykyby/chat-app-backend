package handler

import (
	"github.com/fykyby/chat-app-backend/api/route"
	"github.com/fykyby/chat-app-backend/api/route/ws/chatws"
	"github.com/fykyby/chat-app-backend/internal/database"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/jwtauth/v5"
)

type ApiHandler struct {
	DB        *database.Queries
	TokenAuth *jwtauth.JWTAuth
}

func (h *ApiHandler) Handler(r chi.Router) {

	registerHandler := route.RegisterHandler{
		DB: h.DB,
	}
	r.Post("/register", registerHandler.Register)

	logInHandler := route.LogInHandler{
		DB:        h.DB,
		TokenAuth: h.TokenAuth,
	}
	r.Post("/login", logInHandler.LogIn)

	logOutHandler := route.LogOutHandler{}
	r.Post("/logout", logOutHandler.LogOut)

	chatHandler := route.ChatHandler{
		DB:        h.DB,
		TokenAuth: h.TokenAuth,
	}
	r.Post("/chats", chatHandler.NewChat)
	r.Get("/chats/{id}", chatHandler.GetChat)

	// ws
	chatWsHandler := chatws.ChatWsHandler{
		DB:    h.DB,
		Rooms: make(map[int32]*chatws.Room),
	}
	r.Get("/chats/{id}", chatWsHandler.ConnectToChatWs)
}
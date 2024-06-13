package handler

import (
	"log"
	"net/http"

	"github.com/fykyby/chat-app-backend/database"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/jwtauth/v5"
)

const MESSAGE_ERROR_GENERIC = "An unknown error has occurred"

type ApiHandler struct {
	DB        *database.Queries
	TokenAuth *jwtauth.JWTAuth
}

var db *database.Queries
var tokenAuth *jwtauth.JWTAuth

func (h *ApiHandler) Handler(r chi.Router) {
	db = h.DB
	tokenAuth = h.TokenAuth

	r.Post("/register", postRegister)
	r.Post("/login", postLogIn)
	r.Post("/logout", postLogOut)

	r.Get("/test", func(w http.ResponseWriter, r *http.Request) {
		result, err := db.GetMessagesPage(r.Context(), database.GetMessagesPageParams{
			ChatID: 1,
			Limit:  20,
			Offset: 0,
		})

		if err != nil {
			log.Println(err)
			// http.Error(w, MESSAGE_ERROR_GENERIC, http.StatusInternalServerError)
			SendResponse(w, http.StatusInternalServerError, MESSAGE_ERROR_GENERIC, nil)
			return
		}

		SendResponse(w, http.StatusOK, "Success", result)
	})
}

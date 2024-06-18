package route

import (
	"encoding/json"
	"log"
	"net/http"
	"net/url"
	"os"

	"github.com/fykyby/chat-app-backend/api"
	"github.com/fykyby/chat-app-backend/internal/auth"
	"github.com/fykyby/chat-app-backend/internal/database"
	"github.com/fykyby/chat-app-backend/internal/status"
)

type RegisterHandler struct {
	DB *database.Queries
}

type registerRequest struct {
	Email           string `json:"email"`
	Name            string `json:"name"`
	Password        string `json:"password"`
	PasswordConfirm string `json:"passwordConfirm"`
}

func (h *RegisterHandler) Register(w http.ResponseWriter, r *http.Request) {
	var req registerRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil || req.Email == "" || req.Name == "" || req.Password == "" {
		log.Println("Error decoding request")
		api.SendResponse(w, http.StatusBadRequest, status.MESSAGE_ERROR_GENERIC, nil)
		return
	}

	if req.Password != req.PasswordConfirm {
		log.Println("Passwords do not match")
		api.SendResponse(w, http.StatusBadRequest, status.MESSAGE_REGISTER_PASSWORDS_DONT_MATCH, nil)
		return
	}

	passwordHash, err := auth.HashPassword(req.Password)
	if err != nil {
		log.Println(err)
		api.SendResponse(w, http.StatusInternalServerError, status.MESSAGE_ERROR_GENERIC, nil)
		return
	}

	user, err := h.DB.CreateUser(r.Context(), database.CreateUserParams{
		Email:    req.Email,
		Name:     req.Name,
		Password: passwordHash,
		Avatar:   os.Getenv("AVATAR_API_URL") + url.QueryEscape(req.Name),
	})
	if err != nil {
		log.Println(err)
		api.SendResponse(w, http.StatusConflict, status.MESSAGE_REGISTER_USER_ALREADY_EXISTS, nil)
		return
	}

	api.SendResponse(w, http.StatusOK, status.MESSAGE_REGISTER_SUCCESS, map[string]interface{}{
		"id":     user.ID,
		"name":   user.Name,
		"email":  user.Email,
		"avatar": user.Avatar,
	})
}

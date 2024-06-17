package route

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/fykyby/chat-app-backend/api"
	"github.com/fykyby/chat-app-backend/internal/auth"
	"github.com/fykyby/chat-app-backend/internal/database"
	"github.com/fykyby/chat-app-backend/internal/status"
	"github.com/jackc/pgx/v5/pgtype"
)

type RegisterHandler struct {
	DB *database.Queries
}

type postRegisterRequest struct {
	Email           string `json:"email"`
	Name            string `json:"name"`
	Password        string `json:"password"`
	PasswordConfirm string `json:"passwordConfirm"`
}

func (h *RegisterHandler) Register(w http.ResponseWriter, r *http.Request) {
	var req postRegisterRequest
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
		Avatar:   pgtype.Text{},
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

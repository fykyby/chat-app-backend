package handler

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/fykyby/chat-app-backend/auth"
	"github.com/fykyby/chat-app-backend/database"
	"github.com/jackc/pgx/v5/pgtype"
)

const MESSAGE_REGISTER_PASSWORDS_DONT_MATCH = "Passwords do not match"
const MESSAGE_REGISTER_USER_ALREADY_EXISTS = "User already exists"
const MESSAGE_REGISTER_SUCCESS = "Registration successful"

type postRegisterRequest struct {
	Email           string `json:"email"`
	Name            string `json:"name"`
	Password        string `json:"password"`
	PasswordConfirm string `json:"passwordConfirm"`
}

func postRegister(w http.ResponseWriter, r *http.Request) {
	var req postRegisterRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil || req.Email == "" || req.Name == "" || req.Password == "" {
		log.Println("Error decoding request")
		SendResponse(w, http.StatusBadRequest, MESSAGE_ERROR_GENERIC, nil)
		return
	}

	if req.Password != req.PasswordConfirm {
		log.Println("Passwords do not match")
		SendResponse(w, http.StatusBadRequest, MESSAGE_REGISTER_PASSWORDS_DONT_MATCH, nil)
		return
	}

	passwordHash, err := auth.HashPassword(req.Password)
	if err != nil {
		log.Println(err)
		SendResponse(w, http.StatusInternalServerError, MESSAGE_ERROR_GENERIC, nil)
		return
	}

	user, err := db.CreateUser(r.Context(), database.CreateUserParams{
		Email:    req.Email,
		Name:     req.Name,
		Password: passwordHash,
		Avatar:   pgtype.Text{},
	})
	if err != nil {
		log.Println(err)
		SendResponse(w, http.StatusConflict, MESSAGE_REGISTER_USER_ALREADY_EXISTS, nil)
		return
	}

	SendResponse(w, http.StatusOK, MESSAGE_REGISTER_SUCCESS, map[string]interface{}{
		"id":     user.ID,
		"name":   user.Name,
		"email":  user.Email,
		"avatar": user.Avatar,
	})
}

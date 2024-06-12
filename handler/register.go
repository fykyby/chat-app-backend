package handler

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/fykyby/go-web-template/auth"
	"github.com/fykyby/go-web-template/database"
)

const MESSAGE_REGISTER_PASSWORDS_DONT_MATCH = "Passwords do not match"
const MESSAGE_REGISTER_USER_ALREADY_EXISTS = "User already exists"
const MESSAGE_REGISTER_SUCCESS = "Registration successful"

type postRegisterRequest struct {
	Email           string `json:"email"`
	Password        string `json:"password"`
	PasswordConfirm string `json:"password_confirm"`
}

func postRegister(w http.ResponseWriter, r *http.Request) {
	var req postRegisterRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil || req.Email == "" || req.Password == "" {
		sendResponse(w, http.StatusBadRequest, MESSAGE_ERROR_GENERIC, nil)
		return
	}

	if req.Password != req.PasswordConfirm {
		sendResponse(w, http.StatusBadRequest, MESSAGE_REGISTER_PASSWORDS_DONT_MATCH, nil)
		return
	}

	passwordHash, err := auth.HashPassword(req.Password)
	if err != nil {
		log.Println(err)
		sendResponse(w, http.StatusInternalServerError, MESSAGE_ERROR_GENERIC, nil)
		return
	}

	user, err := db.CreateUser(r.Context(), database.CreateUserParams{
		Email:    req.Email,
		Password: passwordHash,
	})
	if err != nil {
		log.Println(err)
		sendResponse(w, http.StatusConflict, MESSAGE_REGISTER_USER_ALREADY_EXISTS, nil)
		return
	}

	sendResponse(w, http.StatusOK, MESSAGE_REGISTER_SUCCESS, map[string]interface{}{
		"id":    user.ID,
		"email": user.Email,
	})
}

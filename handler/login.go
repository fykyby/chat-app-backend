package handler

import (
	"encoding/json"
	"net/http"

	"github.com/fykyby/chat-app-backend/auth"
)

const MESSAGE_LOGIN_WRONG_CREDENTIALS = "Wrong email or password"
const MESSAGE_LOGIN_SUCCESS = "Login successful"

type postLogInRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func postLogIn(w http.ResponseWriter, r *http.Request) {
	var req postLogInRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil || req.Email == "" || req.Password == "" {
		sendResponse(w, http.StatusBadRequest, MESSAGE_ERROR_GENERIC, nil)
		return
	}

	user, err := db.GetUserByEmail(r.Context(), req.Email)
	if err != nil {
		sendResponse(w, http.StatusUnauthorized, MESSAGE_LOGIN_WRONG_CREDENTIALS, nil)
		return
	}

	passwordsMatch := auth.CheckPasswordHash(req.Password, user.Password)
	if !passwordsMatch {
		sendResponse(w, http.StatusUnauthorized, MESSAGE_LOGIN_WRONG_CREDENTIALS, nil)
		return
	}

	_, tokenString, _ := tokenAuth.Encode(map[string]interface{}{"id": user.ID, "email": user.Email})
	http.SetCookie(w, &http.Cookie{
		Name:     "jwt",
		Value:    tokenString,
		Secure:   true,
		MaxAge:   2592000,
		SameSite: http.SameSiteNoneMode,
		Path:     "/",
		HttpOnly: true,
	})
	w.Header().Set("jwt", tokenString)

	sendResponse(w, http.StatusOK, MESSAGE_LOGIN_SUCCESS, map[string]interface{}{
		"id":    user.ID,
		"email": user.Email,
	})
}

package handler

import (
	"encoding/json"
	"log"
	"net/http"
	"net/url"

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
		log.Println("Error decoding request: ", err)
		sendResponse(w, http.StatusBadRequest, MESSAGE_ERROR_GENERIC, nil)
		return
	}

	user, err := db.GetUserByEmail(r.Context(), req.Email)
	if err != nil {
		log.Println(err)
		sendResponse(w, http.StatusUnauthorized, MESSAGE_LOGIN_WRONG_CREDENTIALS, nil)
		return
	}

	passwordsMatch := auth.CheckPasswordHash(req.Password, user.Password)
	if !passwordsMatch {
		log.Println("Wrong password")
		sendResponse(w, http.StatusUnauthorized, MESSAGE_LOGIN_WRONG_CREDENTIALS, nil)
		return
	}

	_, tokenString, _ := tokenAuth.Encode(map[string]interface{}{
		"id":    user.ID,
		"name":  user.Name,
		"email": user.Email,
	})
	http.SetCookie(w, &http.Cookie{
		Name:     "jwt",
		Value:    tokenString,
		Secure:   true,
		HttpOnly: true,
		MaxAge:   2592000,
		SameSite: http.SameSiteNoneMode,
		Path:     "/",
	})
	w.Header().Set("jwt", tokenString)

	userMap := map[string]interface{}{
		"id":     user.ID,
		"name":   user.Name,
		"email":  user.Email,
		"avatar": user.Avatar,
	}

	userJson, err := json.Marshal(userMap)
	if err != nil {
		log.Println(err)
		sendResponse(w, http.StatusInternalServerError, MESSAGE_ERROR_GENERIC, nil)
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:     "user",
		Value:    url.QueryEscape(string(userJson)),
		Secure:   true,
		HttpOnly: false,
		MaxAge:   2592000,
		SameSite: http.SameSiteNoneMode,
		Path:     "/",
	})

	sendResponse(w, http.StatusOK, MESSAGE_LOGIN_SUCCESS, userMap)
}

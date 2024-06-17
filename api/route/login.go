package route

import (
	"encoding/json"
	"log"
	"net/http"
	"net/url"

	"github.com/fykyby/chat-app-backend/api"
	"github.com/fykyby/chat-app-backend/internal/auth"
	"github.com/fykyby/chat-app-backend/internal/database"
	"github.com/fykyby/chat-app-backend/internal/status"
	"github.com/go-chi/jwtauth/v5"
)

type LogInHandler struct {
	DB        *database.Queries
	TokenAuth *jwtauth.JWTAuth
}

type postLogInRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (h *LogInHandler) LogIn(w http.ResponseWriter, r *http.Request) {
	var req postLogInRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil || req.Email == "" || req.Password == "" {
		log.Println("Error decoding request: ", err)
		api.SendResponse(w, http.StatusBadRequest, status.MESSAGE_ERROR_GENERIC, nil)
		return
	}

	user, err := h.DB.GetUserByEmail(r.Context(), req.Email)
	if err != nil {
		log.Println(err)
		api.SendResponse(w, http.StatusUnauthorized, status.MESSAGE_LOGIN_WRONG_CREDENTIALS, nil)
		return
	}

	passwordsMatch := auth.CheckPasswordHash(req.Password, user.Password)
	if !passwordsMatch {
		log.Println("Wrong password")
		api.SendResponse(w, http.StatusUnauthorized, status.MESSAGE_LOGIN_WRONG_CREDENTIALS, nil)
		return
	}

	_, tokenString, _ := h.TokenAuth.Encode(map[string]interface{}{
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
		api.SendResponse(w, http.StatusInternalServerError, status.MESSAGE_ERROR_GENERIC, nil)
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

	api.SendResponse(w, http.StatusOK, status.MESSAGE_LOGIN_SUCCESS, userMap)
}

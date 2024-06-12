package handler

import (
	"net/http"
	"time"
)

const MESSAGE_LOGOUT_SUCCESS = "Logout successful"

func postLogOut(w http.ResponseWriter, r *http.Request) {
	http.SetCookie(w, &http.Cookie{
		Name:     "jwt",
		Value:    "",
		Expires:  time.Unix(0, 0),
		Path:     "/",
		HttpOnly: true,
	})
	http.SetCookie(w, &http.Cookie{
		Name:     "user",
		Value:    "",
		Expires:  time.Unix(0, 0),
		Path:     "/",
		HttpOnly: false,
	})
	sendResponse(w, http.StatusOK, MESSAGE_LOGOUT_SUCCESS, nil)
}

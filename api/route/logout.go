package route

import (
	"net/http"
	"time"

	"github.com/fykyby/chat-app-backend/api"
	"github.com/fykyby/chat-app-backend/internal/status"
)

type LogOutHandler struct{}

func (h *LogOutHandler) LogOut(w http.ResponseWriter, r *http.Request) {
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

	api.SendResponse(w, http.StatusOK, status.MESSAGE_LOGOUT_SUCCESS, nil)
}

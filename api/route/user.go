package route

import (
	"log"
	"net/http"
	"strconv"

	"github.com/fykyby/chat-app-backend/api"
	"github.com/fykyby/chat-app-backend/internal/auth"
	"github.com/fykyby/chat-app-backend/internal/database"
	"github.com/fykyby/chat-app-backend/internal/model"
	"github.com/fykyby/chat-app-backend/internal/status"
	"github.com/go-chi/jwtauth/v5"
)

const userSearchPageSize = 20

type UserHandler struct {
	DB        *database.Queries
	TokenAuth *jwtauth.JWTAuth
}

func (h *UserHandler) SearchUsers(w http.ResponseWriter, r *http.Request) {
	claimedUser, err := auth.GetClaimedUser(r.Context())
	if err != nil {
		api.SendResponse(w, http.StatusUnauthorized, status.MESSAGE_ERROR_GENERIC, nil)
		return
	}

	query := r.URL.Query().Get("q")
	page_, err := strconv.Atoi(r.URL.Query().Get("page"))
	if err != nil {
		page_ = 1
	}
	page := int32(page_)

	users_, err := h.DB.SearchPublicUsers(r.Context(), database.SearchPublicUsersParams{
		Name:   "%" + query + "%",
		ID:     claimedUser.ID,
		Limit:  userSearchPageSize + 1,
		Offset: (page - 1) * userSearchPageSize,
	})
	if err != nil {
		log.Println(err)
		api.SendResponse(w, http.StatusInternalServerError, status.MESSAGE_ERROR_GENERIC, nil)
		return
	}

	users := []model.PublicUser{}
	for index, user := range users_ {
		if index >= userSearchPageSize {
			break
		}
		users = append(users, model.PublicUser{
			ID:     user.ID,
			Name:   user.Name,
			Avatar: user.Avatar,
		})
	}

	api.SendResponse(w, http.StatusOK, status.MESSAGE_SUCCESS_GENERIC, map[string]interface{}{
		"users":   users,
		"hasMore": len(users) > userSearchPageSize,
	})
}

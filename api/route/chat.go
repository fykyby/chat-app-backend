package route

import (
	"log"
	"net/http"
	"strconv"

	"github.com/fykyby/chat-app-backend/api"
	"github.com/fykyby/chat-app-backend/internal/auth"
	"github.com/fykyby/chat-app-backend/internal/database"
	"github.com/fykyby/chat-app-backend/internal/status"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/jwtauth/v5"
	"github.com/goccy/go-json"
	"github.com/jackc/pgx/v5/pgtype"
)

const chatPageSize = 20

type ChatHandler struct {
	DB        *database.Queries
	TokenAuth *jwtauth.JWTAuth
}

type newChatRequest struct {
	RecipientID int32 `json:"recipientID"`
}

func (h *ChatHandler) GetUserChats(w http.ResponseWriter, r *http.Request) {
	claimedUser, err := auth.GetClaimedUser(r.Context())
	if err != nil {
		log.Println(err)
		api.SendResponse(w, http.StatusUnauthorized, status.MESSAGE_ERROR_GENERIC, nil)
		return
	}

	chats, err := h.DB.GetUserChats(r.Context(), claimedUser.ID)
	if err != nil {
		log.Println(err)
		api.SendResponse(w, http.StatusInternalServerError, status.MESSAGE_ERROR_GENERIC, nil)
		return
	}

	log.Println(chats)

	api.SendResponse(w, http.StatusOK, status.MESSAGE_SUCCESS_GENERIC, chats)
}

func (h *ChatHandler) GetChatMessages(w http.ResponseWriter, r *http.Request) {
	claimedUser, err := auth.GetClaimedUser(r.Context())
	if err != nil {
		log.Println(err)
		api.SendResponse(w, http.StatusUnauthorized, status.MESSAGE_ERROR_GENERIC, nil)
		return
	}

	chatID_, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		log.Println(err)
		api.SendResponse(w, http.StatusBadRequest, status.MESSAGE_ERROR_GENERIC, nil)
		return
	}
	chatID := int32(chatID_)

	_, err = h.DB.GetUsersChat(r.Context(), database.GetUsersChatParams{
		UserID: claimedUser.ID,
		ChatID: chatID,
	})
	if err != nil {
		log.Println(err)
		api.SendResponse(w, http.StatusBadRequest, status.MESSAGE_ERROR_GENERIC, nil)
		return
	}

	messages, err := h.DB.GetMessages(r.Context(), database.GetMessagesParams{
		ChatID: chatID,
		Limit:  chatPageSize,
		Offset: 0,
	})
	if err != nil {
		log.Println(err)
		api.SendResponse(w, http.StatusInternalServerError, status.MESSAGE_ERROR_GENERIC, nil)
		return
	}

	// TODO: join user data to messages

	api.SendResponse(w, http.StatusOK, status.MESSAGE_SUCCESS_GENERIC, messages)
}

func (h *ChatHandler) CreateChat(w http.ResponseWriter, r *http.Request) {
	claimedUser, err := auth.GetClaimedUser(r.Context())
	if err != nil {
		log.Println(err)
		api.SendResponse(w, http.StatusUnauthorized, status.MESSAGE_ERROR_GENERIC, nil)
		return
	}

	var req newChatRequest
	err = json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		log.Println(err)
		api.SendResponse(w, http.StatusBadRequest, status.MESSAGE_ERROR_GENERIC, nil)
		return
	}

	recipient, err := h.DB.GetPublicUser(r.Context(), req.RecipientID)
	if err != nil {
		log.Println(err)
		api.SendResponse(w, http.StatusBadRequest, status.MESSAGE_ERROR_GENERIC, nil)
		return
	}

	chat, err := h.DB.CreateChat(r.Context(), database.CreateChatParams{
		IsGroup: false,
		Avatar: pgtype.Text{
			Valid: false,
		},
		Name: pgtype.Text{
			Valid: false,
		},
	})
	if err != nil {
		log.Println(err)
		api.SendResponse(w, http.StatusInternalServerError, status.MESSAGE_ERROR_GENERIC, nil)
		return
	}

	_, err = h.DB.CreateUserChat(r.Context(), database.CreateUserChatParams{
		UserID: claimedUser.ID,
		ChatID: chat.ID,
	})
	if err != nil {
		log.Println(err)
		h.DB.DeleteChat(r.Context(), chat.ID)

		api.SendResponse(w, http.StatusInternalServerError, status.MESSAGE_ERROR_GENERIC, nil)
		return
	}
	_, err = h.DB.CreateUserChat(r.Context(), database.CreateUserChatParams{
		UserID: recipient.ID,
		ChatID: chat.ID,
	})
	if err != nil {
		log.Println(err)
		h.DB.DeleteChat(r.Context(), chat.ID)

		api.SendResponse(w, http.StatusInternalServerError, status.MESSAGE_ERROR_GENERIC, nil)
		return
	}

	api.SendResponse(w, http.StatusCreated, status.MESSAGE_SUCCESS_GENERIC, chat)
}

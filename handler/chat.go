package handler

import (
	"net/http"
	"strconv"

	"github.com/fykyby/chat-app-backend/auth"
	"github.com/fykyby/chat-app-backend/database"
	"github.com/go-chi/chi/v5"
	"github.com/goccy/go-json"
)

type postChatRequest struct {
	RecipientID int32 `json:"recipientID"`
}

func getChat(w http.ResponseWriter, r *http.Request) {
	claimedUser, err := auth.GetClaimedUser(r.Context())
	if err != nil {
		SendResponse(w, http.StatusUnauthorized, MESSAGE_ERROR_GENERIC, nil)
		return
	}

	chatID_, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		SendResponse(w, http.StatusBadRequest, MESSAGE_ERROR_GENERIC, nil)
		return
	}
	chatID := int32(chatID_)

	_, err = db.GetUsersChat(r.Context(), database.GetUsersChatParams{
		UserID: claimedUser.ID,
		ChatID: chatID,
	})
	if err != nil {
		SendResponse(w, http.StatusBadRequest, MESSAGE_ERROR_GENERIC, nil)
		return
	}

	messages, err := db.GetMessagesPage(r.Context(), database.GetMessagesPageParams{
		ChatID: chatID,
		Limit:  20,
		Offset: 0,
	})
	if err != nil {
		SendResponse(w, http.StatusInternalServerError, MESSAGE_ERROR_GENERIC, nil)
		return
	}

	SendResponse(w, http.StatusOK, "", messages)
}

func postChat(w http.ResponseWriter, r *http.Request) {
	claimedUser, err := auth.GetClaimedUser(r.Context())
	if err != nil {
		SendResponse(w, http.StatusUnauthorized, MESSAGE_ERROR_GENERIC, nil)
		return
	}

	var req postChatRequest
	err = json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		SendResponse(w, http.StatusBadRequest, MESSAGE_ERROR_GENERIC, nil)
		return
	}

	recipient, err := db.GetPublicUser(r.Context(), req.RecipientID)
	if err != nil {
		SendResponse(w, http.StatusBadRequest, MESSAGE_ERROR_GENERIC, nil)
		return
	}

	chat, err := db.CreateChat(r.Context(), recipient.Name)
	if err != nil {
		SendResponse(w, http.StatusInternalServerError, MESSAGE_ERROR_GENERIC, nil)
		return
	}

	_, err = db.CreateUsersChat(r.Context(), database.CreateUsersChatParams{
		UserID: claimedUser.ID,
		ChatID: chat.ID,
	})
	if err != nil {
		db.DeleteChat(r.Context(), chat.ID)

		SendResponse(w, http.StatusInternalServerError, MESSAGE_ERROR_GENERIC, nil)
		return
	}

	SendResponse(w, http.StatusCreated, "", chat)
}

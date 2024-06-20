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
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/jwtauth/v5"
	"github.com/goccy/go-json"
	"github.com/jackc/pgx/v5/pgtype"
)

const CHAT_PAGE_SIZE = 20

type ChatHandler struct {
	DB        *database.Queries
	TokenAuth *jwtauth.JWTAuth
}

type newChatRequest struct {
	RecipientID int32 `json:"recipientID"`
}

func (h *ChatHandler) GetUserChats(w http.ResponseWriter, r *http.Request) {
	claimedUser, err := auth.GetClaimedUser(r.Context(), h.DB)
	if err != nil {
		log.Println(err)
		api.SendResponse(w, http.StatusUnauthorized, status.MESSAGE_ERROR_GENERIC, nil)
		return
	}

	chats_, err := h.DB.GetUserChats(r.Context(), claimedUser.ID)
	if err != nil {
		log.Println(err)
		api.SendResponse(w, http.StatusInternalServerError, status.MESSAGE_ERROR_GENERIC, nil)
		return
	}

	chats := []model.Chat{}

	for _, chat := range chats_ {
		newChat := model.Chat{
			ID:      chat.ID,
			Name:    chat.Name.String,
			Avatar:  chat.Avatar.String,
			IsGroup: chat.IsGroup,
		}

		if !chat.Name.Valid || !chat.Avatar.Valid {
			chatUsers, err := h.DB.GetChatUsers(r.Context(), chat.ID)
			if err != nil {
				log.Println(err)
				api.SendResponse(w, http.StatusInternalServerError, status.MESSAGE_ERROR_GENERIC, nil)
				return
			}

			for _, user := range chatUsers {
				if user.ID != claimedUser.ID {
					if !chat.Name.Valid {
						newChat.Name = user.Name
					}
					if !chat.Avatar.Valid {
						newChat.Avatar = user.Avatar
					}
				}
			}
		}

		chats = append(chats, newChat)
	}

	api.SendResponse(w, http.StatusOK, status.MESSAGE_SUCCESS_GENERIC, chats)
}

func (h *ChatHandler) GetChat(w http.ResponseWriter, r *http.Request) {
	claimedUser, err := auth.GetClaimedUser(r.Context(), h.DB)
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

	page_, err := strconv.Atoi(r.URL.Query().Get("page"))
	if err != nil {
		page_ = 1
	}
	page := int32(page_)

	_, err = h.DB.GetUsersChat(r.Context(), database.GetUsersChatParams{
		UserID: claimedUser.ID,
		ChatID: chatID,
	})
	if err != nil {
		log.Println(err)
		api.SendResponse(w, http.StatusUnauthorized, status.MESSAGE_ERROR_GENERIC, nil)
		return
	}

	chat_, err := h.DB.GetChat(r.Context(), chatID)
	if err != nil {
		log.Println(err)
		api.SendResponse(w, http.StatusInternalServerError, status.MESSAGE_ERROR_GENERIC, nil)
		return
	}

	newChat := model.Chat{
		ID:      chat_.ID,
		Name:    chat_.Name.String,
		Avatar:  chat_.Avatar.String,
		IsGroup: chat_.IsGroup,
	}

	if !chat_.Name.Valid || !chat_.Avatar.Valid {
		chatUsers, err := h.DB.GetChatUsers(r.Context(), chat_.ID)
		if err != nil {
			log.Println(err)
			api.SendResponse(w, http.StatusInternalServerError, status.MESSAGE_ERROR_GENERIC, nil)
			return
		}

		for _, user := range chatUsers {
			if user.ID != claimedUser.ID {
				if !chat_.Name.Valid {
					newChat.Name = user.Name
				}
				if !chat_.Avatar.Valid {
					newChat.Avatar = user.Avatar
				}
			}
		}
	}

	messages_, err := h.DB.GetMessages(r.Context(), database.GetMessagesParams{
		ChatID: chatID,
		Limit:  CHAT_PAGE_SIZE + 1,
		Offset: (page - 1) * CHAT_PAGE_SIZE,
	})
	if err != nil {
		log.Println(err)
		api.SendResponse(w, http.StatusInternalServerError, status.MESSAGE_ERROR_GENERIC, nil)
		return
	}

	messages := []model.Message{}
	for index, message := range messages_ {
		if index >= CHAT_PAGE_SIZE {
			break
		}

		newMessage := model.Message{
			ID:        message.ID,
			Content:   message.Content,
			CreatedAt: message.CreatedAt.Time.String(),
			User: model.PublicUser{
				ID:     message.UserID,
				Name:   message.UserName,
				Avatar: message.UserAvatar,
			},
		}

		messages = append(messages, newMessage)
	}

	response := map[string]interface{}{
		"chat":     newChat,
		"messages": messages,
		"hasMore":  len(messages_) > CHAT_PAGE_SIZE,
	}

	api.SendResponse(w, http.StatusOK, status.MESSAGE_SUCCESS_GENERIC, response)
}

func (h *ChatHandler) CreateChat(w http.ResponseWriter, r *http.Request) {
	claimedUser, err := auth.GetClaimedUser(r.Context(), h.DB)
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

	existingChat, err := h.DB.GetChatOfTwoUsers(r.Context(), database.GetChatOfTwoUsersParams{
		UserID:   claimedUser.ID,
		UserID_2: recipient.ID,
	})
	if err == nil {
		log.Println("Chat already exists")

		chat := model.Chat{
			ID:      existingChat.ID,
			Name:    existingChat.Name.String,
			Avatar:  existingChat.Avatar.String,
			IsGroup: existingChat.IsGroup,
		}

		api.SendResponse(w, http.StatusOK, status.MESSAGE_CHAT_ALREADY_EXISTS, chat)
		return
	}
	log.Println("Creating chat")

	chat_, err := h.DB.CreateChat(r.Context(), database.CreateChatParams{
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
		ChatID: chat_.ID,
	})
	if err != nil {
		log.Println(err)
		h.DB.DeleteChat(r.Context(), chat_.ID)

		api.SendResponse(w, http.StatusInternalServerError, status.MESSAGE_ERROR_GENERIC, nil)
		return
	}

	_, err = h.DB.CreateUserChat(r.Context(), database.CreateUserChatParams{
		UserID: recipient.ID,
		ChatID: chat_.ID,
	})
	if err != nil {
		log.Println(err)
		h.DB.DeleteChat(r.Context(), chat_.ID)

		api.SendResponse(w, http.StatusInternalServerError, status.MESSAGE_ERROR_GENERIC, nil)
		return
	}

	chat := model.Chat{
		ID:      chat_.ID,
		Name:    chat_.Name.String,
		Avatar:  chat_.Avatar.String,
		IsGroup: chat_.IsGroup,
	}

	api.SendResponse(w, http.StatusCreated, status.MESSAGE_SUCCESS_GENERIC, chat)
}

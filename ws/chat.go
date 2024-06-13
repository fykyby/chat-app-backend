package ws

import (
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/fykyby/chat-app-backend/auth"
	"github.com/fykyby/chat-app-backend/handler"
	"github.com/fykyby/chat-app-backend/model"
	"github.com/go-chi/chi/v5"

	// "github.com/go-chi/jwtauth/v5"
	"github.com/gorilla/websocket"
)

type incomingMessage struct {
	UserID  int32  `json:"userID"`
	Content string `json:"content"`
}

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

var rooms = make(map[int32]*room)

func getChatWs(w http.ResponseWriter, r *http.Request) {
	claimedUser, err := auth.GetClaimedUser(r.Context())
	if err != nil {
		handler.SendResponse(w, http.StatusUnauthorized, "Unauthorized", nil)
		return
	}
	// TODO: Check if user is in the chat

	upgrader.CheckOrigin = func(r *http.Request) bool {
		return r.Header.Get("Origin") == os.Getenv("CLIENT_URL")
	}

	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}

	roomID_, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		log.Println(err)
		return
	}
	roomID := int32(roomID_)

	myRoom, ok := rooms[roomID]
	if !ok {
		myRoom = newRoom(roomID)
		rooms[roomID] = myRoom
		go myRoom.run()
	}

	client := &client{
		id:   claimedUser.ID,
		conn: conn,
		room: myRoom,
		send: make(chan model.Message),
	}

	myRoom.register <- client

	go client.readPump()
	go client.writePump()
}

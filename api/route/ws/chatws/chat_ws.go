package chatws

// import (
// 	"log"
// 	"net/http"
// 	"os"
// 	"strconv"

// 	"github.com/fykyby/chat-app-backend/api"
// 	"github.com/fykyby/chat-app-backend/internal/auth"
// 	"github.com/fykyby/chat-app-backend/internal/database"
// 	"github.com/fykyby/chat-app-backend/internal/model"
// 	"github.com/go-chi/chi/v5"

// 	"github.com/gorilla/websocket"
// )

// type ChatWsHandler struct {
// 	DB    *database.Queries
// 	Rooms map[int32]*Room
// }

// type incomingMessage struct {
// 	UserID  int32  `json:"userID"`
// 	Content string `json:"content"`
// }

// var upgrader = websocket.Upgrader{
// 	ReadBufferSize:  1024,
// 	WriteBufferSize: 1024,
// }

// func (h *ChatWsHandler) ConnectToChatWs(w http.ResponseWriter, r *http.Request) {
// 	claimedUser, err := auth.GetClaimedUser(r.Context())
// 	if err != nil {
// 		api.SendResponse(w, http.StatusUnauthorized, "Unauthorized", nil)
// 		return
// 	}
// 	// TODO: Check if user is in the chat

// 	upgrader.CheckOrigin = func(r *http.Request) bool {
// 		return r.Header.Get("Origin") == os.Getenv("CLIENT_URL")
// 	}

// 	conn, err := upgrader.Upgrade(w, r, nil)
// 	if err != nil {
// 		log.Println(err)
// 		return
// 	}

// 	roomID_, err := strconv.Atoi(chi.URLParam(r, "id"))
// 	if err != nil {
// 		log.Println(err)
// 		return
// 	}
// 	roomID := int32(roomID_)

// 	myRoom, ok := h.Rooms[roomID]
// 	if !ok {
// 		myRoom = newRoom(roomID)
// 		h.Rooms[roomID] = myRoom
// 		go myRoom.run()
// 	}

// 	client := &client{
// 		id:   claimedUser.ID,
// 		conn: conn,
// 		room: myRoom,
// 		send: make(chan model.Message),
// 	}

// 	myRoom.register <- client

// 	go client.readPump()
// 	go client.writePump()
// }

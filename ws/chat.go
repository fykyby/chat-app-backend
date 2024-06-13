package ws

import (
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/gorilla/websocket"
)

type incomingMessage struct {
	UserID  int    `json:"userID"`
	Content string `json:"content"`
}

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

var rooms = make(map[int]*room)

func getChatWs(w http.ResponseWriter, r *http.Request) {
	upgrader.CheckOrigin = func(r *http.Request) bool {
		return r.Header.Get("Origin") == os.Getenv("CLIENT_URL")
	}

	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}

	roomID, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		log.Println(err)
		return
	}

	myRoom, ok := rooms[roomID]
	if !ok {
		myRoom = newRoom(roomID)
		rooms[roomID] = myRoom
		go myRoom.run()
	}

	client := &client{
		conn: conn,
		room: myRoom,
		send: make(chan incomingMessage),
	}

	myRoom.register <- client

	go client.readPump()
	go client.writePump()
}

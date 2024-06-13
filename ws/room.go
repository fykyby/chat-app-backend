package ws

import (
	"log"

	"github.com/fykyby/chat-app-backend/model"
)

type room struct {
	id         int32
	clients    map[*client]bool
	broadcast  chan *model.Message
	register   chan *client
	unregister chan *client
}

func newRoom(id int32) *room {
	return &room{
		id:         id,
		clients:    make(map[*client]bool),
		broadcast:  make(chan *model.Message),
		register:   make(chan *client),
		unregister: make(chan *client),
	}

}

func (r *room) run() {
	for {
		select {
		case client := <-r.register:
			log.Println("Client Connected to room ", r.id)

			r.clients[client] = true
		case client := <-r.unregister:
			log.Println("Client Disconnected from room ", r.id)
			if _, ok := r.clients[client]; ok {
				delete(r.clients, client)
				close(client.send)
			}

			if !r.hasClients() {
				log.Println("No clients left, closing room ", r.id)
				r.close()
				return
			}
		case msg := <-r.broadcast:
			for client := range r.clients {
				client.send <- *msg
			}
		}
	}
}

func (r *room) hasClients() bool {
	return len(r.clients) > 0
}

func (r *room) close() {
	close(r.broadcast)
	close(r.register)
	close(r.unregister)

	_, ok := rooms[r.id]
	if ok {
		delete(rooms, r.id)
	}
}
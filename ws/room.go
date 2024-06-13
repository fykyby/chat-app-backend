package ws

import (
	"log"
)

type room struct {
	id         string
	clients    map[*client]bool
	broadcast  chan *message
	register   chan *client
	unregister chan *client
}

func newRoom(id string) *room {
	return &room{
		id:         id,
		clients:    make(map[*client]bool),
		broadcast:  make(chan *message),
		register:   make(chan *client),
		unregister: make(chan *client),
	}

}

func (r *room) run() {
	for {
		select {
		case client := <-r.register:
			log.Println("Client Connected to room " + r.id)

			r.clients[client] = true
		case client := <-r.unregister:
			log.Println("Client Disconnected")
			if _, ok := r.clients[client]; ok {
				delete(r.clients, client)
				close(client.send)
			}

			if !r.hasClients() {
				log.Println("No clients in room, closing room")
				r.close()
				return
			}
		case msg := <-r.broadcast:
			log.Println("Broadcasting Message")

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

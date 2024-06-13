package ws

import (
	"log"

	"github.com/gorilla/websocket"
)

type client struct {
	conn *websocket.Conn
	room *room
	send chan incomingMessage
}

func (c *client) readPump() {
	defer func() {
		c.room.unregister <- c
		c.conn.Close()
	}()

	for {
		msg := incomingMessage{}
		err := c.conn.ReadJSON(&msg)
		if err != nil {
			log.Println(err)
			break
		}

		log.Println(msg)
		// TODO: DO THINGS WITH INCOMING MESSAGE

		c.room.broadcast <- &msg
	}
}

func (c *client) writePump() {
	defer func() {
		c.conn.Close()
	}()

	for msg := range c.send {
		c.conn.WriteJSON(msg)
	}
}

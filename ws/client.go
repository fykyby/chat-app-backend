package ws

import (
	"context"
	"log"

	"github.com/fykyby/chat-app-backend/database"
	"github.com/fykyby/chat-app-backend/model"
	"github.com/gorilla/websocket"
)

type client struct {
	id   int32
	conn *websocket.Conn
	room *room
	send chan model.Message
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

		outgoingMsg, err := db.CreateMessage(context.Background(), database.CreateMessageParams{
			ChatID:  c.room.id,
			UserID:  msg.UserID,
			Content: msg.Content,
		})
		if err != nil {
			log.Println(err)
			break
		}

		publicUser, err := db.GetPublicUser(context.Background(), msg.UserID)
		if err != nil {
			log.Println(err)
			break
		}

		message := model.Message{
			ID:        outgoingMsg.ID,
			Content:   outgoingMsg.Content,
			CreatedAt: outgoingMsg.CreatedAt.Time.String(),
			User: model.PublicUser{
				ID:     publicUser.ID,
				Name:   publicUser.Name,
				Avatar: publicUser.Avatar.String,
			},
		}

		c.room.broadcast <- &message
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

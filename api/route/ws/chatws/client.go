package chatws

import (
	"context"
	"log"

	"github.com/fykyby/chat-app-backend/internal/database"
	"github.com/fykyby/chat-app-backend/internal/model"
	"github.com/gorilla/websocket"
)

type client struct {
	handler *ChatWsHandler
	id      int32
	conn    *websocket.Conn
	room    *Room
	send    chan model.Message
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

		outgoingMsg, err := c.handler.DB.CreateMessage(context.Background(), database.CreateMessageParams{
			ChatID:  c.room.id,
			UserID:  msg.UserID,
			Content: msg.Content,
		})
		if err != nil {
			log.Println(err)
			break
		}

		publicUser, err := c.handler.DB.GetPublicUser(context.Background(), msg.UserID)
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
				Avatar: publicUser.Avatar,
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

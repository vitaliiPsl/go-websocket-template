package sockets

import (
	"encoding/json"
	"log"
	"time"

	"github.com/gorilla/websocket"
)

const (
	pingInterval   = 10 * time.Second
	pongWait       = 15 * time.Second
	maxMessageSize = 1024
)

type Client struct {
	userId   string
	conn     *websocket.Conn
	manager  WebsocketManager
	messages chan WebsocketMessage
}

func NewWebsocketClient(userId string, conn *websocket.Conn, manager WebsocketManager) *Client {
	return &Client{
		userId:   userId,
		conn:     conn,
		manager:  manager,
		messages: make(chan WebsocketMessage, 16),
	}
}

func (c *Client) ReadMessage() {
	log.Printf("Reading messages from the webscoket client. userId=%v...", c.userId)
	defer c.Close()

	c.conn.SetReadLimit(maxMessageSize)
	c.conn.SetReadDeadline(time.Now().Add(pongWait))
	c.conn.SetPongHandler(func(appData string) error {
		c.conn.SetReadDeadline(time.Now().Add(pongWait))
		return nil
	})

	for {
		_, message, err := c.conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("Unexpected close error. userId=%v, err=%v", c.userId, err)
			} else {
				log.Printf("Error while reading websocket message. userId=%v, err=%v", c.userId, err)
			}
			return
		}

		wsMessage := WebsocketMessage{}
		if err = json.Unmarshal(message, &wsMessage); err != nil {
			log.Printf("Error while unmarshaling websocket message. Err=%v\n", err)
			return
		}

		log.Printf("Received websocket message: userId=%v", c.userId)
		c.manager.HandleMessage(c.userId, wsMessage)
	}
}

func (c *Client) WriteMessage() {
	log.Printf("Writing messsages to the websocket client. userId=%v...", c.userId)
	defer c.Close()

	ticker := time.NewTicker(pingInterval)
	defer ticker.Stop()

	for {
		select {
		case message := <-c.messages:
			if err := c.conn.WriteJSON(message); err != nil {
				log.Printf("Error while sending message. userId=%v, err=%v", c.userId, err)
				return
			}
		case <-ticker.C:
			log.Println("Sending ping")
			if err := c.conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				log.Printf("Error while sending ping. userId=%v, err=%v", c.userId, err)
				return
			}
		}
	}
}

func (c *Client) Close() {
	log.Printf("Closing websocket client. userId=%v...", c.userId)
	c.conn.Close()
	c.manager.UnregisterClient(c.userId)
}

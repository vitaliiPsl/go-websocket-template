package sockets

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

type WebsocketHandler interface {
	Handle(ctx *gin.Context)
}

type WebsocketHandlerImpl struct {
	manager WebsocketManager
}

func NewWebsocketHander(manager WebsocketManager) *WebsocketHandlerImpl {
	return &WebsocketHandlerImpl{
		manager: manager,
	}
}

func (h *WebsocketHandlerImpl) Handle(ctx *gin.Context) {
	// should be replaced with proper authentication(e.g. JWT token)
	userId := ctx.Query("user_id")
	if userId == "" {
		log.Println("User id is missing")
		http.Error(ctx.Writer, "User id is required", http.StatusBadRequest)
		return
	}

	conn, err := upgrader.Upgrade(ctx.Writer, ctx.Request, nil)
	if err != nil {
		log.Println("Failed to upgrade websocket connection", err)
		http.Error(ctx.Writer, "Failed to upgrade connection", http.StatusBadRequest)
		return
	}

	client := NewWebsocketClient(userId, conn, h.manager)

	h.manager.RegisterClient(client)

	go client.ReadMessage()
	go client.WriteMessage()
}

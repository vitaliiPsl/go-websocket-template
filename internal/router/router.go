package router

import (
	"net/http"
	"websocket-template/internal/sockets"

	"github.com/gin-gonic/gin"
)

type Router interface {
	Handler() http.Handler
}

type RouterImpl struct {
	Engine *gin.Engine
}

func NewRouter(websocketHandler sockets.WebsocketHandler) *RouterImpl {
	router := gin.Default()

	v1 := router.Group("v1")
	v1.GET("health", func(ctx *gin.Context) {
		ctx.JSON(http.StatusOK, gin.H{"status": "Ok"})
	})

	v1.GET("/ws", websocketHandler.Handle)

	return &RouterImpl{Engine: router}
}

func (r *RouterImpl) Handler() http.Handler {
	return r.Engine
}

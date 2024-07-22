package main

import (
	"log"
	"websocket-template/internal/router"
	"websocket-template/internal/server"
	"websocket-template/internal/sockets"
)

func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	socketManager := sockets.NewWebsocketManager()
	socketHandler := sockets.NewWebsocketHander(socketManager)

	router := router.NewRouter(socketHandler)
	server.Serve(router)
}

package sockets

import (
	"log"
	"sync"
)

type WebsocketManager interface {
	RegisterClient(client *Client)
	UnregisterClient(userId string)
	HandleMessage(userId string, message WebsocketMessage)
}

type WebsocketManagerImpl struct {
	clients map[string]*Client
	mu      sync.Mutex
}

func NewWebsocketManager() *WebsocketManagerImpl {
	return &WebsocketManagerImpl{
		clients: make(map[string]*Client),
	}
}

func (manager *WebsocketManagerImpl) RegisterClient(client *Client) {
	log.Printf("Registering websocket client: userId=%v", client.userId)

	manager.mu.Lock()
	defer manager.mu.Unlock()

	manager.clients[client.userId] = client
}

func (manager *WebsocketManagerImpl) UnregisterClient(userId string) {
	log.Printf("Unregistering websocket client: userId=%v", userId)

	manager.mu.Lock()
	defer manager.mu.Unlock()
	
	delete(manager.clients, userId)
}

func (manager *WebsocketManagerImpl) HandleMessage(userId string, message WebsocketMessage) {
	log.Printf("Handling websocket message: event=%v", message.Event)

	manager.mu.Lock()
	defer manager.mu.Unlock()
	
	for userId, client := range manager.clients {
		log.Printf("Broadcasting message: userId=%v", userId)
		client.messages <- message
	}
}

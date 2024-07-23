package sockets

import (
	"encoding/json"
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

	payload := ConnectedPayload{UserId: client.userId}
	go manager.broadcastEvent(EventTypeConnected, payload)
}

func (manager *WebsocketManagerImpl) UnregisterClient(userId string) {
	log.Printf("Unregistering websocket client: userId=%v", userId)

	manager.mu.Lock()
	defer manager.mu.Unlock()

	delete(manager.clients, userId)

	payload := ConnectedPayload{UserId: userId}
	go manager.broadcastEvent(EventTypeDisconnected, payload)
}

func (manager *WebsocketManagerImpl) HandleMessage(userId string, message WebsocketMessage) {
	log.Printf("Handling websocket message: event=%v", message.Event)

	switch message.Event {
	case EventTypeMessage:
		manager.broadcastMessage(message)
	default:
		log.Printf("Unknown event: event=%v", message.Event)
	}
}

func (manager *WebsocketManagerImpl) broadcastEvent(eventType EventType, payload interface{}) {
	log.Printf("Broadcasting websocket event: event=%v", eventType)

	payloadJSON, err := json.Marshal(payload)
	if err != nil {
		log.Printf("Failed to marshal event payload: event=%v, err=%v", eventType, err)
		return
	}

	websocketMessage := WebsocketMessage{
		Event:   eventType,
		Payload: payloadJSON,
	}
	manager.broadcastMessage(websocketMessage)
}

func (manager *WebsocketManagerImpl) broadcastMessage(message WebsocketMessage) {
	log.Printf("Broadcasting websocket message: event=%v", message.Event)

	manager.mu.Lock()
	defer manager.mu.Unlock()

	for userId, client := range manager.clients {
		log.Printf("Sending broadcast message: userId=%v", userId)
		client.messages <- message
	}
}

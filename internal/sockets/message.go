package sockets

import "encoding/json"

type EventType string

const (
	EventTypeConnected    EventType = "connected"
	EventTypeDisconnected EventType = "disconnected"
	EventTypeMessage      EventType = "message"
)

type WebsocketMessage struct {
	Event   EventType       `json:"event"`
	Payload json.RawMessage `json:"payload,omitempty"`
}

type ConnectedPayload struct {
	UserId string `json:"user_id"`
}

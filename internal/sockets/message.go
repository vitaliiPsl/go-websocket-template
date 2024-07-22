package sockets

import "encoding/json"

type EventType string

const (
	EventTypeMessage EventType = "message"
)

type WebsocketMessage struct {
	Event    EventType       `json:"event"`
	Payload json.RawMessage `json:"payload,omitempty"`
}

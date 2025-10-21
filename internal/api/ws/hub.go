package ws

import (
	"encoding/json"
	"sync"
)

// Hub manages websocket clients per conversation.
type Hub struct {
	mu      sync.RWMutex
	clients map[string]map[*Client]struct{}
}

// NewHub constructs a Hub instance.
func NewHub() *Hub {
	return &Hub{clients: make(map[string]map[*Client]struct{})}
}

// Register adds a client to the hub.
func (h *Hub) Register(client *Client) {
	h.mu.Lock()
	defer h.mu.Unlock()

	convClients := h.clients[client.conversationID]
	if convClients == nil {
		convClients = make(map[*Client]struct{})
		h.clients[client.conversationID] = convClients
	}
	convClients[client] = struct{}{}
}

// Unregister removes a client from the hub.
func (h *Hub) Unregister(client *Client) {
	h.mu.Lock()
	defer h.mu.Unlock()

	if convClients, ok := h.clients[client.conversationID]; ok {
		delete(convClients, client)
		if len(convClients) == 0 {
			delete(h.clients, client.conversationID)
		}
	}
}

// Broadcast sends an event to all clients watching the conversation.
func (h *Hub) Broadcast(conversationID, event string, payload interface{}) {
	message := map[string]interface{}{
		"type": event,
		"data": payload,
	}

	data, err := json.Marshal(message)
	if err != nil {
		return
	}

	h.mu.RLock()
	clients := h.clients[conversationID]
	for client := range clients {
		select {
		case client.send <- data:
		default:
			go client.Close()
		}
	}
	h.mu.RUnlock()
}

// Copyright 2013 The Gorilla WebSocket Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package websocket

import (
	"encoding/json"

	"hinoob.net/learn-go/internal/database"
	"hinoob.net/learn-go/internal/model"
)

// Hub maintains the set of active clients and broadcasts messages to the
// clients.
type Hub struct {
	Clients         map[uint]*Client
	LiveRooms       map[uint]map[*Client]bool
	Register        chan *Client
	Unregister      chan *Client
	DirectMessage   chan *model.Message
	LiveChatMessage chan *LiveMessage
}

// LiveMessage defines a message sent within a live stream chat.
type LiveMessage struct {
	CourseID uint   `json:"course_id"`
	SenderID uint   `json:"sender_id"`
	Content  string `json:"content"`
}

func NewHub() *Hub {
	return &Hub{
		DirectMessage:   make(chan *model.Message),
		LiveChatMessage: make(chan *LiveMessage),
		LiveRooms:       make(map[uint]map[*Client]bool),
		Register:        make(chan *Client),
		Unregister:      make(chan *Client),
		Clients:         make(map[uint]*Client),
	}
}

func (h *Hub) Run() {
	for {
		select {
		case client := <-h.Register:
			h.Clients[client.UserID] = client
		case client := <-h.Unregister:
			for courseID, room := range h.LiveRooms {
				if _, ok := room[client]; ok {
					delete(h.LiveRooms[courseID], client)
				}
			}
			if _, ok := h.Clients[client.UserID]; ok {
				delete(h.Clients, client.UserID)
				close(client.Send)
			}
		case message := <-h.DirectMessage:
			// 1. Save the message to the database
			database.CreateMessage(message)

			// 2. Forward the message to the recipient if they are online
			if client, ok := h.Clients[message.RecipientID]; ok {
				// Marshal the message to JSON to send over WebSocket
				msgBytes, err := json.Marshal(message)
				if err == nil {
					select {
					case client.Send <- msgBytes:
					default:
						close(client.Send)
						delete(h.Clients, client.UserID)
					}
				}
			}
		case message := <-h.LiveChatMessage:
			// Live messages are ephemeral and not saved to the database.
			if room, ok := h.LiveRooms[message.CourseID]; ok {
				msgBytes, err := json.Marshal(message)
				if err == nil {
					for client := range room {
						select {
						case client.Send <- msgBytes:
						default:
							close(client.Send)
							delete(h.Clients, client.UserID)
							delete(room, client)
						}
					}
				}
			}
		}
	}
}

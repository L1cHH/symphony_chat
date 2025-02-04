package server

import (
	"symphony_chat/internal/infrastructure/websocket/client"

	"github.com/google/uuid"
)

type Hub struct {
	chatRooms map[uuid.UUID]map[*client.Client]bool
	broadcast chan []byte
	register chan *client.Client
	unregister chan *client.Client
}

func NewHub() *Hub {
	return &Hub {
		chatRooms: make(map[uuid.UUID]map[*client.Client]bool),
		broadcast: make(chan []byte),
		register: make(chan *client.Client),
		unregister: make(chan *client.Client),
	}
}

func (h *Hub) Run() {
	for {
		select {
		case client := <-h.register:
			h.clients[client] = true
		case client := <-h.unregister:
			if _, exists := h.clients[client]; exists {
				delete(h.clients, client)
				client.CloseBufferChannels()
			}
		}
	}
}

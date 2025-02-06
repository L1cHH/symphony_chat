package chathub

import (
	"symphony_chat/internal/infrastructure/websocket/client"
	"sync"

	"github.com/google/uuid"
)

type Hub struct {
	//Active clients
	activeClients map[uuid.UUID]*client.Client
	//Active chats
	activeChats map[uuid.UUID]map[uuid.UUID]*client.Client

	register chan *client.Client
	unregister chan *client.Client

	mu sync.Mutex
}

func NewHub() *Hub {
	return &Hub {
		activeClients: make(map[uuid.UUID]*client.Client),
		activeChats: make(map[uuid.UUID]map[uuid.UUID]*client.Client),
		register: make(chan *client.Client),
		unregister: make(chan *client.Client),
	}
}


//This method adds client to active clients and adds clients' chats to active chats
func (h *Hub) AddActiveClient(newClient *client.Client, chatIDs ...uuid.UUID) {
	h.mu.Lock()
	h.activeClients[newClient.GetID()] = newClient
	if len(chatIDs) != 0 {
		//user is in chats
		for _, chatID := range chatIDs {
			if _, exists := h.activeChats[chatID]; !exists {
				h.activeChats[chatID] = make(map[uuid.UUID]*client.Client)
				h.activeChats[chatID][newClient.GetID()] = newClient
			} else {
				h.activeChats[chatID][newClient.GetID()] = newClient
			}
		}
	}
	h.mu.Unlock()
}

//This method needs to be used when active client disconnects
func (h *Hub) RemoveActiveClient(client *client.Client, chatIDs ...uuid.UUID) {
	h.mu.Lock()
	for _, chatID := range chatIDs {
		delete(h.activeChats[chatID], client.GetID())
		if len(h.activeChats[chatID]) == 0 {
			delete(h.activeChats, chatID)
		}
	}
	delete(h.activeClients, client.GetID())
	h.mu.Unlock()
}

//This method needs to be used when active client creates new chat
func (h *Hub) AddCreatedChat(chatID uuid.UUID, chatOwner *client.Client) {
	h.mu.Lock()
	h.activeChats[chatID][chatOwner.GetID()] = chatOwner
	h.mu.Unlock()
}

//This method needs to be used when active user deletes chat
func (h *Hub) RemoveActiveChat(chatID uuid.UUID) {
	h.mu.Lock()
	delete(h.activeChats, chatID)
	h.mu.Unlock()
}

func (h *Hub) Run() {
	for {
		select {
		case client := <-h.register:
			h.mu.Lock()
			h.activeClients[client.GetID()] = client
			h.mu.Unlock()
		case client := <-h.unregister:
			h.mu.Lock()
			delete(h.activeClients, client.GetID())
			
			h.mu.Unlock()
		}
	}
}





package chathub

import (
	"encoding/json"
	"log"
	"symphony_chat/internal/infrastructure/websocket/client"
	websocketmessage "symphony_chat/internal/infrastructure/websocket/websocket_message"
	"sync"

	"github.com/google/uuid"
)



type Hub struct {
	//Active clients
	activeClients map[uuid.UUID]*client.Client
	//Active chats
	activeChats map[uuid.UUID]map[uuid.UUID]*client.Client

	

	mu sync.Mutex
}

func NewHub() *Hub {
	return &Hub {
		activeClients: make(map[uuid.UUID]*client.Client),
		activeChats: make(map[uuid.UUID]map[uuid.UUID]*client.Client),
	}
}


//This method adds client to active clients and adds clients' chats to active chats
func (h *Hub) AddActiveClient(newClient *client.Client, chatIDs ...uuid.UUID) {
	h.mu.Lock()
	defer h.mu.Unlock()
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
	
}

//This method needs to be used when active client disconnects
func (h *Hub) RemoveActiveClient(client *client.Client, chatIDs ...uuid.UUID) {
	h.mu.Lock()
	defer h.mu.Unlock()
	for _, chatID := range chatIDs {
		delete(h.activeChats[chatID], client.GetID())
		if len(h.activeChats[chatID]) == 0 {
			delete(h.activeChats, chatID)
		}
	}
	delete(h.activeClients, client.GetID())
}

//This method needs to be used when active client creates new chat
func (h *Hub) AddCreatedChat(chatID uuid.UUID, chatOwner *client.Client) {
	h.mu.Lock()
	defer h.mu.Unlock()
	h.activeChats[chatID][chatOwner.GetID()] = chatOwner
}


//This method needs to be used when active user deletes chat
func (h *Hub) RemoveActiveChat(chatID uuid.UUID) {
	h.mu.Lock()
	defer h.mu.Unlock()
	delete(h.activeChats, chatID)
}

func (h *Hub) HandleMessage(message []byte) {
	var msg websocketmessage.WsMessage 
	if err := json.Unmarshal(message, &msg); err != nil {
		log.Printf("error unmarshalling message: %v", err)
		return
	}
	//TODO: handle message
}

func (h *Hub) Run() {
	
}





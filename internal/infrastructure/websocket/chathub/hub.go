package chathub

import (
	"context"
	"encoding/json"
	"log"
	actions "symphony_chat/internal/domain/chat_actions"
	"symphony_chat/internal/infrastructure/websocket/client"
	websocketmessage "symphony_chat/internal/infrastructure/websocket/websocket_message"
	"symphony_chat/internal/service/chat"
	"sync"

	"github.com/google/uuid"
)



type Hub struct {
	//Active clients
	activeClients map[uuid.UUID]*client.Client
	//Active chats
	activeChats map[uuid.UUID]map[uuid.UUID]*client.Client

	chatService *service.ChatService

	mu sync.RWMutex
}

func NewHub(chatService *service.ChatService) *Hub {
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

func (h *Hub) GetActiveClient(userID uuid.UUID) *client.Client {
	h.mu.RLock()
	defer h.mu.RUnlock()
	return h.activeClients[userID]
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

//This method handles messages from clients
func (h *Hub) HandleMessage(message []byte) {
	var msg websocketmessage.WsMessageRequest 
	if err := json.Unmarshal(message, &msg); err != nil {
		log.Printf("error unmarshalling message: %v", err)
		return
	}

	switch msg.ChatAction {
	case actions.CreateChatAction:
		userID, _ := uuid.Parse(msg.Payload["user_id"].(string))
		chatName, _ := msg.Payload["chat_name"].(string)
		activeClient := h.GetActiveClient(userID)

		chat, err := h.chatService.CreateChat(context.Background(), userID, chatName)
		if err != nil {
			if activeClient.IsStillConnected() {
				activeClient.GetMessageFromServer([]byte(err.Error()))
			}
			return
		}

		h.AddCreatedChat(chat.GetID(), activeClient)
		if activeClient.IsStillConnected() {
			wsRes := websocketmessage.WsMessageResponse {
				ChatActionResult: actions.Success,
				Payload: map[string]interface{} {
					"chat_id": chat.GetID().String(),
					"chat_name": chat.GetName(),
					"chat_created_at": chat.GetCreatedAt(),
				},
			}
			wsResBytes, _ := json.Marshal(wsRes)

			activeClient.GetMessageFromServer(wsResBytes)
		}
    case actions.DeleteChatAction:
		chatID, _ := uuid.Parse(msg.Payload["chat_id"].(string))
		userID, _ := uuid.Parse(msg.Payload["user_id"].(string))
		activeClient := h.GetActiveClient(userID)

		err := h.chatService.DeleteChat(context.Background(), chatID, userID)
		if err != nil {
			if activeClient.IsStillConnected() {
				activeClient.GetMessageFromServer([]byte(err.Error()))
			}
			return
		}

		h.RemoveActiveChat(chatID)

		if activeClient.IsStillConnected() {
			wsRes := websocketmessage.WsMessageResponse {
				ChatActionResult: actions.Success,
				Payload: map[string]interface{} {
					"chat_id": chatID.String(),
				},
			}
			wsResBytes, _ := json.Marshal(wsRes)

			activeClient.GetMessageFromServer(wsResBytes)
		}
	case actions.RenameChatAction:
		chatID, _ := uuid.Parse(msg.Payload["chat_id"].(string))
		userID, _ := uuid.Parse(msg.Payload["user_id"].(string))
		newChatName, _ := msg.Payload["new_chat_name"].(string)
		activeClient := h.GetActiveClient(userID)

		newName, err :=h.chatService.RenameChat(context.Background(), chatID, newChatName, userID)

		if err != nil {
			if activeClient.IsStillConnected(){
				activeClient.GetMessageFromServer([]byte(err.Error()))
			}
			return
		}

		if activeClient.IsStillConnected() {
			wsRes := websocketmessage.WsMessageResponse {
				ChatActionResult: actions.Success,
				Payload: map[string]interface{} {
					"chat_id": chatID,
					"new_chat_name": newName,
				},
			}
			wsResBytes, _ := json.Marshal(wsRes)

			activeClient.GetMessageFromServer(wsResBytes)
		}
	}
}






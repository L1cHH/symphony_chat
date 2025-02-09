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
	"time"

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

func (h *Hub) GetActiveClientsOfChat(chatID uuid.UUID) []*client.Client {
	h.mu.RLock()
	defer h.mu.RUnlock()
	chatClients := make([]*client.Client, 0, len(h.activeChats[chatID]))
	for _, client := range h.activeChats[chatID] {
		chatClients = append(chatClients, client)
	}
	return chatClients
}

//This is the general method for sending events to clients of the 
func (h *Hub) SendWsEventToChatClients(chatID uuid.UUID, clients []*client.Client, wsEvent websocketmessage.WsClientEvent) {
	wsEventBytes, _ := json.Marshal(wsEvent)

	for _, client := range clients {
		if client.IsStillConnected() {
			client.GetMessageFromServer(wsEventBytes)
		}
	}
}

//This is the general method for sending response to the client's request
func (h *Hub) SendWsResponseToClient(client *client.Client, wsResponse websocketmessage.WsMessageResponse) {
	wsResponseBytes, _ := json.Marshal(wsResponse)
	client.GetMessageFromServer(wsResponseBytes)
}

func (h *Hub) AddActiveClientToChat(chatID uuid.UUID, invitedUserID uuid.UUID) {
	h.mu.Lock()
	defer h.mu.Unlock()
	activeClient := h.GetActiveClient(invitedUserID)

	h.activeChats[chatID][invitedUserID] = activeClient
}

//This method needs to be used when active client leaves chat
func (h *Hub) ActiveClientLeftChat(userID uuid.UUID, chatID uuid.UUID) {
	h.mu.Lock()
	defer h.mu.Unlock()

	delete(h.activeChats[chatID], userID)

	//if no more clients in chat, remove chat from active chats
	if len(h.activeChats[chatID]) == 0 {
		delete(h.activeChats, chatID)
	} 
}

func (h *Hub) ActiveClientWasKickedFromChat(chatID uuid.UUID, kickerUserID uuid.UUID,kickedUserID uuid.UUID) {
	h.mu.Lock()
	defer h.mu.Unlock()

	delete(h.activeChats[chatID], kickedUserID)
}

//This method needs to be used when active client disconnects
func (h *Hub) RemoveActiveClient(client *client.Client) {
	h.mu.Lock()
	defer h.mu.Unlock()

	chatIDs, err := h.chatService.GetChatsOfUser(context.Background(), client.GetID())
	if err != nil {
		//TODO: handle error
	}

	for _, chatID := range chatIDs {
		delete(h.activeChats[chatID.GetID()], client.GetID())
		if len(h.activeChats[chatID.GetID()]) == 0 {
			delete(h.activeChats, chatID.GetID())
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

		newName, err := h.chatService.RenameChat(context.Background(), chatID, newChatName, userID)

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

			go h.SendWsResponseToClient(activeClient, wsRes)
		
		}

		wsEvent := websocketmessage.WsClientEvent {
			EventType: actions.ChatNameUpdatedEvent,
			Payload: map[string]interface{} {
				"chat_id": chatID,
				"user_id": userID,
				"new_chat_name": newName,
			},
		}
		go h.SendWsEventToChatClients(chatID, h.GetActiveClientsOfChat(chatID), wsEvent)
	case actions.LeaveChatAction:

		chatID, _ := uuid.Parse(msg.Payload["chat_id"].(string))
		userID, _ := uuid.Parse(msg.Payload["user_id"].(string))
		activeClient := h.GetActiveClient(userID)

		err := h.chatService.LeaveChat(context.Background(), chatID, userID)
		if err != nil {
			if activeClient.IsStillConnected() {
				activeClient.GetMessageFromServer([]byte(err.Error()))
			}
			return
		}

		h.ActiveClientLeftChat(userID, chatID)
		if activeClient.IsStillConnected() {
			wsRes := websocketmessage.WsMessageResponse {
				ChatActionResult: actions.Success,
				Payload: map[string]interface{} {
					"chat_id": chatID,
				},
			}
			go h.SendWsResponseToClient(activeClient, wsRes)
		}

		wsEvent := websocketmessage.WsClientEvent {
			EventType: actions.UserLeftChatEvent,
			Payload: map[string]interface{}{
				"user_id": userID,
				"chat_id": chatID,
				"left_at": time.Now(),
			},
		}
		go h.SendWsEventToChatClients(chatID, h.GetActiveClientsOfChat(chatID), wsEvent)
	case actions.AddMemberToChatAction:
		chatID, _ := uuid.Parse(msg.Payload["chat_id"].(string))
		userID, _ := uuid.Parse(msg.Payload["user_id"].(string))
		invitedUserID, _ := uuid.Parse(msg.Payload["invited_user_id"].(string))

		activeClient := h.GetActiveClient(userID)

		err := h.chatService.AddUserToChat(context.Background(), chatID, userID, invitedUserID)
		if err != nil {
			if activeClient.IsStillConnected() {
				activeClient.GetMessageFromServer([]byte(err.Error()))
			}
			return
		}

		h.AddActiveClientToChat(chatID, invitedUserID)

		if activeClient.IsStillConnected() {
			wsRes := websocketmessage.WsMessageResponse {
				ChatActionResult: actions.Success,
				Payload: map[string]interface{} {
					"chat_id": chatID,
					"inviter_user_id": userID,
					"invited_user_id": invitedUserID,
				},
			}

			go h.SendWsResponseToClient(activeClient, wsRes)
		}

		wsEvent := websocketmessage.WsClientEvent {
			EventType: actions.UserSentMessageEvent,
			Payload: map[string]interface{} {
				"chat_id": chatID,
				"inviter_user_id": userID,
				"invited_user_id": invitedUserID,
			},
		}

		go h.SendWsEventToChatClients(chatID, h.GetActiveClientsOfChat(chatID), wsEvent)
	case actions.RemoveMemberFromChatAction:
		chatID, _ := uuid.Parse(msg.Payload["chat_id"].(string))
		userID, _ := uuid.Parse(msg.Payload["user_id"].(string))
		removingUserID, _ := uuid.Parse(msg.Payload["removing_user_id"].(string))

		activeClient := h.GetActiveClient(userID)

		err := h.chatService.RemoveUserFromChat(context.Background(), chatID, userID, removingUserID)
		if err != nil {
			if activeClient.IsStillConnected() {
				activeClient.GetMessageFromServer([]byte(err.Error()))
			}
			return
		}

		h.ActiveClientWasKickedFromChat(chatID, userID, removingUserID)

		if activeClient.IsStillConnected() {
			wsRes := websocketmessage.WsMessageResponse {
				ChatActionResult: actions.Success,
				Payload: map[string]interface{} {
					"chat_id": chatID,
					"kicker_user_id": userID,
					"kicked_user_id": removingUserID,
				},
			}

			go h.SendWsResponseToClient(activeClient, wsRes)
		}

		wsEvent := websocketmessage.WsClientEvent {
			EventType: actions.UserWasKickedFromChatEvent,
			Payload: map[string]interface{} {
				"chat_id": chatID,
				"kicker_user_id": userID,
				"kicked_user_id": removingUserID,
			},
		}
	
		go h.SendWsEventToChatClients(chatID, h.GetActiveClientsOfChat(chatID), wsEvent)
	case actions.PromoteUserToChatAdminAction:
		chatID, _ := uuid.Parse(msg.Payload["chat_id"].(string))
		userID, _ := uuid.Parse(msg.Payload["user_id"].(string))
		promotedUserID, _ := uuid.Parse(msg.Payload["promoted_user_id"].(string))

		activeClient := h.GetActiveClient(userID)

		err := h.chatService.PromoteUserToChatAdmin(context.Background(), chatID, userID, promotedUserID)
		if err != nil {
			if activeClient.IsStillConnected() {
				activeClient.GetMessageFromServer([]byte(err.Error()))
			}
			return
		}

		if activeClient.IsStillConnected() {
			wsRes := websocketmessage.WsMessageResponse {
				ChatActionResult: actions.Success,
				Payload: map[string]interface{} {
					"chat_id": chatID,
					"promoter_user_id": userID,
					"promoted_user_id": promotedUserID,
				},
			}

			go h.SendWsResponseToClient(activeClient, wsRes)
		}

		wsEvent := websocketmessage.WsClientEvent {
			EventType: actions.UserWasPromotedToChatAdminEvent,
			Payload: map[string]interface{} {
				"chat_id": chatID,
				"promoter_user_id": userID,
				"promoted_user_id": promotedUserID,
			},
		}
	
		go h.SendWsEventToChatClients(chatID, h.GetActiveClientsOfChat(chatID), wsEvent)
	}
}






package chats

import (
	"symphony_chat/internal/domain/chats"
	"symphony_chat/internal/domain/messages"
	"symphony_chat/internal/domain/users"

	"github.com/google/uuid"
)

type ChatService struct {
	chatUserRepo users.ChatUserRepository
	directMsgRepo messages.DirectMessageRepository
	chatRepo chats.ChatRepository
}

type ChatConfiguration func(cs *ChatService) error


func NewChatService(configs ...ChatConfiguration) (*ChatService, error) {
	cs := &ChatService{}
	
	for _, cfg := range configs {
		err := cfg(cs)
		if err != nil {
			return nil, err
		}
	}

	return cs, nil
}

func (cs *ChatService) CreateChatUser(authUserID uuid.UUID, username string) users.ChatUser {
	return users.NewChatUser(authUserID, username)
}


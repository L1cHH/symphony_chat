package service

import (
	"symphony_chat/internal/domain/users"
	"symphony_chat/internal/domain/chat"
	"symphony_chat/internal/domain/roles"
)

type ChatService struct {
	chatUserRepo   users.ChatUserRepository
	chatRepo       chat.ChatRepository
	chatRolesRepo  roles.ChatRoleRepository
}

type ChatServiceConfiguration func(*ChatService) error

func WithChatUserRepository(chatUserRepo users.ChatUserRepository) ChatServiceConfiguration {
	return func(s *ChatService) error {
		s.chatUserRepo = chatUserRepo
		return nil
	}
}

func WithChatRepository(chatRepo chat.ChatRepository) ChatServiceConfiguration {
	return func(s *ChatService) error {
		s.chatRepo = chatRepo
		return nil
	}
}

func WithChatRolesRepository(chatRolesRepo roles.ChatRoleRepository) ChatServiceConfiguration {
	return func(s *ChatService) error {
		s.chatRolesRepo = chatRolesRepo
		return nil
	}
}

func NewChatService(configs ...ChatServiceConfiguration) (*ChatService, error) {
	cs := &ChatService{}

	for _, cfg := range configs {
		err := cfg(cs)
		if err != nil {
			return nil, err
		}
	}

	return cs, nil
}


	
package service

import (
	"context"
	"symphony_chat/internal/application/transaction"
	"symphony_chat/internal/domain/chat"
	"symphony_chat/internal/domain/messages"
	"symphony_chat/internal/domain/roles"
	"symphony_chat/internal/domain/users"

	"github.com/google/uuid"
)

type ChatService struct {
	chatUserRepo        users.ChatUserRepository
	chatRepo            chat.ChatRepository
	chatRolesRepo       roles.ChatRoleRepository
	chatMessageRepo     messages.ChatMessageRepository
	transactionManager  transaction.TransactionManager
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

func WithChatMessageRepository(chatMessageRepo messages.ChatMessageRepository) ChatServiceConfiguration {
	return func(s *ChatService) error {
		s.chatMessageRepo = chatMessageRepo
		return nil
	}
}

func WithTransactionManager(tm transaction.TransactionManager) ChatServiceConfiguration {
	return func(s *ChatService) error {
		s.transactionManager = tm
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


func (cs *ChatService) CreateChat(ctx context.Context, createrID uuid.UUID, chatName string) (chat.Chat, error) {
	createdChat, err := chat.NewChat(chatName)
	if err != nil {
		return chat.Chat{}, err
	}

	err = cs.transactionManager.WithinTransaction(ctx, func(txCtx context.Context) error {
		if err = cs.chatRepo.AddChat(txCtx, createdChat); err != nil {
			return err
		}


		

		return nil
	})

	if err != nil {
		return chat.Chat{}, err
	}

	return createdChat, nil
}



	
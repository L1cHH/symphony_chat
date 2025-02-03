package service

import (
	"context"
	"symphony_chat/internal/application/transaction"
	"symphony_chat/internal/domain/chat"
	"symphony_chat/internal/domain/chat_participant"
	"symphony_chat/internal/domain/messages"
	"symphony_chat/internal/domain/roles"
	"symphony_chat/internal/domain/users"
	"time"

	"github.com/google/uuid"
)

type ChatService struct {
	chatUserRepo        users.ChatUserRepository
	chatParticipantRepo chatparticipant.ChatParticipantRepository
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

func WithChatParticipantRepository(chatParticipantRepo chatparticipant.ChatParticipantRepository) ChatServiceConfiguration {
	return func(s *ChatService) error {
		s.chatParticipantRepo = chatParticipantRepo
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

		if err = cs.CreateChatOwner(txCtx, createdChat.GetID(), createrID); err != nil {
			return err
		}
		
		return nil
	})

	if err != nil {
		return chat.Chat{}, err
	}

	return createdChat, nil
}

func (cs *ChatService) DeleteChat(ctx context.Context, chatID uuid.UUID, deletingInitiatorID uuid.UUID) error {

	isOwner, err := cs.IsOwner(ctx, chatID, deletingInitiatorID)
	if err != nil {
		return err
	}

	if !isOwner {
		return roles.ErrInsufficientPermissions
	}

	err = cs.transactionManager.WithinTransaction(ctx, func(txCtx context.Context) error {

		//Deleting chat participants
		if err := cs.chatParticipantRepo.DeleteAllChatParticipants(txCtx, chatID); err != nil {
			return err
		}

		//Deleting chat messages
		if err := cs.chatMessageRepo.DeleteAllChatMessagesByChatID(txCtx, chatID); err != nil {
			return err
		}

		//Deleting chat
		if err := cs.chatRepo.DeleteChat(txCtx, chatID); err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		return err
	}

	return nil
}



func (cs *ChatService) CreateChatOwner(ctx context.Context, chatID uuid.UUID, userID uuid.UUID) error {
	chatOwner := chatparticipant.NewChatParticipant(chatID, userID, roles.OwnerChatRole.GetID(), time.Now())
	err := cs.chatParticipantRepo.AddChatParticipant(ctx, chatOwner)
	if err != nil {
		return err
	}
	return nil
}

func (cs *ChatService) IsOwner(ctx context.Context, chatID uuid.UUID, userID uuid.UUID) (bool, error) {
	chatParticipant, err := cs.chatParticipantRepo.GetChatParticipantByIDs(ctx, chatID, userID)
	if err != nil {
		return false, err
	}
	return chatParticipant.GetRoleID() == roles.OwnerChatRole.GetID(), nil
}



	
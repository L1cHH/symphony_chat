package service

import (
	"context"
	"slices"
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

	isEnoughPermissions, err := cs.IsUserHasEnoughPermissions(ctx, chatID, deletingInitiatorID, roles.PermissionDeleteChat)
	if err != nil {
		return err
	}

	if !isEnoughPermissions {
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

func (cs *ChatService) RenameChat(ctx context.Context, chatID uuid.UUID, newName string, renamingInitiatorID uuid.UUID) (string, error) {

	isEnoughPermissions, err := cs.IsUserHasEnoughPermissions(ctx, chatID, renamingInitiatorID, roles.PermissionUpdateChatName)
	if err != nil {
		return "", err
	}

	if !isEnoughPermissions {
		return "", roles.ErrInsufficientPermissions
	}

	err = cs.transactionManager.WithinTransaction(ctx, func(txCtx context.Context) error {
		if err := cs.chatRepo.UpdateChatName(txCtx, chatID, newName); err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		return "", err
	}

	return newName, nil
}

func (cs *ChatService) AddUserToChat(ctx context.Context, chatID uuid.UUID, inviterUserID uuid.UUID, invitedUserID uuid.UUID) error {
	isEnoughPermissions, err := cs.IsUserHasEnoughPermissions(ctx, chatID, inviterUserID, roles.PermissionAddMember)
	if err != nil {
		return err
	}

	if !isEnoughPermissions {
		return roles.ErrInsufficientPermissions
	}

	err = cs.transactionManager.WithinTransaction(ctx, func(txCtx context.Context) error {
		if err := cs.CreateChatMember(txCtx, chatID, invitedUserID); err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		return err
	}



	return nil 
}

func (cs *ChatService) RemoveUserFromChat(ctx context.Context, chatID uuid.UUID, removerUserID uuid.UUID, removedUserID uuid.UUID) error {
	isEnoughPermissions, err := cs.IsUserHasEnoughPermissions(ctx, chatID, removerUserID, roles.PermissionRemoveMember)
	if err != nil {
		return err
	}

	if !isEnoughPermissions {
		return roles.ErrInsufficientPermissions
	}

	err = cs.transactionManager.WithinTransaction(ctx, func(txCtx context.Context) error {
		if err := cs.chatParticipantRepo.DeleteChatParticipant(txCtx, chatID, removedUserID); err != nil {
			return err
		}

		//We are not deleting messages of the removed user from the chat
		return nil
	})

	if err != nil {
		return err
	}

	return nil
}

func (cs *ChatService) PromoteUserToChatAdmin(ctx context.Context, chatID uuid.UUID, promoterUserID uuid.UUID, promotedUserID uuid.UUID) error {
	isEnoughPermissions, err := cs.IsUserHasEnoughPermissions(ctx, chatID, promoterUserID, roles.PermissionManageRoles)
	if err != nil {
		return err
	}

	if !isEnoughPermissions {
		return roles.ErrInsufficientPermissions
	}

	err = cs.transactionManager.WithinTransaction(ctx, func(txCtx context.Context) error {
		if err := cs.chatParticipantRepo.UpdateChatParticipantRole(txCtx, chatID, promotedUserID, roles.AdminChatRole.GetID()); err != nil {
			return err
		}
		return nil
	})

	if err != nil {
		return err
	}

	return nil
}

func (cs *ChatService) DemoteChatAdminToChatMember(ctx context.Context, chatID uuid.UUID, demoterUserID uuid.UUID, adminUserID uuid.UUID) error {
	isEnoughPermissions, err := cs.IsUserHasEnoughPermissions(ctx, chatID, demoterUserID, roles.PermissionManageRoles)
	if err != nil {
		return err
	}

	if !isEnoughPermissions {
		return roles.ErrInsufficientPermissions
	}

	err = cs.transactionManager.WithinTransaction(ctx, func(txCtx context.Context) error {
		if err := cs.chatParticipantRepo.UpdateChatParticipantRole(txCtx, chatID, adminUserID, roles.MemberChatRole.GetID()); err != nil {
			return err
		}
		return nil
	})

	if err != nil {
		return err
	}

	return nil
}

func (cs *ChatService) LeaveChat(ctx context.Context, chatID uuid.UUID, leavingUserID uuid.UUID) error {
	err := cs.transactionManager.WithinTransaction(ctx, func(txCtx context.Context) error {
		return cs.chatParticipantRepo.DeleteChatParticipant(txCtx, chatID, leavingUserID)
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

func (cs *ChatService) CreateChatMember(ctx context.Context, chatID uuid.UUID, userID uuid.UUID) error {
	chatMember := chatparticipant.NewChatParticipant(chatID, userID, roles.MemberChatRole.GetID(), time.Now())
	err := cs.chatParticipantRepo.AddChatParticipant(ctx, chatMember)
	if err != nil {
		return err
	}
	return nil
}

func (cs *ChatService) IsUserHasEnoughPermissions(ctx context.Context, chatID uuid.UUID, userID uuid.UUID, requiredPermissions ...roles.Permission) (bool, error) {
	chatParticipant, err := cs.chatParticipantRepo.GetChatParticipantByIDs(ctx, chatID, userID)
	if err != nil {
		return false, err
	}

	chatRole, err := cs.chatRolesRepo.GetChatRoleByID(ctx, chatParticipant.GetRoleID())
	if err != nil {
		return false, err
	}

	userPermissions := chatRole.GetPermissions()

	for _, requiredPermission := range requiredPermissions {
		if !slices.Contains(userPermissions, requiredPermission) {
			return false, nil
		}
	}

	return true, nil
}




	
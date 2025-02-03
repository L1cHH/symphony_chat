package chatparticipant

import (
	"context"
	"time"

	"github.com/google/uuid"
)

type ChatParticipant struct {
	chatID uuid.UUID
	userID uuid.UUID
	roleID uuid.UUID
	joinedAt time.Time
}

func (c ChatParticipant) GetChatID() uuid.UUID {
	return c.chatID
}

func (c ChatParticipant) GetUserID() uuid.UUID {
	return c.userID
}

func (c ChatParticipant) GetRoleID() uuid.UUID {
	return c.roleID
}

func (c ChatParticipant) GetJoinedAt() time.Time {
	return c.joinedAt
}

func NewChatParticipant(chatID uuid.UUID, userID uuid.UUID, roleID uuid.UUID, joinedAt time.Time) ChatParticipant {
	return ChatParticipant {
		chatID: chatID,
		userID: userID,
		roleID: roleID,
		joinedAt: joinedAt,
	}
}

func ChatParticipantFromDB(chatID uuid.UUID, userID uuid.UUID, roleID uuid.UUID, joinedAt time.Time) ChatParticipant {
	return ChatParticipant {
		chatID: chatID,
		userID: userID,
		roleID: roleID,
		joinedAt: joinedAt,
	}
}

type ChatParticipantRepository interface {
	GetChatParticipantByIDs(ctx context.Context, chatID uuid.UUID, userID uuid.UUID) (ChatParticipant, error)
	GetAllChatParticipantsByChatID(ctx context.Context, chatID uuid.UUID) ([]ChatParticipant, error)

	AddChatParticipant(ctx context.Context, chatParticipant ChatParticipant) error

	UpdateChatParticipantRole(ctx context.Context, chatID uuid.UUID, userID uuid.UUID, newRoleID uuid.UUID) error
	
	DeleteChatParticipant(ctx context.Context, chatID uuid.UUID, userID uuid.UUID) error
	DeleteAllChatParticipants(ctx context.Context, chatID uuid.UUID) error
}
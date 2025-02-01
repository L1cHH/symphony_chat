package users

import (
	"context"
	"time"

	"github.com/google/uuid"
)

type ChatUser struct {
	id         uuid.UUID
	username   string
	status     UserStatus
	createdAt  time.Time
	lastSeenAt time.Time
}

type UserStatus string

const (
	Online  UserStatus = "online"
	Offline UserStatus = "offline"
)

func (ch ChatUser) GetID() uuid.UUID {
	return ch.id
}

func (ch ChatUser) GetUsername() string {
	return ch.username
}

func (ch ChatUser) GetStatus() UserStatus {
	return ch.status
}

func (ch ChatUser) GetCreatedAt() time.Time {
	return ch.createdAt
}

func (ch ChatUser) GetLastSeenAt() time.Time {
	return ch.lastSeenAt
}

func NewChatUser(authUSerID uuid.UUID, username string, status UserStatus, createdAt time.Time, lastSeenAt time.Time) ChatUser {
	return ChatUser{
		id:         authUSerID,
		username:   username,
		status:     status,
		createdAt:  createdAt,
		lastSeenAt: lastSeenAt,
	}
}

func ChatUserFromDB(id uuid.UUID, username string, status UserStatus, createdAt time.Time, lastSeenAt time.Time) ChatUser {
	return ChatUser{
		id:         id,
		username:   username,
		status:     status,
		createdAt:  createdAt,
		lastSeenAt: lastSeenAt,
	}
}

type ChatUserRepository interface {
	GetChatUserByID(ctx context.Context, chatUserId uuid.UUID) (ChatUser, error)
	GetChatUserByUsername(ctx context.Context, username string) (ChatUser, error)
	AddChatUser(ctx context.Context, chatUser ChatUser) error
	DeleteChatUserByID(ctx context.Context, chatUserId uuid.UUID) error
	UpdateUsername(ctx context.Context, chatUserId uuid.UUID, newUsername string) error
	UpdateStatus(ctx context.Context, chatUserId uuid.UUID, newStatus UserStatus) error
	UpdateLastSeenAt(ctx context.Context, chatUserId uuid.UUID, newLastSeenAt time.Time) error
}

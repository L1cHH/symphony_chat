package users

import (
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

func NewChatUser(authUSerID uuid.UUID, username string) ChatUser {
	return ChatUser{
		id:         authUSerID,
		username:   username,
		status:     Offline,
		createdAt:  time.Now(),
		lastSeenAt: time.Now(),
	}
}

type ChatUserRepository interface {
	GetChatUserByID(uuid.UUID) (ChatUser, error)
	GetChatUserByUsername(string) (ChatUser, error)
	AddChatUser(ChatUser) error
	DeleteChatUserByID(uuid.UUID) error
	UpdateChatUser(uuid.UUID) error
}

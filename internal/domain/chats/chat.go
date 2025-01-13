package chats

import (
	"time"

	"github.com/google/uuid"
)

type Chat struct {
	id        uuid.UUID
	userOneID uuid.UUID
	userTwoID uuid.UUID
	createdAt time.Time
}

func (c Chat) GetID() uuid.UUID {
	return c.id
}

func (c Chat) GetUserOneID() uuid.UUID {
	return c.userOneID
}

func (c Chat) GetUserTwoID() uuid.UUID {
	return c.userTwoID
}

func (c Chat) GetCreatedAt() time.Time {
	return c.createdAt
}

func NewChat(userOneID uuid.UUID, userTwoID uuid.UUID) Chat {
	return Chat{
		id:        uuid.New(),
		userOneID: userOneID,
		userTwoID: userTwoID,
		createdAt: time.Now(),
	}
}

type ChatRepository interface {
	GetChatByID(uuid.UUID) (Chat, error)
	UpdateChat(uuid.UUID) error
	AddChat(Chat) error
	DeleteChat(uuid.UUID) error
}

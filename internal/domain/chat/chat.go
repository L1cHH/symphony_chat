package chat

import (
	"context"
	"time"

	"github.com/google/uuid"
)

type Chat struct {
	id 			uuid.UUID
	name		string
	createdAt   time.Time
	updatedAt   time.Time
}


func (c Chat) GetID() uuid.UUID {
	return c.id 
}

func (c Chat) GetName() string {
	return c.name
}

func (c Chat) GetCreatedAt() time.Time {
	return c.createdAt
}

func (c Chat) GetUpdatedAt() time.Time {
	return c.updatedAt
}

func NewChat(name string) (Chat, error) {
	if len(name) == 0 {
		return Chat{}, ErrWrongChatName
	}

	if len(name) > 15 {
		return Chat{}, ErrWrongChatName
	}
	
	return Chat {
		id: uuid.New(),
		name: name,
		createdAt: time.Now(),
		updatedAt: time.Now(),
	}, nil
}

func ChatFromDB(id uuid.UUID, name string, createdAt time.Time, updatedAt time.Time) Chat {
	return Chat {
		id: id, 
		name: name,
		createdAt: createdAt,
		updatedAt: updatedAt,
	}
}

type ChatRepository interface {
	GetChatByID(context.Context, uuid.UUID) (Chat, error)
	AddChat(context.Context, Chat) error
	UpdateChatName(context.Context, uuid.UUID, string) error
	UpdateChatUpdatedAt(ctx context.Context, chatID uuid.UUID, updatedAt time.Time) error
	DeleteChat(context.Context, uuid.UUID) error
}
package messages

import (
	"context"
	"time"

	"github.com/google/uuid"
)

type ChatMessage struct {
	id uuid.UUID
	chatID uuid.UUID
	senderID uuid.UUID
	content string
	createdAt time.Time
	status MessageStatus
}

type MessageStatus string 

const (
	Sent MessageStatus = "sent"
	Received MessageStatus = "received"
	Edited MessageStatus = "edited" 
	Read MessageStatus = "read"
)

func (cm ChatMessage) GetID() uuid.UUID {
	return cm.id
}

func (cm ChatMessage) GetChatID() uuid.UUID {
	return cm.chatID
}

func (cm ChatMessage) GetSenderID() uuid.UUID {
	return cm.senderID
}

func (cm ChatMessage) GetContent() string {
	return cm.content
}

func (cm ChatMessage) GetCreatedAt() time.Time {
	return cm.createdAt
}

func (cm ChatMessage) GetStatus() MessageStatus {
	return cm.status
}

func NewChatMessage(chatID uuid.UUID, senderID uuid.UUID, content string, createdAt time.Time, status MessageStatus) ChatMessage {
	return ChatMessage{
		id: uuid.New(),
		chatID: chatID,
		senderID: senderID,
		content: content,
		createdAt: createdAt,
		status: status,
	}
}

func ChatMessageFromDB(id uuid.UUID, chatID uuid.UUID, senderID uuid.UUID, content string, createdAt time.Time, status MessageStatus) ChatMessage {
	return ChatMessage{
		id: id,
		chatID: chatID,
		senderID: senderID,
		content: content,
		createdAt: createdAt,
		status: status,
	}
}

type ChatMessageRepository interface {
	GetChatMessageById(ctx context.Context, messageID uuid.UUID) (ChatMessage, error)
	GetChatMessagesByChatId(ctx context.Context, chatID uuid.UUID) ([]ChatMessage, error)
	GetChatMessagesByContentAndChatID(ctx context.Context, content string, chatID uuid.UUID) ([]ChatMessage, error)

	AddChatMessage(ctx context.Context, message ChatMessage) error

	UpdateChatMessageContent(ctx context.Context, messageID uuid.UUID, content string) error
	UpdateChatMessageStatus(ctx context.Context, messageID uuid.UUID, status MessageStatus) error
	
	DeleteChatMessage(ctx context.Context, messageID uuid.UUID) error
	DeleteAllChatMessagesByChatID(ctx context.Context, chatID uuid.UUID) error
}
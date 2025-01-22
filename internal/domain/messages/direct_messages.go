package messages

import (
	"time"

	"github.com/google/uuid"
)

type DirectMessage struct {
	id          uuid.UUID
	chatID      uuid.UUID
	senderID    uuid.UUID
	recipientID uuid.UUID
	message     string
	createdAt   time.Time
	status      MessageStatus
	isEdited    bool
}

type MessageStatus string

const (
	Sent        MessageStatus = "sent"
	Read        MessageStatus = "read"
	Undelivered MessageStatus = "undelivered"
)

func (dm DirectMessage) GetID() uuid.UUID {
	return dm.id
}

func (dm DirectMessage) GetChatID() uuid.UUID {
	return dm.chatID
}

func (dm DirectMessage) GetSenderID() uuid.UUID {
	return dm.senderID
}

func (dm DirectMessage) GetRecipientID() uuid.UUID {
	return dm.recipientID
}

func (dm DirectMessage) GetMessageContent() string {
	return dm.message
}

func (dm DirectMessage) GetCreatedAt() time.Time {
	return dm.createdAt
}

func (dm DirectMessage) GetIsEdited() bool {
	return dm.isEdited
}

func (dm DirectMessage) GetStatus() MessageStatus {
	return dm.status
}

func NewDirectMessage(chatID uuid.UUID, senderID uuid.UUID, recipientID uuid.UUID, message string) DirectMessage {
	return DirectMessage{
		id:          uuid.New(),
		chatID:      chatID,
		senderID:    senderID,
		recipientID: recipientID,
		message:     message,
		createdAt:   time.Now(),
		status:      Sent,
		isEdited:    false,
	}
}

type DirectMessageRepository interface {
	//Get-methods
	GetDirectMessageByID(directMessageID uuid.UUID) (DirectMessage, error)
	GetDirectMessagesBySenderID(senderID uuid.UUID, chatID uuid.UUID) ([]DirectMessage, error)
	GetDirectMessagesByRecipientID(recipientID uuid.UUID, chatID uuid.UUID) ([]DirectMessage, error)
	GetDirectMessagesByChatID(chatID uuid.UUID) ([]DirectMessage, error)
	//Create-methods
	AddDirectMessage(DirectMessage) error
	//Delete-methods
	DeleteDirectMessage(directMessageID uuid.UUID) error
	DeleteDirectMessagesByChatID(chatID uuid.UUID) error
	DeleteDirectMessagesBySenderID(chatID uuid.UUID, senderID uuid.UUID) error
	//Update-methods
	UpdateDirectMessage(directMessageID uuid.UUID, updatedMessage string) error
}

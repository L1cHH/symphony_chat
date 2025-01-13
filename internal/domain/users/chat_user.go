package users

import "github.com/google/uuid"

type ChatUser struct {
	id uuid.UUID
	username string

}

func (ch ChatUser) GetID() uuid.UUID {
	return ch.id
}

func (ch ChatUser) GetUsername() string {
	return ch.username
}

func NewChatUser(username string) ChatUser {
	return ChatUser{
		id: uuid.New(),
		username: username, 
	}
}

type ChatUserRepository interface {
	GetChatUserByID(uuid.UUID) (ChatUser, error) 
	AddChatUser(ChatUser) error
	DeleteChatUserByID(uuid.UUID) error 
	UpdateChatUser(uuid.UUID) error
}


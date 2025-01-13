package chats

import (
	"errors"
	"symphony_chat/internal/domain/chats"
	"symphony_chat/internal/domain/messages"
	"symphony_chat/internal/domain/users"

	"github.com/google/uuid"
)

var (
	//Problem with creating a chat
	ErrProblemWithCreatingChat = errors.New("error occurs while creating chat")
	//Problem with creating chatUser
	ErrProblemWithCreatingChatUser = errors.New("error occurs while creating chatUser")
)

type ChatService struct {
	chatUserRepo users.ChatUserRepository
	directMsgRepo messages.DirectMessageRepository
	chatRepo chats.ChatRepository
}

type ChatConfiguration func(cs *ChatService) error


func NewChatService(configs ...ChatConfiguration) (*ChatService, error) {
	cs := &ChatService{}
	
	for _, cfg := range configs {
		err := cfg(cs)
		if err != nil {
			return nil, err
		}
	}

	return cs, nil
}


func (cs *ChatService) CreateChatUser(authUserID uuid.UUID, username string) (users.ChatUser, error) {
	chatUser := users.NewChatUser(authUserID, username)

	err := cs.chatUserRepo.AddChatUser(chatUser)
	if err != nil {
		return users.ChatUser{}, errors.New(ErrProblemWithCreatingChatUser.Error() + ": " + err.Error())
	}

	return chatUser, nil
}

func (cs *ChatService) CreateChat(userOneID uuid.UUID, userTwoID uuid.UUID) (chats.Chat, error) {
	chat := chats.NewChat(userOneID, userTwoID)
	err := cs.chatRepo.AddChat(chat)
	if err != nil {
		return chats.Chat{}, errors.New(ErrProblemWithCreatingChat.Error() + ": " + err.Error())
	}

	return chat, nil
}



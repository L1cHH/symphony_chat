package chats

import (
	"errors"
	"fmt"
	"symphony_chat/internal/domain/chats"
	"symphony_chat/internal/domain/messages"
	"symphony_chat/internal/domain/users"

	"github.com/google/uuid"
)

var (
	//Problem with creating a chat
	ErrWithCreatingChat = errors.New("error occurs while creating chat")
	//Problem with creating chatUser
	ErrWithCreatingChatUser = errors.New("error occurs while creating chatUser")
	//Problem with creating direct message
	ErrWithCreatingDirectMessage = errors.New("error occurs while creating direct message")
	//Problem with sending direct message
	ErrWithSendingDirectMessage = errors.New("error occurs while sending direct message")
	//Empty message
	ErrEmptyMessage = errors.New("message is empty")
	//User Ids are the same
	ErrIncorrectUserIds = errors.New("user ids are the same")
)

type ChatService struct {
	chatUserRepo  users.ChatUserRepository
	directMsgRepo messages.DirectMessageRepository
	chatRepo      chats.ChatRepository
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
		return users.ChatUser{}, errors.New(ErrWithCreatingChatUser.Error() + ": " + err.Error())
	}

	return chatUser, nil
}

func (cs *ChatService) CreateChat(userOneID uuid.UUID, userTwoID uuid.UUID) (chats.Chat, error) {
	chat := chats.NewChat(userOneID, userTwoID)
	err := cs.chatRepo.AddChat(chat)
	if err != nil {
		return chats.Chat{}, errors.New(ErrWithCreatingChat.Error() + ": " + err.Error())
	}

	return chat, nil
}

func (cs *ChatService) CreateDirectMessage(chatID uuid.UUID, senderID uuid.UUID, receiverID uuid.UUID, text string) (messages.DirectMessage, error) {

	dm := messages.NewDirectMessage(chatID, senderID, receiverID, text)

	err := cs.directMsgRepo.AddDirectMessage(dm)
	if err != nil {
		return messages.DirectMessage{}, fmt.Errorf("%w: %v", ErrWithCreatingDirectMessage, err)
	}

	return dm, nil
}

func (cs *ChatService) SendDirectMessage(senderID uuid.UUID, recipientID uuid.UUID, messageContent string) error {

	if messageContent == "" {
		return ErrEmptyMessage
	}

	if senderID == recipientID {
		return ErrIncorrectUserIds
	}

	_, err := cs.chatUserRepo.GetChatUserByID(senderID)
	if err != nil {
		return fmt.Errorf("%w: %v", ErrWithSendingDirectMessage, err)
	}

	_, err = cs.chatUserRepo.GetChatUserByID(recipientID)
	if err != nil {
		return fmt.Errorf("%w: %v", ErrWithSendingDirectMessage, err)
	}

	chat, err := cs.chatRepo.GetChatByUsersID(senderID, recipientID)
	if err != nil {
		switch {
		case errors.Is(err, chats.ErrChatDoesNotExist):
			chat, err = cs.CreateChat(senderID, recipientID)
			if err != nil {
				return fmt.Errorf("%w: %v", ErrWithSendingDirectMessage, err)
			}
		default:
			return fmt.Errorf("%w: %v", ErrWithSendingDirectMessage, err)
		}
	}

	_, err = cs.CreateDirectMessage(chat.GetID(), senderID, recipientID, messageContent)
	if err != nil {
		return fmt.Errorf("%w: %v", ErrWithSendingDirectMessage, err)
	}

	return nil
}

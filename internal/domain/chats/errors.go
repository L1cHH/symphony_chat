package chats

import "errors"

var (
	ErrProblemWithDB    = errors.New("error occurs while working with db")
	ErrChatDoesNotExist = errors.New("chat does not exist")
)

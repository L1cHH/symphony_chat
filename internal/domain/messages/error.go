package messages 

type ChatMessageError struct {
	Code  string
	Message string
	Err error
}

func (e *ChatMessageError) Error() string {
	if e.Err != nil {
		return e.Message + ": " + e.Err.Error()
	}
	return e.Message
}

var (
	ErrChatMessageNotFound = &ChatMessageError {
		Code: "CHAT_MESSAGE_NOT_FOUND",
		Message: "chat message with that id not found in storage",
	}
)
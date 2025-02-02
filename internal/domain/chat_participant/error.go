package chatparticipant

type ChatParticipantError struct {
	Code       string
	Message    string
	Err        error
}

func (e *ChatParticipantError) Error() string {
	if e.Err != nil {
		return e.Message + ": " + e.Err.Error()
	}
	return e.Message
}

var (
	ErrChatParticipantNotFound = &ChatParticipantError{
		Code: "CHAT_PARTICIPANT_NOT_FOUND",
		Message: "chat participant not found",
	}
)
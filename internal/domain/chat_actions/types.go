package actions

type ChatActionType string

const (
	//Clients actions

	//Chat actions
	LeaveChatAction ChatActionType = "LEAVE_CHAT"
	CreateChatAction ChatActionType = "CREATE_CHAT"
	RenameChatAction ChatActionType = "RENAME_CHAT"
	DeleteChatAction ChatActionType = "DELETE_CHAT"

	//Members actions
	AddMemberToChatAction ChatActionType = "ADD_MEMBER_TO_CHAT"
	RemoveMemberFromChatAction ChatActionType = "REMOVE_MEMBER_FROM_CHAT"
	PromoteUserToChatAdminAction ChatActionType = "PROMOTE_USER_TO_CHAT_ADMIN"
	DemoteChatAdminToChatMemberAction ChatActionType = "DEMOTE_CHAT_ADMIN_TO_CHAT_MEMBER"

	//Messages actions
	SendMessageAction ChatActionType = "SEND_MESSAGE"
	DeleteMessageAction ChatActionType = "DELETE_MESSAGE"
	EditMessageAction ChatActionType = "EDIT_MESSAGE"

)

type ChatActionResult string

const (
	Success ChatActionResult = "SUCCESS"
	Failed ChatActionResult = "FAILED"
)

type EventType string 

const (
	UserEnteredChatEvent EventType = "USER_ENTERED_CHAT"
	UserLeftChatEvent EventType = "USER_LEFT_CHAT"
	ChatNameUpdatedEvent EventType = "CHAT_NAME_UPDATED"
	UserSentMessageEvent EventType = "USER_SENT_MESSAGE"
	UserEditedMessageEvent EventType = "USER_EDITED_MESSAGE"
	UserDeletedMessageEvent EventType = "USER_DELETED_MESSAGE"
	UserWasKickedFromChatEvent EventType = "USER_WAS_KICKED_FROM_CHAT"
	UserWasPromotedToChatAdminEvent EventType = "USER_WAS_PROMOTED_TO_CHAT_ADMIN"
	UserWasDemotedFromChatAdminEvent EventType = "USER_WAS_DEMOTED_FROM_CHAT_ADMIN"
)
package constant

import "time"

const (
	// 一天
	UserInboxExpiredTime = time.Hour * 24

	UserInboxKey = "inbox:%d" // uid
)

const (
	InboxTypeUser         = 0
	InboxTypeConversation = 1
)

const (
	ConversationIdPatternSingle = "0:%d:%d" // 私聊
)
package util

import (
	"strconv"
	"strings"

	"github.com/ChowRobin/fantim/constant"
)

func GetReceiver(conversationType int32, conversationId string, sender int64) int64 {
	strSlice := strings.Split(conversationId, ":")
	if len(strSlice) < 3 {
		return 0
	}

	switch conversationType {
	case constant.ConversationTypeSingle:
		user1, user2 := strSlice[1], strSlice[2]
		userId1, _ := strconv.ParseInt(user1, 10, 64)
		userId2, _ := strconv.ParseInt(user2, 10, 64)
		if userId1 == sender {
			return userId2
		} else {
			return userId1
		}
	case constant.ConversationTypeGroup:
		groupId, _ := strconv.ParseInt(strSlice[2], 10, 64)
		return groupId
	}

}

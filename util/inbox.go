package util

import (
	"strconv"
	"strings"
)

func GetReceiver(conversationId string, sender int64) int64 {
	strSlice := strings.Split(conversationId, ":")
	if len(strSlice) < 3 {
		return 0
	}
	user1, user2 := strSlice[1], strSlice[2]
	userId1, _ := strconv.ParseInt(user1, 10, 64)
	userId2, _ := strconv.ParseInt(user2, 10, 64)
	if userId1 == sender {
		return userId2
	} else {
		return userId1
	}
}

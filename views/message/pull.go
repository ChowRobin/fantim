package message

import (
	"fmt"
	"log"
	"strconv"

	"github.com/ChowRobin/fantim/constant"
	"github.com/ChowRobin/fantim/model/bo"

	"github.com/ChowRobin/fantim/constant/status"
	"github.com/ChowRobin/fantim/model/vo"
	"github.com/gin-gonic/gin"
)

func Pull(c *gin.Context) interface{} {
	req := &vo.MessagePullRequest{}
	resp := &vo.MessagePullResponse{}

	userId := c.GetInt64("user_id")
	// 校验参数合法性
	err := c.BindJSON(req)
	if err != nil {
		log.Printf("[Pull] BindJSON failed. err=%v", err)
		return status.FillResp(resp, status.ErrInvalidParam)
	}

	inbox := &bo.Inbox{
		InboxType: int8(req.InboxType),
	}
	switch req.InboxType {
	case constant.InboxTypeUser:
		inbox.Key = fmt.Sprintf(constant.UserInboxKey, userId)
	case constant.InboxTypeConversation:
		inbox.Key = req.ConversationId
		msgId, err := strconv.ParseInt(req.MsgIdStr, 10, 64)
		if err != nil {
			return status.FillResp(resp, status.ErrInvalidParam)
		}
		req.Cursor = msgId
	}

	msgList, err := inbox.Pull(req.Cursor, int64(req.Count))
	if err != nil {
		log.Printf("[message.Pull] inbox.Pull failed. err=%v", err)
		return status.FillResp(resp, status.ErrServiceInternal)
	}

	resp.MessageList = msgList

	if req.Count == -1 {
		resp.HasMore = false
	} else {
		resp.HasMore = int32(len(msgList)) == req.Count
	}

	return status.FillResp(resp, status.Success)
}

package message

import (
	"log"

	"github.com/ChowRobin/fantim/constant/status"
	"github.com/ChowRobin/fantim/model/vo"
	"github.com/ChowRobin/fantim/service"
	"github.com/gin-gonic/gin"
)

func Send(c *gin.Context) interface{} {
	req := &vo.MessageSendRequest{}
	resp := &vo.MessageSendResponse{}

	userId := c.GetInt64("user_id")

	// 校验参数合法性
	err := c.BindJSON(req)
	if err != nil {
		log.Printf("[Send] BindJSON failed. err=%v", err)
		return status.FillResp(resp, status.ErrInvalidParam)
	}
	if req.MsgType == 0 || req.Content == "" {
		log.Printf("[Send] invalid param. req=%v", req)
		return status.FillResp(resp, status.ErrInvalidParam)
	}
	// todo 校验能否发消息

	msgVo := &vo.MessageBody{
		ConversationType: req.ConversationType,
		ConversationId:   req.ConversationId,
		MsgType:          req.MsgType,
		Content:          req.Content,
		Ext:              req.Ext,
		Sender:           userId,
	}
	msgId, err := service.SendMessage(c, msgVo)
	if err != nil {
		log.Printf("[Send] service.SendMessage failed. err=%v", err)
		return status.FillResp(resp, status.ErrServiceInternal)
	}

	resp.MsgId = msgId
	return status.FillResp(resp, status.Success)
}

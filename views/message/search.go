package message

import (
	"log"
	"strconv"

	"github.com/ChowRobin/fantim/constant/status"
	"github.com/ChowRobin/fantim/model/bo"
	"github.com/ChowRobin/fantim/model/po"
	"github.com/ChowRobin/fantim/model/vo"
	"github.com/gin-gonic/gin"
)

func Search(c *gin.Context) interface{} {
	resp := &vo.MessageSearchResponse{}

	key := c.Query("key")
	cursor := c.Query("cursor")
	convId := c.Query("conversation_id")
	count, _ := strconv.Atoi(c.Query("count"))

	if count == 0 || key == "" {
		return status.FillResp(resp, status.ErrInvalidParam)
	}

	midx, newCursor, err := bo.SearchMessageIndex(c, convId, key, cursor, count)
	if err != nil {
		log.Printf("[Message.Search] SearchMessageIndex failed. err=%v", err)
		return status.FillResp(resp, status.ErrServiceInternal)
	}
	msgIdList := make([]int64, 0, len(midx))
	for _, m := range midx {
		msgIdList = append(msgIdList, m.MsgId)
	}
	if len(msgIdList) > 0 {
		mPos, err := po.MultiGetMessage(c, convId, msgIdList)
		if err != nil {
			log.Printf("[Message.Search] po.MultiGetMessage failed. err=%v", err)
			return status.FillResp(resp, status.ErrServiceInternal)
		}
		resp.MessageList = po.MessagePoListToVo(mPos)
		resp.Cursor = newCursor
	}

	return status.FillResp(resp, status.Success)
}

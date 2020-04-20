package bo

import (
	"context"
	"errors"
	"fmt"

	"github.com/ChowRobin/fantim/constant"
	"github.com/ChowRobin/fantim/model/po"
	"github.com/ChowRobin/fantim/model/vo"
	"github.com/ChowRobin/fantim/util"
)

type Inbox struct {
	Ctx       context.Context
	InboxType int8   // 0->用户链，1->会话链
	Key       string // 0->redis key 1->conversation_id
}

func (b *Inbox) Append(m *vo.MessageBody) (int64, error) {
	if b.Key == "" {
		return 0, errors.New("inbox key is nil")
	}
	switch b.InboxType {
	case constant.InboxTypeUser:
		msg := &UserMessage{*m}
		return msg.Add(b.Key)
	case constant.InboxTypeConversation:
		msgId := util.GenId()
		msgPo := po.MessageRecord{
			MsgId:            m.MsgId,
			Sender:           m.Sender,
			ConversationId:   m.ConversationId,
			ConversationType: int8(m.ConversationType),
			MsgType:          m.MsgType,
			Content:          m.Content,
			Ext:              util.ToJsonString(m.Ext),
		}
		err := msgPo.Create(b.Ctx)
		return msgId, err
	default:
		return 0, fmt.Errorf("invalid inbox_type %d", b.InboxType)
	}
}

func (b *Inbox) Pull(cursor, count int64) ([]*vo.MessageBody, error) {
	switch b.InboxType {
	case constant.InboxTypeUser:
		var start, stop int64
		if cursor == -1 && count == -1 { // 拉取全部
			start = 0
			stop = -1
		} else if cursor > 0 && count < 0 { // 逆向拉历史消息
			start = cursor + count + 1
			stop = cursor
		} else if cursor >= 0 && count > 0 { // 正向拉新消息
			start = cursor
			stop = cursor + count - 1
		} else {
			return nil, fmt.Errorf("invalid cursor %d count %d", cursor, count)
		}
		msgList, err := PullMessage(b.Key, start, stop)
		return UserMessageListToVo(msgList), err
	case constant.InboxTypeConversation:
		msgPoList, err := po.ListByConversationAndMsgId(b.Ctx, b.Key, cursor, count)
		if err != nil {
			return nil, err
		}
		return po.MessagePoListToVo(msgPoList), err
	default:
		return nil, fmt.Errorf("invalid inbox_type %d", b.InboxType)
	}
}

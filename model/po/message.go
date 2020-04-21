package po

import (
	"context"
	"encoding/json"
	"strconv"
	"time"

	"github.com/ChowRobin/fantim/model/vo"

	"github.com/ChowRobin/fantim/client"
	"github.com/jinzhu/gorm"
)

type MessageRecord struct {
	Id     int64 `gorm:"primary_key"`
	MsgId  int64 `gorm:"column:msg_id"`
	Sender int64 `gorm:"column:sender"`

	ConversationId   string `gorm:"column:conversation_id"`
	ConversationType int8   `gorm:"column:conversation_type"`
	MsgType          int32  `gorm:"column:msg_type"`
	Content          string `gorm:"column:content"`
	Ext              string `gorm:"column:ext"`

	CreateTime *time.Time `gorm:"column:create_time"`
	UpdateTime *time.Time `gorm:"column:update_time"`
}

func (*MessageRecord) TableName() string {
	return "im_message"
}

func (m *MessageRecord) Create(ctx context.Context) error {
	conn, err := client.DBConn(ctx)
	if err != nil {
		return err
	}

	defer client.CloseDBConn(ctx)
	return conn.Create(m).Error
}

func (m *MessageRecord) Update(ctx context.Context) error {
	conn, err := client.DBConn(ctx)
	if err != nil {
		return err
	}

	defer client.CloseDBConn(ctx)
	return conn.Model(m).Update(m).Error
}

func MessagePoListToVo(msgList []*MessageRecord) []*vo.MessageBody {
	result := make([]*vo.MessageBody, 0, len(msgList))
	for _, msg := range msgList {
		extMap := make(map[string]string)
		_ = json.Unmarshal([]byte(msg.Ext), &extMap)
		result = append(result, &vo.MessageBody{
			ConversationType: int32(msg.ConversationType),
			ConversationId:   msg.ConversationId,
			MsgType:          msg.MsgType,
			Content:          msg.Content,
			Ext:              extMap,
			CreateTime:       msg.CreateTime.Unix(),
			Sender:           msg.Sender,
			MsgId:            msg.MsgId,
			MsgIdStr:         strconv.FormatInt(msg.MsgId, 10),
		})
	}
	return result
}

func ListMessageByConversation(ctx context.Context, conversationId string, count int64) ([]*MessageRecord, error) {
	conn, err := client.DBConn(ctx)
	if err != nil {
		return nil, err
	}
	defer client.CloseDBConn(ctx)
	var list []*MessageRecord
	err = conn.Where("conversation_id=?", conversationId).
		Order("msg_id desc").Limit(count).Find(&list).Error
	if err == gorm.ErrRecordNotFound {
		return list, nil
	}
	return list, err
}

func ListByConversationAndMsgId(ctx context.Context, conversationId string, msgId, count int64) ([]*MessageRecord, error) {
	conn, err := client.DBConn(ctx)
	if err != nil {
		return nil, err
	}
	var list []*MessageRecord
	defer client.CloseDBConn(ctx)
	conn = conn.Debug()
	err = conn.Where("conversation_id=? and msg_id < ?", conversationId, msgId).Limit(count).
		Order("msg_id desc").Find(&list).Error
	if err == gorm.ErrRecordNotFound {
		return list, nil
	}
	return list, err
}

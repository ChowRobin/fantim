package service

import (
	"context"
	"errors"
	"fmt"
	"log"
	"strconv"
	"sync"
	"time"

	"github.com/ChowRobin/fantim/manager"

	"github.com/ChowRobin/fantim/constant"
	"github.com/ChowRobin/fantim/model/bo"
	"github.com/ChowRobin/fantim/model/vo"
	"github.com/ChowRobin/fantim/util"
)

func SendMessage(ctx context.Context, msg *vo.MessageBody) (msgId int64, err error) {
	var sender, receiver int64
	sender = msg.Sender
	receiver = util.GetReceiver(msg.ConversationId, sender)
	// 生成消息id
	msg.MsgId = util.GenId()
	if msg.MsgId == 0 {
		return 0, errors.New("msgId is 0")
	}
	msg.MsgIdStr = strconv.FormatInt(msg.MsgId, 10)
	msg.CreateTime = time.Now().Unix()

	// 写入会话链
	convInbox := &bo.Inbox{
		Ctx:       ctx,
		InboxType: constant.InboxTypeConversation,
		Key:       msg.ConversationId,
	}
	_, err = convInbox.Append(msg)
	if err != nil {
		return 0, err
	}

	// 写入用户链 兼容群聊需要抽出
	// 接收方用户链
	recvInbox := &bo.Inbox{
		Ctx:       ctx,
		InboxType: constant.InboxTypeUser,
		Key:       fmt.Sprintf(constant.UserInboxKey, sender),
	}
	recvIndex, err := recvInbox.Append(msg)
	log.Printf("[SendMessage] receiver %d inbox index=%d", sender, recvIndex)
	// 发送方用户链
	sendInbox := &bo.Inbox{
		Ctx:       ctx,
		InboxType: constant.InboxTypeUser,
		Key:       fmt.Sprintf(constant.UserInboxKey, sender),
	}
	sendIndex, err := sendInbox.Append(msg)
	log.Printf("[SendMessage] sender %d inbox index=%d", receiver, sendIndex)

	// 长链推通知
	wg := &sync.WaitGroup{}
	wg.Add(1)
	go func() {
		defer wg.Done()
		err = manager.PushMessage(sender, &vo.PushMessage{
			Body:  msg,
			Index: int32(sendIndex),
		})
	}()
	wg.Add(1)
	go func() {
		defer wg.Done()
		err = manager.PushMessage(receiver, &vo.PushMessage{
			Body:  msg,
			Index: int32(recvIndex),
		})
	}()
	wg.Wait()

	return msg.MsgId, nil
}

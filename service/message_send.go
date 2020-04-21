package service

import (
	"context"
	"fmt"
	"log"
	"strconv"
	"sync"
	"time"

	"github.com/ChowRobin/fantim/constant"
	"github.com/ChowRobin/fantim/constant/status"
	"github.com/ChowRobin/fantim/manager"
	"github.com/ChowRobin/fantim/model/bo"
	"github.com/ChowRobin/fantim/model/vo"
	"github.com/ChowRobin/fantim/util"
)

func SendMessage(ctx context.Context, msg *vo.MessageBody) (msgId int64, es *status.ErrStatus) {
	var sender, receiver int64
	sender = msg.Sender
	if msg.Receiver != 0 {
		receiver = msg.Receiver
		// 生成conversationId
		msg.ConversationId = util.GenConversationId(sender, receiver)
	} else if msg.ConversationId != "" {
		receiver = util.GetReceiver(msg.ConversationId, sender)

		msg.Receiver = receiver
	}

	if s := checkMsg(msg); s != status.Success {
		es = s
		return
	}

	// 生成消息id
	msg.MsgId = util.GenId()
	if msg.MsgId == 0 {
		es = status.ErrServiceInternal
		return
	}
	msg.MsgIdStr = strconv.FormatInt(msg.MsgId, 10)
	msg.CreateTime = time.Now().Unix()

	// 写入会话链
	convInbox := &bo.Inbox{
		Ctx:       ctx,
		InboxType: constant.InboxTypeConversation,
		Key:       msg.ConversationId,
	}
	_, err := convInbox.Append(msg)
	if err != nil {
		log.Printf("[service.SendMessage] convInbox append failed. err=%v", err)
		es = status.ErrServiceInternal
		return
	}

	// 写入用户链 兼容群聊需要抽出
	// 接收方用户链
	recvInbox := &bo.Inbox{
		Ctx:       ctx,
		InboxType: constant.InboxTypeUser,
		Key:       fmt.Sprintf(constant.UserInboxKey, receiver),
	}
	recvIndex, err := recvInbox.Append(msg)
	if err != nil {
		log.Printf("[service.SendMessage] recvIndex append failed. err=%v", err)
		es = status.ErrServiceInternal
		return
	}
	log.Printf("[SendMessage] receiver %d inbox index=%d", sender, recvIndex)
	// 发送方用户链
	sendInbox := &bo.Inbox{
		Ctx:       ctx,
		InboxType: constant.InboxTypeUser,
		Key:       fmt.Sprintf(constant.UserInboxKey, sender),
	}
	sendIndex, err := sendInbox.Append(msg)
	if err != nil {
		log.Printf("[service.SendMessage] sendIndex append failed. err=%v", err)
		es = status.ErrServiceInternal
		return
	}
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

	return msg.MsgId, status.Success
}

func checkMsg(msg *vo.MessageBody) *status.ErrStatus {
	if msg.ConversationType == 0 { // 私聊
		checkConvId := util.GenConversationId(msg.Sender, msg.Receiver)
		if checkConvId != msg.ConversationId {
			return status.ErrInvalidParam
		}
	}
	return status.Success
}

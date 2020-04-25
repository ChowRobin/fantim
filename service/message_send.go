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
	"github.com/ChowRobin/fantim/model/po"
	"github.com/ChowRobin/fantim/model/vo"
	"github.com/ChowRobin/fantim/util"
)

func SendMessage(ctx context.Context, msg *vo.MessageBody) (msgId int64, es *status.ErrStatus) {
	var sender, receiver int64
	sender = msg.Sender
	if msg.Receiver != 0 {
		receiver = msg.Receiver
		// 生成conversationId
		msg.ConversationId = util.GenConversationId(msg.ConversationType, sender, receiver)
	} else if msg.ConversationId != "" {
		receiver = util.GetReceiver(msg.ConversationType, msg.ConversationId, sender)

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

	if msg.ConversationType == constant.ConversationTypeSingle {
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
	} else if msg.ConversationType == constant.ConversationTypeGroup {
		err = HandleGroupMessage(ctx, msg)
		if err != nil {
			log.Printf("[SendMessage] HandleGroupMessage failed. err=%v", err)
			es = status.ErrServiceInternal
			return
		}
	}

	return msg.MsgId, status.Success
}

func checkMsg(msg *vo.MessageBody) *status.ErrStatus {
	checkConvId := util.GenConversationId(msg.ConversationType, msg.Sender, msg.Receiver)
	if msg.ConversationType == constant.ConversationTypeSingle {
		if checkConvId != msg.ConversationId {
			return status.ErrInvalidParam
		}
	}
	return status.Success
}

func HandleGroupMessage(ctx context.Context, msg *vo.MessageBody) error {
	// 查询群用户
	members, err := po.ListMembersByGroupId(ctx, msg.Receiver)
	if err != nil {
		return err
	}

	// 写入全部用户链 and 推送通知
	wg := &sync.WaitGroup{}
	for _, member := range members {
		wg.Add(1)
		go func(m *po.GroupMember) {
			defer wg.Done()
			inbox := &bo.Inbox{
				Ctx:       ctx,
				InboxType: constant.InboxTypeUser,
				Key:       fmt.Sprintf(constant.UserInboxKey, m.UserId),
			}
			index, err := inbox.Append(msg)
			if err != nil {
				log.Printf("[HandleGroupMessage] inbox.Append failed. err=%v userId=%d, msg=%+v", err, m.UserId, msg)
			}
			err = manager.PushMessage(m.UserId, &vo.PushMessage{
				Body:  msg,
				Index: int32(index),
			})
			if err != nil {
				log.Printf("[HandleGroupMessage] manager.PushMessage failed. err=%v userId=%d, msg=%+v", err, m.UserId, msg)
			}
		}(&member.GroupMember)
	}
	wg.Wait()
	return nil
}

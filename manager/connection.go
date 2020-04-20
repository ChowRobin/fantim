package manager

import (
	"github.com/ChowRobin/fantim/model/bo"
	"github.com/ChowRobin/fantim/model/vo"
	"github.com/gorilla/websocket"
)

// 单机先采用本地缓存链接关系，分布式采用redis
var (
	UserConnRouter map[int64]*bo.LConnectionGroup
)

func init() {
	UserConnRouter = make(map[int64]*bo.LConnectionGroup)
}

// 注册长连接
func RegisterUserLongConn(userId, connId int64, conn *websocket.Conn) error {
	connGroup, ok := UserConnRouter[userId]
	if !ok {
		connGroup = &bo.LConnectionGroup{}
		connGroup.Register(connId, conn)
		UserConnRouter[userId] = connGroup
	} else if connGroup != nil {
		connGroup.Register(connId, conn)
	}

	return nil
}

// 移除连接
func DeleteUserLongConn(userId, connId int64) {
	connGroup, ok := UserConnRouter[userId]
	if ok && connGroup != nil {
		delete(connGroup.ConnMap, connId)
	}
}

// 长链推送消息
func PushMessage(userId int64, msg *vo.PushMessage) error {
	connGroup, ok := UserConnRouter[userId]
	if !ok {
		return nil
	}
	return connGroup.BroadCast(msg)
}

package connection

import (
	"context"
	"encoding/json"
	"log"
	"net/http"

	"github.com/ChowRobin/fantim/constant/status"
	"github.com/ChowRobin/fantim/manager"
	"github.com/ChowRobin/fantim/model/vo"
	"github.com/ChowRobin/fantim/service"
	"github.com/ChowRobin/fantim/util"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

var upGrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func Handle(c *gin.Context) {
	ws, err := upGrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		return
	}
	defer ws.Close()
	connId := util.GenId()
	var userId int64
	for {
		var msg vo.PushMessage
		//websocket接受信息
		mt, msgBytes, err := ws.ReadMessage()
		// 关闭
		if mt == -1 {
			manager.DeleteUserLongConn(userId, connId)
			break
		}
		if err != nil {
			log.Printf("[connection.Handle] receive failed:%v messageType=%d", err, mt)
			continue
		}
		log.Printf("[connection.Handle] messageType=%v, msg=%v", mt, string(msgBytes))
		_ = json.Unmarshal(msgBytes, &msg)

		// -1为关闭长连接
		if msg.PushType == -1 {
			manager.DeleteUserLongConn(userId, connId)
			break
		}

		if msg.Body == nil {
			log.Printf("[connection.Handle] msg.Body is nil. msg=%v", msg)
			continue
		}

		// 初始化uid消息
		if msg.PushType == 0 {
			userId = msg.Body.Sender
			if userId == 0 {
				log.Printf("[websocket.Handle] userId is 0")
				continue
			}
		}

		if err = manager.RegisterUserLongConn(userId, connId, ws); err != nil {
			log.Printf("[connection.Handle] RegisterUserLongConn failed. err=%v", err)
		}

		if msg.PushType == 1 {
			msg.Body.Sender = userId
			ctx := context.Background()
			_, es := service.SendMessage(ctx, msg.Body)
			if es != status.Success {
				log.Printf("[connection.Handle] service.SendMessage failed. err=%v", es.Msg)
			}
		}
	}
}

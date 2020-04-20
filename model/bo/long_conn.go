package bo

import (
	"log"
	"time"

	"github.com/gorilla/websocket"
)

// 长链接
type LConnection struct {
	Id      int64
	Conn    *websocket.Conn
	Expired *time.Time
}

func (c *LConnection) IsExpired() bool {
	return c.Expired.Before(time.Now())
}

type LConnectionGroup struct {
	ConnMap map[int64]*LConnection
}

func (g *LConnectionGroup) Register(connId int64, conn *websocket.Conn) {
	timeout := time.Now().Add(time.Minute * 5)
	if g.ConnMap == nil {
		g.ConnMap = make(map[int64]*LConnection)
	}
	oldConn, ok := g.ConnMap[connId]
	if ok {
		oldConn.Expired = &timeout
	} else {
		g.ConnMap[connId] = &LConnection{
			Id:      connId,
			Conn:    conn,
			Expired: &timeout,
		}
	}
	_ = conn.SetReadDeadline(timeout)
}

func (g *LConnectionGroup) BroadCast(msg interface{}) error {
	if len(g.ConnMap) == 0 {
		log.Println("LongConnectGroup len is 0")
		return nil
	}
	for id, conn := range g.ConnMap {
		err := conn.Conn.WriteJSON(msg)
		if err != nil {
			if conn.IsExpired() {
				delete(g.ConnMap, id)
			} else {
				log.Printf("[LConnectionGroup.BroadCast] websocket message send failed. err=%v", err)
			}
		}
	}
	return nil
}

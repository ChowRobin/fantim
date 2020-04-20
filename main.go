package main

import (
	"github.com/ChowRobin/fantim/middleware"
	"github.com/ChowRobin/fantim/views/connection"
	"github.com/ChowRobin/fantim/views/message"
	"github.com/ChowRobin/fantim/views/relation"
	"github.com/ChowRobin/fantim/views/user"
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()

	store := cookie.NewStore([]byte("secret"))
	r.Use(sessions.Sessions("session_id", store))

	needLogin := middleware.ApiOption{
		Key: "login",
		Ok:  true,
	}
	r.POST("/user/login/", middleware.ApiDecorator(user.Login))
	r.POST("/user/sign/up/", middleware.ApiDecorator(user.SignUp))
	r.GET("/user/info/", middleware.ApiDecorator(user.Info, needLogin))

	r.POST("/message/send/", middleware.ApiDecorator(message.Send, needLogin))
	r.GET("/message/pull/", middleware.ApiDecorator(message.Pull, needLogin))
	r.GET("/websocket/create/", connection.Handle)

	r.POST("/relation/apply/create/", middleware.ApiDecorator(relation.CreateApply, needLogin))
	r.POST("/relation/apply/update/", middleware.ApiDecorator(relation.UpdateApply, needLogin))
	r.GET("/relation/apply/list/", middleware.ApiDecorator(relation.ListApply, needLogin))
	r.GET("/relation/friend/list/", middleware.ApiDecorator(relation.ListFriend, needLogin))

	_ = r.Run() // 监听并在 0.0.0.0:8080 上启动服务
}

package main

import (
	"github.com/ChowRobin/fantim/middleware"
	"github.com/ChowRobin/fantim/views/connection"
	"github.com/ChowRobin/fantim/views/group"
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
	/**
	store.Options(sessions.Options{
		Path:     "/",
		Domain:   "localhost:3000",
		MaxAge:   int(time.Hour * 24 * 3),
		Secure:   false,
		HttpOnly: false,
		SameSite: 0,
	})

	*/
	r.Use(sessions.Sessions("session_id", store))
	r.Use(middleware.Cors())

	needLogin := middleware.ApiOption{
		Key: "login",
		Ok:  true,
	}
	r.POST("/user/login/", middleware.ApiDecorator(user.Login))
	r.POST("/user/sign/up/", middleware.ApiDecorator(user.SignUp))
	r.GET("/user/info/", middleware.ApiDecorator(user.Info, needLogin))

	r.POST("/message/send/", middleware.ApiDecorator(message.Send, needLogin))
	r.POST("/message/pull/", middleware.ApiDecorator(message.Pull, needLogin))
	r.GET("/websocket/create/", connection.Handle)

	r.POST("/relation/apply/create/", middleware.ApiDecorator(relation.CreateApply, needLogin))
	r.POST("/relation/apply/update/", middleware.ApiDecorator(relation.UpdateApply, needLogin))
	r.GET("/relation/apply/list/", middleware.ApiDecorator(relation.ListApply, needLogin))
	r.GET("/relation/friend/list/", middleware.ApiDecorator(relation.ListFriend, needLogin))

	r.POST("/group/create/", middleware.ApiDecorator(group.Create, needLogin))
	r.GET("/group/list/", middleware.ApiDecorator(group.List, needLogin))
	r.GET("/group/member/list/", middleware.ApiDecorator(group.ListMembers, needLogin))
	r.GET("/group/search/", middleware.ApiDecorator(group.Search, needLogin))
	r.POST("/group/apply/update/", middleware.ApiDecorator(relation.UpdateGroupApply, needLogin))

	_ = r.Run() // 监听并在 0.0.0.0:8080 上启动服务
}

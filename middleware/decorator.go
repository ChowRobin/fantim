package middleware

import (
	"log"
	"strconv"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

type ApiFunc func(*gin.Context) interface{}

type ApiOption struct {
	Key   string
	Value string
	Ok    bool
}

func ApiDecorator(apiFunc ApiFunc, options ...ApiOption) func(*gin.Context) {
	return func(c *gin.Context) {
		var resp interface{}
		defer func() {
			c.JSON(200, resp)
		}()
		for _, op := range options {
			switch op.Key {
			case "login":
				if op.Ok {
					// 校验登陆态
					s := sessions.Default(c)
					uid := s.Get("user_id")

					if c.Query("user_id") != "" {
						uid, _ = strconv.ParseInt(c.Query("user_id"), 10, 64)
					}

					if uid == nil || uid == 0 {
						resp = map[string]interface{}{
							"status_code":    5,
							"status_message": "user not login",
						}
						return
					}
					log.Printf("[ApiDecorator] uid=%d", uid)
					c.Set("user_id", uid)
				}
			}
		}

		resp = apiFunc(c)
	}
}

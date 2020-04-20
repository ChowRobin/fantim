package user

import (
	"log"

	"github.com/ChowRobin/fantim/constant/status"
	"github.com/ChowRobin/fantim/model/po"
	"github.com/ChowRobin/fantim/model/vo"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

func Login(c *gin.Context) interface{} {
	req := &vo.LoginRequest{}
	resp := &vo.LoginResponse{}

	// 校验参数合法性
	err := c.BindJSON(req)
	if err != nil {
		log.Printf("[Login] BindJSON failed. err=%v", err)
		return status.FillResp(resp, status.ErrInvalidParam)
	}

	user := &po.User{
		UserId: req.UserId,
	}
	err = user.GetByUserId(c)
	if err != nil {
		log.Printf("[Login] po.user.GetByUserId failed. err=%v", err)
		return status.ErrInvalidParam
	}
	if user.Password != req.Password {
		return status.ErrInvalidPassword
	}

	s := sessions.Default(c)
	s.Set("user_id", user.UserId)
	err = s.Save()
	if err != nil {
		log.Printf("[Login] session.Save failed. err=%v", err)
		return status.FillResp(resp, status.ErrServiceInternal)
	}
	resp.UserInfo = &vo.User{
		UserId:   user.UserId,
		Nickname: user.Nickname,
	}

	return status.FillResp(resp, status.Success)
}

package user

import (
	"log"

	"github.com/ChowRobin/fantim/model/po"

	"github.com/ChowRobin/fantim/constant/status"
	"github.com/ChowRobin/fantim/model/vo"
	"github.com/gin-gonic/gin"
)

func SignUp(c *gin.Context) interface{} {
	req := &vo.SignUpRequest{}
	resp := &vo.SignUpResponse{}

	// 校验参数合法性
	err := c.BindJSON(req)
	if err != nil {
		log.Printf("[SignUp] BindJSON failed. err=%v", err)
		return status.FillResp(resp, status.ErrInvalidParam)
	}

	user := &po.User{
		UserId:   req.UserId,
		Password: req.Password,
		Nickname: req.Nickname,
	}
	err = user.GetByUserId(c)
	if err != nil {
		log.Printf("[SignUp] po.user.GetByUserId failed. err=%v", err)
		return status.FillResp(resp, status.ErrServiceInternal)
	}
	if user.Id != 0 {
		return status.FillResp(resp, status.ErrDuplicateUserId)
	}

	err = user.Create(c)
	if err != nil {
		log.Printf("[SignUp] po.user.Create failed. err=%v", err)
		return status.FillResp(resp, status.ErrServiceInternal)
	}

	return status.FillResp(resp, status.Success)
}

package user

import (
	"log"
	"strconv"

	"github.com/ChowRobin/fantim/constant/status"
	"github.com/ChowRobin/fantim/model/po"
	"github.com/ChowRobin/fantim/model/vo"
	"github.com/gin-gonic/gin"
)

func Info(c *gin.Context) interface{} {
	resp := &vo.UserInfoResponse{}

	userIdStr := c.Query("user_id")
	userId, _ := strconv.ParseInt(userIdStr, 10, 64)
	if userId == 0 {
		log.Printf("[User.Info] user_id is 0")
		return status.FillResp(resp, status.ErrInvalidParam)
	}

	userPo := po.User{
		UserId: userId,
	}
	err := userPo.GetByUserId(c)
	if err != nil {
		log.Printf("[User.Info] po user GetByUserId failed. err=%v", err)
		return status.FillResp(resp, status.ErrServiceInternal)
	}
	if userPo.Id == 0 {
		return status.FillResp(resp, status.ErrInvalidParam)
	}

	resp.UserInfo = &vo.User{
		UserId:   userPo.UserId,
		Nickname: userPo.Nickname,
		Avatar:   userPo.Avatar,
	}

	return status.FillResp(resp, status.Success)
}

package convert

import (
	"github.com/ChowRobin/fantim/model/po"
	"github.com/ChowRobin/fantim/model/vo"
)

func UserPoToVo(userPo *po.User) *vo.User {
	if userPo == nil {
		return nil
	}
	return &vo.User{
		UserId:   userPo.UserId,
		Nickname: userPo.Nickname,
		Avatar:   userPo.Avatar,
	}
}

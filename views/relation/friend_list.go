package relation

import (
	"log"

	"github.com/ChowRobin/fantim/constant"
	"github.com/ChowRobin/fantim/constant/status"
	"github.com/ChowRobin/fantim/model/po"
	"github.com/ChowRobin/fantim/model/vo"
	"github.com/gin-gonic/gin"
)

func ListFriend(c *gin.Context) interface{} {
	resp := &vo.FriendListResponse{}

	// 校验参数合法性
	userId := c.GetInt64("user_id")

	relations, err := po.ListUserRelationByFrom(c, userId, constant.RelationTypeFriend)
	if err != nil {
		log.Printf("[ListFriend] ListUserRelationByFrom failed. err=%v", err)
		return status.FillResp(resp, status.ErrServiceInternal)
	}
	if len(relations) == 0 {
		return status.FillResp(resp, status.Success)
	}

	userIds := make([]int64, 0, len(relations))
	for _, rel := range relations {
		userIds = append(userIds, rel.ToUserId)
	}

	userMap, err := po.MultiGetUserMapByUserId(c, userIds)
	if err != nil {
		log.Printf("[ListFriend] MultiGetUserMapByUserId failed. err=%v", err)
		return status.FillResp(resp, status.ErrServiceInternal)
	}

	for _, rel := range relations {
		if userPo, ok := userMap[rel.ToUserId]; ok {
			userVo := &vo.User{
				UserId:   userPo.UserId,
				Nickname: userPo.Nickname,
			}
			resp.Friends = append(resp.Friends, userVo)
		}
	}

	resp.TotalNum = int32(len(resp.Friends))

	return status.FillResp(resp, status.Success)
}

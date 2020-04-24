package group

import (
	"log"

	"github.com/ChowRobin/fantim/constant/status"
	"github.com/ChowRobin/fantim/model/po"
	"github.com/ChowRobin/fantim/model/vo"
	"github.com/ChowRobin/fantim/util"
	"github.com/gin-gonic/gin"
)

func Create(c *gin.Context) interface{} {
	req := &vo.GroupCreateRequest{}
	resp := &vo.GroupCreateResponse{}

	userId := c.GetInt64("user_id")

	// 校验参数合法性
	err := c.BindJSON(req)
	if err != nil {
		log.Printf("[CreateGroup] BindJSON failed. err=%v", err)
		return status.FillResp(resp, status.ErrInvalidParam)
	}

	groupId := util.GenId()
	groupPo := &po.Group{
		GroupId:     groupId,
		OwnerId:     userId,
		Name:        req.Name,
		Avatar:      req.Avatar,
		Description: req.Description,
	}
	err = groupPo.Create(c)
	if err != nil {
		log.Printf("[CreateGroup] groupPo.Create failed. err=%v", err)
		return status.FillResp(resp, status.ErrServiceInternal)
	}
	// 拉入初始化群成员
	groupMembers := make([]*po.GroupMember, 0, len(req.Members)+1)
	// 加入群主
	groupMembers = append(groupMembers, &po.GroupMember{
		GroupId: groupId,
		UserId:  userId,
		Role:    2, // 群主
	})
	// 加入群员
	for _, m := range req.Members {
		groupMembers = append(groupMembers, &po.GroupMember{
			GroupId: groupId,
			UserId:  m,
			Role:    1, // 群成员
		})
	}
	err = po.MultiAddMemberInGroup(c, groupMembers)
	if err != nil {
		log.Printf("[CreateGroup] po.MultiAddMemberInGroup failed. err=%v", err)
		return status.FillResp(resp, status.ErrServiceInternal)
	}

	resp.GroupId = groupId

	return status.FillResp(resp, status.Success)
}

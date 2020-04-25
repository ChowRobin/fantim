package group

import (
	"log"
	"strconv"

	"github.com/ChowRobin/fantim/constant/status"
	"github.com/ChowRobin/fantim/model/po"
	"github.com/ChowRobin/fantim/model/vo"
	"github.com/gin-gonic/gin"
)

func List(c *gin.Context) interface{} {
	resp := &vo.GroupListResponse{}

	userId := c.GetInt64("user_id")

	groupPoList, err := po.ListGroupByUserId(c, userId)
	if err != nil {
		log.Printf("[ListGroup] po.ListGroupByUserId failed. err=%v", err)
		return status.FillResp(resp, status.ErrServiceInternal)
	}
	result := make([]*vo.GroupInfo, 0, len(groupPoList))
	for _, p := range groupPoList {
		result = append(result, &vo.GroupInfo{
			GroupId:     p.GroupMember.GroupId,
			GroupIdStr:  strconv.FormatInt(p.GroupMember.GroupId, 10),
			OwnerUid:    p.OwnerId,
			Name:        p.Name,
			Avatar:      p.Avatar,
			Description: p.Description,
			UserRole:    int32(p.Role),
		})
	}
	resp.Groups = result

	return status.FillResp(resp, status.Success)
}

package group

import (
	"log"
	"strconv"

	"github.com/ChowRobin/fantim/constant/status"
	"github.com/ChowRobin/fantim/model/convert"
	"github.com/ChowRobin/fantim/model/po"
	"github.com/ChowRobin/fantim/model/vo"
	"github.com/gin-gonic/gin"
)

func ListMembers(c *gin.Context) interface{} {
	resp := &vo.GroupMemberListResponse{}

	// todo 校验用户是否是群成员，否则不可查询

	groupId, _ := strconv.ParseInt(c.Query("group_id"), 10, 64)
	if groupId == 0 {
		return status.FillResp(resp, status.ErrInvalidParam)
	}

	memberPoList, err := po.ListMembersByGroupId(c, groupId)
	if err != nil {
		log.Printf("[ListGroupMembers] po.ListMembersByGroupId failed. err=%v", err)
		return status.FillResp(resp, status.ErrServiceInternal)
	}
	result := make([]*vo.GroupMember, 0, len(memberPoList))
	for _, m := range memberPoList {
		result = append(result, &vo.GroupMember{
			UserInfo: convert.UserPoToVo(&m.User),
			UserRole: int32(m.Role),
		})
	}
	resp.Members = result

	return status.FillResp(resp, status.Success)
}

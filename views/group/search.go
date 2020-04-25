package group

import (
	"log"
	"strconv"

	"github.com/ChowRobin/fantim/constant/status"
	"github.com/ChowRobin/fantim/model/po"
	"github.com/ChowRobin/fantim/model/vo"
	"github.com/gin-gonic/gin"
)

func Search(c *gin.Context) interface{} {
	resp := &vo.GroupSearchResponse{}

	userId := c.GetInt64("user_id")
	searchName := c.Query("name")
	page, _ := strconv.Atoi(c.Query("page"))
	pageSize, _ := strconv.Atoi(c.Query("page_size"))
	// 参数校验
	if searchName == "" || page == 0 || pageSize == 0 || pageSize > 100 {
		return status.FillResp(resp, status.ErrInvalidParam)
	}

	// 查询总数
	totalNum, err := po.CountGroupByName(c, searchName)
	if err != nil {
		log.Printf("[GroupSearch] po.CountGroupByName failed. err=%v", err)
		return status.FillResp(resp, status.ErrServiceInternal)
	}
	if totalNum == 0 {
		return status.FillResp(resp, status.Success)
	}

	if (page-1)*pageSize >= int(totalNum) {
		return status.FillResp(resp, status.ErrInvalidPageParam)
	}

	// 搜索
	groups, err := po.SearchGroupByName(c, searchName, int32(page), int32(pageSize))
	if err != nil {
		log.Printf("[GroupSearch] po.SearchGroupByName failed. err=%v", err)
		return status.FillResp(resp, status.ErrServiceInternal)
	}
	if len(groups) == 0 {
		return status.FillResp(resp, status.Success)
	}

	// 查询用户的群列表
	ownGroupList, err := po.ListGroupByUserId(c, userId)
	if err != nil {
		log.Printf("[GroupSearch] po.ListGroupByUserId failed. err=%v", err)
		return status.FillResp(resp, status.ErrServiceInternal)
	}
	ownGroupMap := make(map[int64]*po.GroupMemberWithGroup)
	for _, ownGroup := range ownGroupList {
		ownGroupMap[ownGroup.Group.GroupId] = ownGroup
	}

	// 拼装数据
	result := make([]*vo.GroupInfo, 0, len(groups))
	for _, g := range groups {
		groupVo := &vo.GroupInfo{
			GroupId:     g.GroupId,
			GroupIdStr:  strconv.FormatInt(g.GroupId, 10),
			OwnerUid:    g.OwnerId,
			Name:        g.Name,
			Avatar:      g.Avatar,
			Description: g.Description,
			UserRole:    0,
		}
		ownInfo := ownGroupMap[g.GroupId]
		if ownInfo != nil {
			groupVo.UserRole = int32(ownInfo.Role)
		}
		result = append(result, groupVo)
	}

	resp.Groups = result
	resp.TotalNum = 1

	return status.FillResp(resp, status.Success)
}

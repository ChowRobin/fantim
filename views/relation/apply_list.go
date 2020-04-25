package relation

import (
	"fmt"
	"log"
	"strconv"

	"github.com/ChowRobin/fantim/constant"

	"github.com/ChowRobin/fantim/constant/status"
	"github.com/ChowRobin/fantim/model/po"
	"github.com/ChowRobin/fantim/model/vo"
	"github.com/gin-gonic/gin"
)

func ListApply(c *gin.Context) interface{} {
	resp := &vo.RelationApplyListResponse{}

	// 校验参数合法性
	userId := c.GetInt64("user_id")
	fromUserId, _ := strconv.ParseInt(c.Query("from_user_id"), 10, 64)
	applyType, _ := strconv.Atoi(c.Query("apply_type"))
	page, _ := strconv.Atoi(c.Query("page"))
	pageSize, _ := strconv.Atoi(c.Query("page_size"))
	queryStatusStr := c.QueryArray("status")
	var queryStatus []int32
	for _, statusStr := range queryStatusStr {
		if s, err := strconv.Atoi(statusStr); err == nil {
			queryStatus = append(queryStatus, int32(s))
		}
	}
	// 参数校验
	if page == 0 || pageSize == 0 || pageSize > 100 {
		return status.FillResp(resp, status.ErrInvalidParam)
	}

	var toIds []int64
	if fromUserId != 0 {
		// 不允许查其他人的申请列表
		if fromUserId != userId {
			return status.FillResp(resp, status.ErrInvalidParam)
		}
	} else {
		if applyType == 1 { // 好友申请
			toIds = append(toIds, userId)
		} else if applyType == 2 { // 加群申请
			groups, err := po.ListGroupByCondition(c, userId, []int32{2})
			if err != nil {
				log.Printf("[ListApply] ListGroupByCondition failed. err=%v", err)
				return status.FillResp(resp, status.ErrServiceInternal)
			}
			for _, g := range groups {
				toIds = append(toIds, g.GroupId)
			}
		}
	}

	// 查询总数
	totalNum, err := po.CountUserRelationApplyPageByCondition(c, fromUserId, toIds, queryStatus, int32(applyType))
	if err != nil {
		log.Printf("[ListApply] po.CountUserRelationApplyPageByCondition failed. err=%v", err)
		return status.FillResp(resp, status.ErrServiceInternal)
	}
	if totalNum == 0 {
		return status.FillResp(resp, status.Success)
	}

	if (page-1)*pageSize >= int(totalNum) {
		return status.FillResp(resp, status.ErrInvalidPageParam)
	}

	// 分页查询记录
	applyPoList, err := po.ListUserRelationApplyPageByCondition(c, fromUserId, toIds, queryStatus, int32(applyType), int32(page), int32(pageSize))
	if err != nil {
		log.Printf("[ListApply] po.ListUserRelationApplyPageByCondition failed. err=%v", err)
		return status.FillResp(resp, status.ErrServiceInternal)
	}

	groupMap := make(map[int64]*vo.GroupInfo)
	if len(applyPoList) > 0 && applyType == constant.RelationApplyTypeGroup {
		reqGroupIds := make([]int64, 0, len(applyPoList))
		for _, apply := range applyPoList {
			reqGroupIds = append(reqGroupIds, apply.ToUserId)
		}
		groupsPo, err := po.MultiGetGroup(c, reqGroupIds)
		if err != nil {
			log.Printf("[ListApply] po.MultiGetGroup failed. err=%v", err)
			return status.FillResp(resp, status.ErrServiceInternal)
		}
		for _, g := range groupsPo {
			groupMap[g.GroupId] = &vo.GroupInfo{
				GroupId:     g.GroupId,
				OwnerUid:    g.OwnerId,
				Name:        g.Name,
				Avatar:      g.Avatar,
				Description: g.Description,
				GroupIdStr:  fmt.Sprintf("%d", g.GroupId),
			}
		}
	}

	for _, applyPo := range applyPoList {
		applyVo := &vo.RelationApply{
			FromUserId: applyPo.FromUserId,
			ToUserId:   applyPo.ToUserId,
			ApplyType:  int32(applyPo.ApplyType),
			Status:     int32(applyPo.Status),
			Content:    applyPo.Content,
			ApplyId:    applyPo.Id,
		}
		if applyType == constant.RelationApplyTypeGroup {
			applyVo.GroupInfo = groupMap[applyPo.ToUserId]
		}
		resp.ApplyList = append(resp.ApplyList, applyVo)
	}
	resp.TotalNum = totalNum

	return status.FillResp(resp, status.Success)
}

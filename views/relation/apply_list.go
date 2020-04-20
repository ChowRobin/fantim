package relation

import (
	"log"
	"strconv"

	"github.com/ChowRobin/fantim/constant/status"
	"github.com/ChowRobin/fantim/model/po"
	"github.com/ChowRobin/fantim/model/vo"
	"github.com/gin-gonic/gin"
)

func ListApply(c *gin.Context) interface{} {
	resp := &vo.RelationApplyListResponse{}

	// 校验参数合法性
	userId := c.GetInt64("user_id")
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

	// 查询总数
	totalNum, err := po.CountUserRelationApplyPageByCondition(c, nil, &userId, queryStatus, int32(applyType))
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
	applyPoList, err := po.ListUserRelationApplyPageByCondition(c, nil, &userId, queryStatus, int32(applyType), int32(page), int32(pageSize))
	if err != nil {
		log.Printf("[ListApply] po.ListUserRelationApplyPageByCondition failed. err=%v", err)
		return status.FillResp(resp, status.ErrServiceInternal)
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
		resp.ApplyList = append(resp.ApplyList, applyVo)
	}
	resp.TotalNum = totalNum

	return status.FillResp(resp, status.Success)
}

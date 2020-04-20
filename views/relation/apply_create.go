package relation

import (
	"log"

	"github.com/ChowRobin/fantim/constant"

	"github.com/ChowRobin/fantim/model/po"

	"github.com/ChowRobin/fantim/constant/status"
	"github.com/ChowRobin/fantim/model/vo"
	"github.com/gin-gonic/gin"
)

func CreateApply(c *gin.Context) interface{} {
	req := &vo.RelationApplyCreateRequest{}
	resp := &vo.RelationApplyCreateResponse{}

	userId := c.GetInt64("user_id")

	// 校验参数合法性
	err := c.BindJSON(req)
	if err != nil {
		log.Printf("[CreateApply] BindJSON failed. err=%v", err)
		return status.FillResp(resp, status.ErrInvalidParam)
	}

	if userId == req.ToUserId || req.ApplyType != constant.RelationApplyTypeFriend {
		return status.FillResp(resp, status.ErrInvalidParam)
	}

	applyPo := &po.UserRelationApply{
		FromUserId: userId,
		ToUserId:   req.ToUserId,
		ApplyType:  int8(req.ApplyType),
		Status:     int8(0),
	}
	err = applyPo.GetByCondition(c)
	// 已有未处理申请
	if applyPo.Id != 0 {
		resp.ApplyId = applyPo.Id
		return status.FillResp(resp, status.Success)
	}

	err = applyPo.Create(c)
	if err != nil {
		log.Printf("[CreateApply] applyPo.Create failed. err=%v", err)
		return status.FillResp(resp, status.ErrServiceInternal)
	}

	resp.ApplyId = applyPo.Id
	return status.FillResp(resp, status.Success)
}

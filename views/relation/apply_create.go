package relation

import (
	"log"
	"strconv"

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

	var toId int64
	if req.ApplyType == constant.RelationApplyTypeFriend {
		if userId == req.ToUserId {
			return status.FillResp(resp, status.ErrInvalidParam)
		}
		toId = req.ToUserId
	} else if req.ApplyType == constant.RelationApplyTypeGroup {
		toId, _ = strconv.ParseInt(req.GroupIdStr, 10, 64)
	}

	applyPo := &po.UserRelationApply{
		FromUserId: userId,
		ToUserId:   toId,
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

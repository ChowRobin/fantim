package relation

import (
	"log"

	"github.com/ChowRobin/fantim/client"
	"github.com/ChowRobin/fantim/constant"
	"github.com/ChowRobin/fantim/constant/status"
	"github.com/ChowRobin/fantim/model/po"
	"github.com/ChowRobin/fantim/model/vo"
	"github.com/gin-gonic/gin"
)

// 好友申请
func UpdateApply(c *gin.Context) interface{} {
	req := &vo.RelationApplyUpdateRequest{}
	resp := &vo.RelationApplyUpdateResponse{}

	userId := c.GetInt64("user_id")
	// 校验参数合法性
	err := c.BindJSON(req)
	if err != nil {
		log.Printf("[UpdateApply] BindJSON failed. err=%v", err)
		return status.FillResp(resp, status.ErrInvalidParam)
	}

	// 判断是from操作，还是to操作
	var isFromOp, isToOp bool
	switch req.Status {
	case constant.RelationApplyStatusPass, constant.RelationApplyStatusReject:
		isToOp = true
	case constant.RelationApplyStatusCancel:
		isFromOp = true
	default:
		log.Printf("[UpdateApply] status is not legal. status=%d", req.Status)
		return status.FillResp(resp, status.ErrInvalidParam)
	}

	applyPo := &po.UserRelationApply{
		Id: req.ApplyId,
	}
	err = applyPo.GetById(c)
	if err != nil {
		log.Printf("[UpdateApply] applyPo.GetById failed. err=%v", err)
		return status.FillResp(resp, status.ErrServiceInternal)
	}
	if applyPo.FromUserId == 0 || applyPo.ToUserId == 0 {
		log.Printf("[UpdateApply] apply id is not exsist. id=%d", applyPo.Id)
		return status.FillResp(resp, status.ErrInvalidParam)
	}

	var ok bool
	if isFromOp {
		ok = applyPo.FromUserId == userId
	} else if isToOp {
		ok = applyPo.ToUserId == userId
	}
	if !ok {
		log.Printf("[UpdateApply] userId from to op check failed. userId=%d req=%+v", userId, req)
		return status.FillResp(resp, status.ErrInvalidParam)
	}

	// 开启事务
	err = client.StartDBTransaction(c)
	if err != nil {
		log.Printf("[UpdateApply] StartDBTransaction failed. err=%v", err)
	} else {
		// 提交事务
		defer func() {
			err = client.CommitDBTransaction(c)
			if err != nil {
				log.Printf("[UpdateApply] CommitDBTransaction failed. err=%v", err)
			}
		}()
	}

	// 修改状态
	applyPo.Status = int8(req.Status)
	err = applyPo.Update(c)
	if err != nil {
		log.Printf("[UpdateApply] applyPo.Update failed. err=%v", err)
		return status.FillResp(resp, status.ErrServiceInternal)
	}
	fromRelationPo := &po.UserRelation{}
	toRelationPo := &po.UserRelation{}
	// 修改关联数据
	switch req.Status {
	case constant.RelationApplyStatusPass:
		fromRelationPo.FromUserId = applyPo.FromUserId
		fromRelationPo.ToUserId = applyPo.ToUserId
		fromRelationPo.Status = constant.RelationTypeFriend

		toRelationPo.FromUserId = applyPo.ToUserId
		toRelationPo.ToUserId = applyPo.FromUserId
		toRelationPo.Status = constant.RelationTypeFriend

		err = fromRelationPo.GetByCondition(c)
		if err != nil {
			log.Printf("[UpdateApply] fromRelationPo.GetByCondition failed. err=%v", err)
			return status.FillResp(resp, status.ErrServiceInternal)
		}
		if fromRelationPo.Id == 0 {
			err = fromRelationPo.Create(c)
			if err != nil {
				log.Printf("[UpdateApply] fromRelationPo.Create failed. err=%v", err)
				return status.FillResp(resp, status.ErrServiceInternal)
			}
		}

		err = toRelationPo.GetByCondition(c)
		if err != nil {
			log.Printf("[UpdateApply] fromRelationPo.GetByCondition failed. err=%v", err)
			return status.FillResp(resp, status.ErrServiceInternal)
		}
		if toRelationPo.Id == 0 {
			err = toRelationPo.Create(c)
			if err != nil {
				log.Printf("[UpdateApply] toRelationPo.Create failed. err=%v", err)
				return status.FillResp(resp, status.ErrServiceInternal)
			}
		}
	}

	return status.FillResp(resp, status.Success)
}

// 群聊申请
func UpdateGroupApply(c *gin.Context) interface{} {
	req := &vo.RelationApplyUpdateRequest{}
	resp := &vo.RelationApplyUpdateResponse{}

	userId := c.GetInt64("user_id")
	// 校验参数合法性
	err := c.BindJSON(req)
	if err != nil {
		log.Printf("[UpdateGroupApply] BindJSON failed. err=%v", err)
		return status.FillResp(resp, status.ErrInvalidParam)
	}

	// 判断是用户操作，还是群管操作
	var isFromOp, isToOp bool
	var userGroupRole int8
	switch req.Status {
	case constant.RelationApplyStatusPass, constant.RelationApplyStatusReject:
		isToOp = true

	case constant.RelationApplyStatusCancel:
		isFromOp = true
	default:
		log.Printf("[UpdateGroupApply] status is not legal. status=%d", req.Status)
		return status.FillResp(resp, status.ErrInvalidParam)
	}

	applyPo := &po.UserRelationApply{
		Id: req.ApplyId,
	}
	err = applyPo.GetById(c)
	if err != nil {
		log.Printf("[UpdateGroupApply] applyPo.GetById failed. err=%v", err)
		return status.FillResp(resp, status.ErrServiceInternal)
	}
	if applyPo.FromUserId == 0 || applyPo.ToUserId == 0 {
		log.Printf("[UpdateGroupApply] apply id is not exsist. id=%d", applyPo.Id)
		return status.FillResp(resp, status.ErrInvalidParam)
	}

	// 查询操作用户的群身份
	gm := &po.GroupMember{
		GroupId: applyPo.ToUserId,
		UserId:  userId,
	}
	userGroupRole, err = gm.GetMemberRole(c)
	if err != nil {
		log.Printf("[UpdateGroupApply] GetMemberRole failed. err=%v", err)
		return status.FillResp(resp, status.ErrServiceInternal)
	}

	var ok bool
	if isFromOp {
		ok = applyPo.FromUserId == userId
	} else if isToOp {
		if userGroupRole == 2 || userGroupRole == 3 {
			ok = true
		}
	}
	if !ok {
		log.Printf("[UpdateGroupApply] userId from to op check failed. userId=%d req=%+v", userId, req)
		return status.FillResp(resp, status.ErrInvalidParam)
	}

	// todo 优化整体连接关闭方法，使用事务
	/**
	err = client.StartDBTransaction(c)
	if err != nil {
		log.Printf("[UpdateGroupApply] StartDBTransaction failed. err=%v", err)
	} else {
		// 提交事务
		defer func() {
			err = client.CommitDBTransaction(c)
			if err != nil {
				log.Printf("[UpdateGroupApply] CommitDBTransaction failed. err=%v", err)
			}
		}()
	}

	*/

	// 修改状态
	applyPo.Status = int8(req.Status)
	err = applyPo.Update(c)
	if err != nil {
		log.Printf("[UpdateGroupApply] applyPo.Update failed. err=%v", err)
		return status.FillResp(resp, status.ErrServiceInternal)
	}
	groupMemberPo := &po.GroupMember{
		GroupId: applyPo.ToUserId,
		UserId:  applyPo.FromUserId,
	}
	// 修改关联数据
	switch req.Status {
	case constant.RelationApplyStatusPass:
		role, err := groupMemberPo.GetMemberRole(c)
		if err != nil {
			log.Printf("[UpdateGroupApply] groupMemberPo.GetMemberRole failed. err=%v", err)
			return status.FillResp(resp, status.ErrServiceInternal)
		}
		if role == 0 {
			groupMemberPo.Role = 1 // 群成员
			err = groupMemberPo.Create(c)
			if err != nil {
				log.Printf("[UpdateGroupApply] groupMemberPo.Create failed. err=%v", err)
				return status.FillResp(resp, status.ErrServiceInternal)
			}
		}

	}

	return status.FillResp(resp, status.Success)
}

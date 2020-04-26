package po

import (
	"context"
	"time"

	"github.com/ChowRobin/fantim/client"
	"github.com/jinzhu/gorm"
)

type UserRelationApply struct {
	Id         int64  `gorm:"primary_key"`
	FromUserId int64  `gorm:"column:from_uid"`
	ToUserId   int64  `gorm:"column:to_uid"`     // 存储uid 或 群id
	ApplyType  int8   `gorm:"column:apply_type"` // 1->好友申请， 2->群申请
	Status     int8   `gorm:"column:status"`
	Content    string `gorm:"column:content"`

	CreateTime *time.Time `gorm:"column:create_time"`
	UpdateTime *time.Time `gorm:"column:update_time"`
}

func (*UserRelationApply) TableName() string {
	return "user_relation_apply"
}

func (ua *UserRelationApply) Create(ctx context.Context) error {
	conn, err := client.DBConn(ctx)
	if err != nil {
		return err
	}
	defer conn.Close()
	return conn.Create(ua).Error
}

func (ua *UserRelationApply) Update(ctx context.Context) error {
	conn, err := client.DBConn(ctx)
	if err != nil {
		return err
	}
	defer conn.Close()
	err = conn.Model(ua).Where("id=?", ua.Id).UpdateColumn("status", ua.Status).Error
	if err != nil {
		conn.Rollback()
	}
	return err
}

func (ua *UserRelationApply) GetById(ctx context.Context) error {
	conn, err := client.DBConn(ctx)
	if err != nil {
		return err
	}
	defer conn.Close()
	return conn.Model(ua).Where("id=?", ua.Id).First(ua).Error
}

func (ua *UserRelationApply) GetByCondition(ctx context.Context) error {
	conn, err := client.DBConn(ctx)
	if err != nil {
		return err
	}
	defer conn.Close()
	return conn.Model(ua).Where("from_uid=? and to_uid=? and apply_type=? and status=?",
		ua.FromUserId, ua.ToUserId, ua.ApplyType, ua.Status).First(ua).Error
}

func ListUserRelationApplyPageByCondition(ctx context.Context, fromUid int64, toIds []int64, queryStatus []int32, applyType, page, pageSize int32) ([]*UserRelationApply, error) {
	conn, err := client.DBConn(ctx)
	if err != nil {
		return nil, err
	}
	defer conn.Close()
	if len(toIds) == 0 {
		conn = conn.Where("from_uid=?", fromUid)
	} else {
		conn = conn.Where("to_uid in (?)", toIds)
	}
	conn = conn.Where("apply_type=?", applyType)
	if len(queryStatus) > 0 {
		conn = conn.Where("status in (?)", queryStatus)
	}
	var result []*UserRelationApply
	err = conn.Model(&UserRelationApply{}).Order("create_time desc").Offset((page - 1) * pageSize).Limit(pageSize).Find(&result).Error
	if err == gorm.ErrRecordNotFound {
		return nil, nil
	}
	return result, err
}

func CountUserRelationApplyPageByCondition(ctx context.Context, fromUid int64, toIds []int64, queryStatus []int32, applyType int32) (int32, error) {
	conn, err := client.DBConn(ctx)
	if err != nil {
		return 0, err
	}
	defer conn.Close()
	if len(toIds) == 0 {
		conn = conn.Where("from_uid=?", fromUid)
	} else {
		conn = conn.Where("to_uid in (?)", toIds)
	}
	conn = conn.Where("apply_type=?", applyType)
	if len(queryStatus) > 0 {
		conn = conn.Where("status in (?)", queryStatus)
	}
	var result int32
	err = conn.Model(&UserRelationApply{}).Count(&result).Error
	if err == gorm.ErrRecordNotFound {
		return 0, nil
	}
	return result, err
}

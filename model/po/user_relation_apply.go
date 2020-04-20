package po

import (
	"context"
	"errors"
	"time"

	"github.com/ChowRobin/fantim/client"
	"github.com/jinzhu/gorm"
)

type UserRelationApply struct {
	Id         int64  `gorm:"primary_key"`
	FromUserId int64  `gorm:"column:from_uid"`
	ToUserId   int64  `gorm:"column:to_uid"`
	ApplyType  int8   `gorm:"column:apply_type"`
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
	return conn.Create(ua).Error
}

func (ua *UserRelationApply) Update(ctx context.Context) error {
	conn, err := client.DBConn(ctx)
	if err != nil {
		return err
	}
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
	return conn.Model(ua).Where("id=?", ua.Id).First(ua).Error
}

func (ua *UserRelationApply) GetByCondition(ctx context.Context) error {
	conn, err := client.DBConn(ctx)
	if err != nil {
		return err
	}
	return conn.Model(ua).Where("from_uid=? and to_uid=? and apply_type=? and status=?",
		ua.FromUserId, ua.ToUserId, ua.ApplyType, ua.Status).First(ua).Error
}

func ListUserRelationApplyPageByCondition(ctx context.Context, fromUid, toUid *int64, queryStatus []int32, applyType, page, pageSize int32) ([]*UserRelationApply, error) {
	conn, err := client.DBConn(ctx)
	if err != nil {
		return nil, err
	}
	if fromUid != nil {
		conn = conn.Where("from_uid=?", fromUid)
	} else if toUid != nil {
		conn = conn.Where("to_uid=?", toUid)
	} else {
		return nil, errors.New("from to uid both nil")
	}
	conn = conn.Where("apply_type=?", applyType)
	if len(queryStatus) > 0 {
		conn = conn.Where("status in (?)", queryStatus)
	}
	var result []*UserRelationApply
	err = conn.Model(&UserRelationApply{}).Offset((page - 1) * pageSize).Limit(pageSize).Find(&result).Error
	if err == gorm.ErrRecordNotFound {
		return nil, nil
	}
	return result, err
}

func CountUserRelationApplyPageByCondition(ctx context.Context, fromUid, toUid *int64, queryStatus []int32, applyType int32) (int32, error) {
	conn, err := client.DBConn(ctx)
	if err != nil {
		return 0, err
	}
	if fromUid != nil {
		conn = conn.Where("from_uid=?", fromUid)
	} else if toUid != nil {
		conn = conn.Where("to_uid=?", toUid)
	} else {
		return 0, errors.New("from to uid both nil")
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

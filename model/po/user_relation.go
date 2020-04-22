package po

import (
	"context"
	"time"

	"github.com/ChowRobin/fantim/client"
	"github.com/jinzhu/gorm"
)

type UserRelation struct {
	Id         int64 `gorm:"primary_key"`
	FromUserId int64 `gorm:"column:from_uid"`
	ToUserId   int64 `gorm:"column:to_uid"`
	Status     int8  `gorm:"column:status"`

	CreateTime *time.Time `gorm:"column:create_time"`
	UpdateTime *time.Time `gorm:"column:update_time"`
}

func (*UserRelation) TableName() string {
	return "user_relation"
}

func (ur *UserRelation) Create(ctx context.Context) error {
	conn, err := client.DBConn(ctx)
	if err != nil {
		return err
	}
	err = conn.Create(ur).Error
	if err != nil {
		conn.Rollback()
	}
	return err
}

func (ur *UserRelation) Update(ctx context.Context) error {
	conn, err := client.DBConn(ctx)
	if err != nil {
		return err
	}
	err = conn.Model(ur).Where("from_uid=? AND to_uid=?", ur.FromUserId, ur.ToUserId).
		UpdateColumn("status", ur.Status).Error
	if err != nil {
		conn.Rollback()
	}
	return err
}

func (ur *UserRelation) GetByCondition(ctx context.Context) error {
	conn, err := client.DBConn(ctx)
	if err != nil {
		return err
	}

	defer conn.Close()
	err = conn.Model(ur).Where("from_uid=? AND to_uid=?", ur.FromUserId, ur.ToUserId).First(ur).Error
	if err == gorm.ErrRecordNotFound {
		return nil
	}
	return err
}

func ListUserRelationByFrom(ctx context.Context, fromUserId int64, status int32) ([]*UserRelation, error) {
	conn, err := client.DBConn(ctx)
	if err != nil {
		return nil, err
	}
	defer conn.Close()
	var rels []*UserRelation
	err = conn.Model(&UserRelation{}).Where("from_uid=? and status=?", fromUserId, status).Find(&rels).Error
	if err == gorm.ErrRecordNotFound {
		return nil, nil
	}
	return rels, err
}

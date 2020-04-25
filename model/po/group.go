package po

import (
	"context"
	"fmt"
	"time"

	"github.com/ChowRobin/fantim/client"
	"github.com/jinzhu/gorm"
)

type Group struct {
	Id          int64  `gorm:"primary_key"`
	GroupId     int64  `gorm:"column:group_id"`
	OwnerId     int64  `gorm:"column:owner_uid"`
	Name        string `gorm:"column:name"`
	Avatar      string `gorm:"column:avatar"`
	Description string `gorm:"column:description"`

	CreateTime *time.Time `gorm:"column:create_time"`
	UpdateTime *time.Time `gorm:"column:update_time"`
}

func (*Group) TableName() string {
	return "im_group_base"
}

func (g *Group) Create(ctx context.Context) error {
	conn, err := client.DBConn(ctx)
	if err != nil {
		return err
	}
	defer conn.Close()
	return conn.Create(g).Error
}

func (g *Group) GetByGroupId(ctx context.Context) error {
	conn, err := client.DBConn(ctx)
	if err != nil {
		return err
	}
	defer conn.Close()
	err = conn.Model(g).Where("group_id=?", g.GroupId).First(g).Error
	if err == gorm.ErrRecordNotFound {
		return nil
	}
	return err
}

func MultiGetGroup(ctx context.Context, groupIds []int64) (groups []*Group, err error) {
	conn, err := client.DBConn(ctx)
	if err != nil {
		return
	}
	defer conn.Close()

	conn = conn.Debug()
	err = conn.Model(&Group{}).Where("group_id in (?)", groupIds).Find(&groups).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			err = nil
			return
		}
	}
	return
}

func SearchGroupByName(ctx context.Context, searchName string, page, pageSize int32) (groups []*Group, err error) {
	conn, err := client.DBConn(ctx)
	if err != nil {
		return
	}
	defer conn.Close()

	conn = conn.Debug()
	statement := fmt.Sprintf("MATCH (name) AGAINST('%s*' in boolean mode)", searchName)
	err = conn.Model(&Group{}).Where(statement).
		Offset((page - 1) * pageSize).Limit(pageSize).Find(&groups).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			err = nil
			return
		}
	}

	return
}

func CountGroupByName(ctx context.Context, searchName string) (totalNum int64, err error) {
	conn, err := client.DBConn(ctx)
	if err != nil {
		return
	}
	defer conn.Close()

	conn = conn.Debug()
	statement := fmt.Sprintf("MATCH (name) AGAINST('%s*' in boolean mode)", searchName)
	err = conn.Model(&Group{}).Where(statement).Count(&totalNum).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			err = nil
			return
		}
	}

	return
}

package po

import (
	"context"
	"time"

	"github.com/jinzhu/gorm"

	"github.com/ChowRobin/fantim/client"
)

type User struct {
	Id       int64  `gorm:"primary_key"`
	UserId   int64  `gorm:"column:uid"`
	Password string `gorm:"column:password"`
	Nickname string `gorm:"column:nickname"`
	Avatar   string `gorm:"column:avatar"`

	CreateTime *time.Time `gorm:"column:create_time"`
	UpdateTime *time.Time `gorm:"column:update_time"`
}

func (*User) TableName() string {
	return "user_base"
}

func (u *User) Create(ctx context.Context) error {
	conn, err := client.DBConn(ctx)
	if err != nil {
		return err
	}
	defer conn.Close()
	return conn.Create(u).Error
}

func (u *User) GetByUserId(ctx context.Context) error {
	conn, err := client.DBConn(ctx)
	if err != nil {
		return err
	}
	defer conn.Close()
	err = conn.Model(u).Where("uid=?", u.UserId).First(&u).Error
	if err == gorm.ErrRecordNotFound {
		return nil
	}
	return err
}

func MultiGetUserByUserId(ctx context.Context, userIds []int64) ([]*User, error) {
	conn, err := client.DBConn(ctx)
	if err != nil {
		return nil, err
	}
	defer conn.Close()
	var result []*User
	err = conn.Model(&User{}).Where("uid in (?)", userIds).Find(&result).Error
	if err == gorm.ErrRecordNotFound {
		return nil, nil
	}
	return result, err
}

func MultiGetUserMapByUserId(ctx context.Context, userIds []int64) (map[int64]*User, error) {
	users, err := MultiGetUserByUserId(ctx, userIds)
	if err != nil {
		return nil, err
	}
	result := make(map[int64]*User, len(users))
	for _, u := range users {
		result[u.UserId] = u
	}
	return result, nil
}

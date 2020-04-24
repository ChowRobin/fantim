package po

import (
	"context"
	"fmt"
	"time"

	"github.com/ChowRobin/fantim/client"
	"github.com/jinzhu/gorm"
)

type GroupMember struct {
	Id      int64 `gorm:"primary_key"`
	GroupId int64 `gorm:"column:group_id"`
	UserId  int64 `gorm:"column:uid"`
	Role    int8  `gorm:"column:role"`

	CreateTime *time.Time `gorm:"column:create_time"`
	UpdateTime *time.Time `gorm:"column:update_time"`
}

type GroupMemberWithUser struct {
	GroupMember
	User
}

type GroupMemberWithGroup struct {
	GroupMember
	Group
}

func (*GroupMember) TableName() string {
	return "im_group_member"
}

func (g *GroupMember) Create(ctx context.Context) error {
	conn, err := client.DBConn(ctx)
	if err != nil {
		return err
	}
	defer conn.Close()
	return conn.Create(g).Error
}

func (g *GroupMember) GetMemberRole(ctx context.Context) (role int8, err error) {
	conn, err := client.DBConn(ctx)
	if err != nil {
		return
	}
	defer conn.Close()
	err = conn.Model(g).Where("group_id = ? and uid = ?", g.GroupId, g.UserId).Find(g).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return 0, nil
		}
		return
	}
	if g != nil {
		role = g.Role
	}
	return
}

func MultiAddMemberInGroup(ctx context.Context, members []*GroupMember) error {
	if len(members) == 0 {
		return nil
	}
	conn, err := client.DBConn(ctx)
	if err != nil {
		return err
	}
	defer conn.Close()
	conn = conn.Debug()
	insertSql := "INSERT INTO `im_group_member` (`group_id`, `uid`, `role`) VALUES "
	insertPattern := "(%d, %d, %d)"
	for i := 0; i < len(members); i++ {
		if i > 0 {
			insertSql += ","
		}
		m := members[i]
		insertSql += fmt.Sprintf(insertPattern, m.GroupId, m.UserId, m.Role)
	}
	err = conn.Exec(insertSql).Error
	return err
}

func ListMembersByGroupId(ctx context.Context, groupId int64) (members []*GroupMemberWithUser, err error) {
	conn, err := client.DBConn(ctx)
	if err != nil {
		return
	}
	defer conn.Close()

	conn = conn.Debug()
	conn = conn.Table("im_group_member gm").Joins("left join user_base u on gm.uid=u.uid")
	err = conn.Where("gm.group_id = ?", groupId).Find(&members).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			err = nil
			return
		}
	}
	return
}

func ListGroupByUserId(ctx context.Context, userId int64) (groups []*GroupMemberWithGroup, err error) {
	conn, err := client.DBConn(ctx)
	if err != nil {
		return
	}
	defer conn.Close()

	conn = conn.Debug()
	conn = conn.Table("im_group_member gm").Joins("left join im_group_base g on gm.group_id=g.group_id")
	err = conn.Where("gm.uid = ?", userId).Find(&groups).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			err = nil
			return
		}
	}
	return
}

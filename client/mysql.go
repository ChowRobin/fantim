package client

import (
	"context"
	"fmt"

	"github.com/gin-gonic/gin"

	"github.com/ChowRobin/fantim/constant"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
)

func DBConn(ctx context.Context) (*gorm.DB, error) {
	if ctx != nil {
		connItf := ctx.Value("DBConn")
		conn, ok := connItf.(*gorm.DB)
		if ok {
			return conn, nil
		}
	}
	return dbConn(
		constant.MYSQL_CONF.UserName,
		constant.MYSQL_CONF.Password,
		constant.MYSQL_CONF.Host,
		constant.MYSQL_CONF.DbName,
		constant.MYSQL_CONF.Port,
	)
}

func dbConn(MyUser, Password, Host, Db string, Port int32) (*gorm.DB, error) {
	connArgs := fmt.Sprintf("%s:%s@(%s:%d)/%s?charset=utf8&parseTime=True&loc=Local", MyUser, Password, Host, Port, Db)
	db, err := gorm.Open("mysql", connArgs)
	if err != nil {
		return nil, err
	}
	db.SingularTable(true)
	return db, err
}

func StartDBTransaction(ctx *gin.Context) error {
	db, err := DBConn(ctx)
	if err != nil {
		return err
	}
	tx := db.Begin()
	ctx.Set("DBConn", tx)
	return nil
}

func CommitDBTransaction(ctx context.Context) error {
	conn, err := DBConn(ctx)
	if err != nil {
		return err
	}
	conn.Commit()
	conn.Close()
	return nil
}

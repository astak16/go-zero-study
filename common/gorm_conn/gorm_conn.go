package gorm_conn

import (
	"fmt"
	"go-zero-study/model"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func GormConn(MysqlDataSource string) *gorm.DB {
	db, err := gorm.Open(mysql.Open(MysqlDataSource), &gorm.Config{})
	if err != nil {
		panic("数据库连接失败，" + err.Error())
	}
	fmt.Println("数据库连接成功")
	db.AutoMigrate(&model.User{})
	return db
}

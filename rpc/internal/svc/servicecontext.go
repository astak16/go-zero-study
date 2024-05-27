package svc

import (
	"go-zero-study/common/gorm_conn"
	"go-zero-study/rpc/internal/config"

	"gorm.io/gorm"
)

type ServiceContext struct {
	Config config.Config
	DB     *gorm.DB
}

func NewServiceContext(c config.Config) *ServiceContext {
	return &ServiceContext{
		Config: c,
		DB:     gorm_conn.GormConn(c.Mysql.DataSource),
	}
}

package model

import "gorm.io/gorm"

type User struct {
	gorm.Model
	ID       uint32 `gorm:"index,primarykey" json:"id"`
	Username string `gorm:"size:32" json:"username"`
	Password string `gorm:"size:64" json:"password"`
}

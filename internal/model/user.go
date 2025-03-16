package model

import (
	"gorm.io/gorm"
	"gorm.io/plugin/optimisticlock"
)

type User struct {
	gorm.Model
	UserID string `gorm:"type:varchar(150);column:user_id;unique;not null;comment:User ID"`
	UserName string `gorm:"type:varchar(32);column:username;unique;not null; comment:User Name"`
	Avatar string `gorm:"type:varchar(255);column:avatar;default:null;comment:Avatar"`
	Email string `gorm:"type:varchar(150);column:email;unique;not null;comment:Email"`
	Password string `gorm:"type:varchar(128);column:password;not null;comment:Password"`
	Status string `gorm:"type:varchar(15);column:status;default:offline;comment:User Status"`
	Version optimisticlock.Version
}

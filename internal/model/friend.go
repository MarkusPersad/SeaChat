package model

import (
	"gorm.io/gorm"
	"gorm.io/plugin/optimisticlock"
)

type Friend struct {
	gorm.Model
	FGID   string `gorm:"type:varchar(150);column:fg_id;unique;not null;comment:Friend or Group ID"`
	UserID string `gorm:"type:varchar(150);column:user_id;unique;not null;comment:User ID"`
	FriendID string `gorm:"type:varchar(150);column:friend_id;unique;not null;comment:Friend ID"`
	Status string `gorm:"type:varchar(15);column:status;not null;default:'friendly';comment:Friend Status"`
	Version optimisticlock.Version
}

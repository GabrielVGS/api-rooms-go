package models

import (
	"time"

	"gorm.io/gorm"
)

type Reservation struct {
	gorm.Model
	UserID    uint `gorm:"foreignKey:UserID"`
	RoomID    uint `gorm:"foreignKey:RoomID"`
	StartTime time.Time
	EndTime   time.Time
}

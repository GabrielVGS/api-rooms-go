package models

import "gorm.io/gorm"

type RoomMember struct {
	gorm.Model
	UserID uint `json:"user_id" gorm:"primaryKey"`
	RoomID uint `json:"room_id" gorm:"primaryKey"`
	Role   string `json:"role" gorm:"default:'member'"` // member, admin
	User   User `json:"user"`
	Room   Room `json:"room"`
}
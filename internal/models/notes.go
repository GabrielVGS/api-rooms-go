package models

import "gorm.io/gorm"

type Note struct {
	gorm.Model
	UserID  uint   `json:"user_id"`
	RoomID  uint   `json:"room_id"`
	Title   string `json:"title"`
	Content string `json:"content" gorm:"type:text"`
	User    User   `json:"user"`
	Room    Room   `json:"room"`
}
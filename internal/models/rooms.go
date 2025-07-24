package models

import "gorm.io/gorm"

type Room struct {
	gorm.Model
	Name        string       `json:"name"`
	Description string       `json:"description"`
	Subject     string       `json:"subject"`
	Capacity    int          `json:"capacity"`
	CreatedBy   uint         `json:"created_by"`
	Members     []RoomMember `json:"members" gorm:"foreignKey:RoomID"`
	Notes       []Note       `json:"notes" gorm:"foreignKey:RoomID"`
}

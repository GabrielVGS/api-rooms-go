package models

import "gorm.io/gorm"

type Room struct {
	gorm.Model
	Name        string `json:"name"`
	Description string `json:"description"`
	Capacity    int    `json:"capacity"`
}

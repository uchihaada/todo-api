package models

import "gorm.io/gorm"

type Todo struct {
	gorm.Model
	Title     string `json:"title" gorm:"not null"`          // Title is required
	Completed bool   `json:"completed" gorm:"default:false"` // Default value is false
}

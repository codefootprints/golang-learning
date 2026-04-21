package models

import "gorm.io/gorm"

type Task struct {
	gorm.Model
	Title  string `json:"title" binding:"required"`
	Status string `json:"status" binding:"required,oneof=Pending Done"`
}

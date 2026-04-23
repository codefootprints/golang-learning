package models

import "gorm.io/gorm"

type Task struct {
	gorm.Model
	Title  string `json:"title" binding:"required" gorm:"uniqueIndex"`
	Status string `json:"status" binding:"required,oneof=Pending Done"`
	UserID uint   `json:"user_id"` // Foreign Key ke tabel User

	// Tambahkan ini agar GORM tahu Task ini milik User
	User User `gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE" json:"-"`

	// Override DeletedAt agar tidak muncul di JSON
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}

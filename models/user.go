package models

import "gorm.io/gorm"

type User struct {
	gorm.Model
	Username string `gorm:"unique;not null" json:"username"`
	Password string `gorm:"not null" json:"-"`                        // "-" agar password tidak pernah tampil d JSON
	Tasks    []Task `gorm:"foreignKey:UserID" json:"tasks,omitempty"` // Relasi ke Task
}

package config

import (
	"fmt"
	"my-project/models"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var DB *gorm.DB

func InitDB() {
	dsn := "host=db user=user password=password dbname=taskdb port=5432 sslmode=disable"
	var err error
	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})

	if err != nil {
		panic("Gagal konek ke database")
	}
	fmt.Println("Berhasil konek ke database!")
	DB.AutoMigrate(&models.Task{})
}

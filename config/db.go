package config

import (
	"fmt"
	"my-project/models"
	"os"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var DB *gorm.DB

func InitDB() {
	dbHost := os.Getenv("DB_HOST")
	if dbHost == "" {
		dbHost = "localhost"
	}

	dbUser := os.Getenv("DB_USER")
	if dbUser == "" {
		{
			dbUser = "user"
		}
	}

	dsn := fmt.Sprintf("host=%s user=%s password=password dbname=taskdb port=5432 sslmode=disable", dbHost, dbUser)

	var err error
	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})

	if err != nil {
		panic("Gagal konek ke database:" + err.Error())
	}
	fmt.Printf("Berhasil konek ke database (Host: %s)\n!", dbHost)
	DB.AutoMigrate(&models.Task{}, &models.User{})

	DB.Exec("CREATE UNIQUE INDEX IF NOT EXISTS idx_title_active ON tasks (title) WHERE deleted_at IS NULL")
}

package main

import (
	"my-project/config"
	"my-project/handlers"
	"my-project/models"

	"github.com/gin-gonic/gin" // Framework API populer
)

// Pindahkan logika ini ke luar main() agar bisa dipanggil oleh file test
func CreateTaskLogic(task models.Task) (models.Task, error) {
	result := config.DB.Create(&task)
	return task, result.Error
}

// Fungsi logika untuk menghapus task
func DeleteTaskLogic(id string) (int64, error) {
	result := config.DB.Delete(&models.Task{}, id)
	return result.RowsAffected, result.Error
}

func GetTaskLogic(limit int) ([]models.Task, error) {
	var tasks []models.Task
	// GORM akan menambahkan "LIMIT" pada query SQL-nya
	result := config.DB.Limit(limit).Find(&tasks)
	return tasks, result.Error
}

func main() {
	config.InitDB()

	router := gin.Default()

	// Endpoint untuk mengambil semua data (GET)
	router.GET("/tasks", handlers.GetTasks)

	// Endpoint untuk mengambil satu data berdasarkan ID (GET)
	router.GET("/tasks/:id", handlers.GetTaskById)

	// Endpoint untuk menambah data baru (POST)
	router.POST("/tasks", handlers.CreateTask)

	// Endpoint untuk update satu task
	router.PUT("/tasks/:id", handlers.UpdateTask)

	// Endpoint untuk menghapus data (DELETE)
	router.DELETE("/tasks/:id", handlers.DeleteTask)

	// Endpoint untuk menghapus SEMUA data
	router.DELETE("/tasks", handlers.DeleteAllTasks)

	router.Run(":8080") // Jalankan di port 8080
}

package main

import (
	"my-project/config"
	"my-project/handlers"
	"my-project/middlewares"
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

	// Rute publik (bisa diakses tanpa login)
	router.POST("/register", handlers.Register)
	router.POST("/login", handlers.Login)

	// Kelompok rute yang butuh Middleware
	authorized := router.Group("/")
	authorized.Use(middlewares.AuthMiddleware())
	{
		authorized.GET("/tasks", handlers.GetTasks)
		authorized.GET("/tasks/:id", handlers.GetTaskById)
		authorized.POST("/tasks", handlers.CreateTask)
		authorized.PUT("/tasks/:id", handlers.UpdateTask)
		authorized.DELETE("/tasks/:id", handlers.DeleteTask)
		authorized.DELETE("/tasks", handlers.DeleteAllTasks)
	}

	router.Run(":8080") // Jalankan di port 8080
}

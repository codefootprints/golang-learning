package handlers

import (
	"fmt"
	"my-project/config"
	"my-project/models"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// Menampilkan semua task (dengan limit)
func GetTasks(c *gin.Context) {
	limitStr := c.DefaultQuery("limit", "10")
	limit, _ := strconv.Atoi(limitStr)
	fmt.Println("hello")

	var tasks []models.Task
	config.DB.Limit(limit).Find(&tasks)
	c.JSON(http.StatusOK, tasks)
}

// Membuat task baru
func CreateTask(c *gin.Context) {
	var task models.Task
	if err := c.ShouldBindJSON(&task); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Validasi gagal!",
			"details": err.Error(),
		})
		return
	}
	config.DB.Create(&task)
	c.JSON(http.StatusCreated, task)
}

// Menghapus task berdasarkan ID
func DeleteTask(c *gin.Context) {
	id := c.Param("id")

	// Kita gunakan unscoped untuk Hard Delete
	result := config.DB.Unscoped().Delete(&models.Task{}, id)

	if result.RowsAffected == 0 {
		c.JSON(http.StatusNotFound, gin.H{"message": "Task tidak ditemukan"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Task berhasil dihapus"})
}

func GetTaskById(c *gin.Context) {
	id := c.Param("id")

	var task models.Task

	if err := config.DB.First(&task, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Task tidak ditemukan"})
		return
	}

	c.JSON(http.StatusOK, task)
}

func UpdateTask(c *gin.Context) {
	id := c.Param("id")
	var task models.Task

	if err := config.DB.First(&task, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Task tidak ditemukan"})
		return
	}

	var updateData models.Task
	if err := c.ShouldBindJSON(&updateData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	config.DB.Model(&task).Updates(updateData)

	c.JSON(http.StatusOK, task)
}

func DeleteAllTasks(c *gin.Context) {
	config.DB.Where("1 = 1").Unscoped().Delete(&models.Task{})
	c.JSON(http.StatusOK, gin.H{"message": "Semua task telah dibersihkan!"})
}

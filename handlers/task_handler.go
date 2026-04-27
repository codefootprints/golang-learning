package handlers

import (
	"fmt"
	"my-project/config"
	"my-project/models"
	"my-project/repositories"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// Inisialisasi repository (nanti ini bisa dipindah ke main)
var taskRepo = repositories.TaskRepository{}

type UpdateTaskInput struct {
	Title  string `json:"title"`
	Status string `json:"status" binding:"oneof=Pending Done"`
}

// Menampilkan semua task (dengan limit)
func GetTasks(c *gin.Context) {
	val, exists := c.Get("currentUserID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "Identitas user tidak ditemukan",
		})
		return
	}
	userID := val.(uint)

	// 1. Ambil input
	pageStr := c.DefaultQuery("page", "1")
	limitStr := c.DefaultQuery("limit", "5")
	sort := c.DefaultQuery("sort", "id")
	order := c.DefaultQuery("order", "desc") // Default diperbaiki ke 'desc'
	search := c.Query("search")
	status := c.Query("status")

	// 2. Validasi Parameter (Guard Clause)
	if (sort != "id" && sort != "title") || (order != "asc" && order != "desc") {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Parameter invalid. Use sort: id/title and order: asc/desc.",
		})
		return
	}

	// 3. Konversi data
	page, _ := strconv.Atoi(pageStr)
	limit, _ := strconv.Atoi(limitStr)

	// Sanitize
	if page < 1 {
		page = 1
	}

	if limit < 1 {
		limit = 5
	} else if limit > 100 {
		limit = 100
	}

	offset := (page - 1) * limit

	var tasks []models.Task
	query := config.DB.
		Model(&models.Task{}).
		Where("user_id = ?", userID)

	// 4. Logika Bisnis (Filter & Search)
	if search != "" {
		query = query.Where("title ILIKE ?", "%"+search+"%")
	}
	if status != "" {
		query = query.Where("status = ?", status)
	}

	// 5. Eksekusi Database (Method Chaining)
	result := query.Limit(limit).
		Offset(offset).
		Order(fmt.Sprintf("%s %s", sort, order)).
		Find(&tasks)

	// 6. Response Handling
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal mengambil data"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"page":  page,
		"limit": limit,
		"count": len(tasks), // Tambahan opsional untuk memudahkan pengecekan
		"data":  tasks,
	})
}

// Membuat task baru
func CreateTask(c *gin.Context) {
	// Ambil ID User dari context (hasil kerja AuthMiddleware)
	val, exists := c.Get("currentUserID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "Identitas user tidak ditemukan",
		})
		return
	}
	userID := val.(uint)

	var task models.Task
	if err := c.ShouldBindJSON(&task); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Validasi gagal!",
			"details": err.Error(),
		})
		return
	}
	task.UserID = userID

	result := config.DB.Create(&task)

	if result.Error != nil {
		c.JSON(http.StatusConflict, gin.H{
			"error":   "Gagal membuat task",
			"details": "Judul ini sudah ada dan masih aktif. Gunakan judul lain.",
		})
		return
	}

	c.JSON(http.StatusCreated, task)
}

// Menghapus task berdasarkan ID
func DeleteTask(c *gin.Context) {
	// Ambil ID User dari context (hasil kerja AuthMiddleware)
	val, exists := c.Get("currentUserID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "Identitas user tidak ditemukan",
		})
		return
	}
	userID := val.(uint)

	taskId := c.Param("id")

	// Tambahkan DB.Debug() jika ingin melihat plain query string
	result := config.DB.Where("user_id = ? AND id = ?", userID, taskId).Delete(&models.Task{})

	if result.RowsAffected == 0 {
		c.JSON(http.StatusNotFound, gin.H{"message": "Task tidak ditemukan"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Task berhasil dihapus"})
}

func GetTaskById(c *gin.Context) {
	currentUserID, exists := c.Get("currentUserID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "Identitas user tidak ditemukan",
		})
		return
	}
	id := c.Param("id")

	task, err := taskRepo.GetByID(id, currentUserID.(uint))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Task tidak ditemukan"})
		return
	}

	c.JSON(http.StatusOK, task)
}

func UpdateTask(c *gin.Context) {
	currentUserID, _ := c.Get("currentUserID")
	id := c.Param("id")

	// Cari data lewat repository
	task, err := taskRepo.GetByID(id, currentUserID.(uint))

	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Task tidak ditemukan"})
		return
	}

	// Ambil input DTO (Data Transfer Object)
	var updateData UpdateTaskInput
	if err := c.ShouldBindJSON(&updateData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Update lewat repository
	err = taskRepo.Update(&task, models.Task{
		Title:  updateData.Title,
		Status: updateData.Status,
	})

	if err != nil {
		// Membungkus Error
		wrappedErr := fmt.Errorf("Update task failed for ID %s: %w", id, err)

		// Cetak di log server
		fmt.Println(wrappedErr)

		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Gagal update data",
		})
		return
	}

	c.JSON(http.StatusOK, task)
}

func DeleteAllTasks(c *gin.Context) {
	// Perintah GORM untuk menghapus semua record di tabel tasks
	// Kita pakai Where("1 = 1") karena GORM mencegah penghapusan tanpa kondisi (Safety Feature)
	// Kita gunakan Unscoped() untuk Hard Delete
	result := config.DB.Where("1 = 1").Unscoped().Delete(&models.Task{})

	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal menghapus semua data"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": fmt.Sprintf("Berhasil menghapus seluruh data (%d task terhapus)", result.RowsAffected),
	})
}

package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"my-project/config"
	"my-project/handlers"
	"my-project/models"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/gin-gonic/gin"
)

func TestMain(m *testing.M) {
	// Jalankan initDB satu kali saja di awal
	config.InitDB()

	// Jalankan semua test (TestCreate..., TestDelete...)
	exitCode := m.Run()

	// Keluar dengan code yang sesuai
	os.Exit(exitCode)
}

func TestCreateTaskLogic(t *testing.T) {
	// 1. Setup data dummy
	testTask := models.Task{
		Title:  "Cleanup Test Task",
		Status: "Testing",
	}

	// 2. Eksekusi
	created, err := CreateTaskLogic(testTask)
	if err != nil {
		t.Fatalf("Gagal membuat task: %v", err)
	}

	// --- BAGIAN CLEANUP ---
	// Perintah ini akan dijalankan OTOMATIS setelah fungsi test selesai
	t.Cleanup(func() {
		config.DB.Unscoped().Delete(&created)
		// Unscoped() digunakan agar data benar-benar terhapus dari tabel (bukan soft delete)
	})
	// ----------------------

	// 3. Assert (Validasi)
	if created.ID == 0 {
		t.Error("Harusnya ID tidak nol setelah disimpan")
	}

	if created.Title != testTask.Title {
		t.Errorf("Ekspektasi %s, dapat %s", testTask.Title, created.Title)
	}
}

func TestDeleteTaskLogic(t *testing.T) {
	// 1. Buat data dulu untuk dihapus
	task, _ := CreateTaskLogic(models.Task{Title: "Task Mau Dihapus", Status: "Temp"})
	idStr := fmt.Sprintf("%d", task.ID)

	// 2. Jalankan fungsi Delete
	rows, err := DeleteTaskLogic(idStr)

	// 3. Validasi (Assert)
	if err != nil {
		t.Errorf("Harusnya tidak error saat delete, tapi dapat: %v", err)
	}
	if rows == 0 {
		t.Error("Harusnya ada 1 baris yang terhapus, tapi dapat 0")
	}

	// 4. Double Check: Pastikan sudah tidak ada di DB
	var checkTask models.Task
	result := config.DB.First(&checkTask, task.ID)
	if result.Error == nil {
		t.Error("Data masih ditemukan di DB, harusnya sudah terhapus!")
	}
}

func TestCreateTaskValidation(t *testing.T) {
	// 1. Tambahkan ini untuk membungkam log debug Gin
	gin.SetMode(gin.ReleaseMode)

	// Buat router khusus untuk test ini
	// router := gin.Default()

	// 2. Gunakan gin.New() agar tidak ada log request otomatis
	router := gin.New()
	router.Use(gin.Recovery()) // Agar tetap aman jika ada panic

	router.POST("/tasks", handlers.CreateTask)

	// Skenario: Mengirim Status yang salah (bukan Pending/Done)
	invalidData := map[string]string{
		"title":  "Coba Validasi",
		"status": "Ngawur",
	}
	body, _ := json.Marshal(invalidData)

	req, _ := http.NewRequest("POST", "/tasks", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// Assert: Kita ingin statusnya 400 (Bad Request)
	if w.Code != http.StatusBadRequest {
		t.Errorf("Harusnya ditolak (400), tapi malah lolos (%d)", w.Code)
	}
}

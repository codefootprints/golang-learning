package main

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin" // Framework API populer
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// Ini adalah "Blueprint" data kita (Struct)
type Task struct {
	ID     uint   `gorm:"primaryKey;autoIncrement" json:"id"`
	Title  string `json:"title"`
	Status string `json:"status"`
}

var db *gorm.DB

func initDB() {
	dsn := "host=db user=user password=password dbname=taskdb port=5432 sslmode=disable"
	var err error
	for i := 0; i < 10; i++ {
		db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
		if err == nil {
			fmt.Println("Berhasil konek ke database!")
			break
		}
		fmt.Printf("Database belum siap (percobaan %d)... menunggu 2 detik \n", i+1)
		time.Sleep(2 * time.Second)
	}

	if err != nil {
		panic("Gagal konek ke database setelah beberapa kali mencoba")
	}

	db.AutoMigrate(&Task{}) // Otomatis bikin tabel
}

func main() {
	initDB()
	router := gin.Default()

	// Endpoint untuk mengambil semua data (GET)
	router.GET("/tasks", func(c *gin.Context) {
		var tasks []Task
		db.Find(&tasks)
		c.JSON(http.StatusOK, tasks)
	})

	// Endpoint untuk menambah data baru (POST)
	router.POST("/tasks", func(c *gin.Context) {
		var newTask Task
		if err := c.ShouldBindJSON(&newTask); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		db.Create(&newTask) // Simpan permanen ke Postgres
		c.JSON(http.StatusCreated, newTask)
	})

	router.Run(":8080") // Jalankan di port 8080
}

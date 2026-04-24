package handlers

import (
	"fmt"
	"my-project/config"
	"my-project/models"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

func Register(c *gin.Context) {
	var input models.User
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	// Hash Password sebelum simpan
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)
	fmt.Printf("DB Password: [%s]\n", hashedPassword)
	fmt.Printf("Plain Password: [%s]\n", input.Password)
	input.Password = string(hashedPassword)

	if err := config.DB.Create(&input).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Gagal membuat user",
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "Registrasi berhasil",
	})
}

func Login(c *gin.Context) {
	var input struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Input tidak valid",
		})
		return
	}

	// Cari user di database
	var user models.User
	if err := config.DB.Where("username = ?", input.Username).First(&user).Error; err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "Username atau password salah",
		})
		return
	}

	fmt.Printf("Input Password: [%s]\n", input.Password)
	fmt.Printf("DB Password: [%s]\n", user.Password)

	// Cek password (Bcrypt)
	err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(input.Password))
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "Username atau password salah",
		})
		return
	}

	// Generate JWT
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id":  user.ID,
		"username": user.Username,
		"exp":      time.Now().Add(time.Hour * 24).Unix(),
	})

	tokenString, _ := token.SignedString([]byte("RAHASIA_NEGARA")) // Gunakan secret key yang kuat

	c.JSON(http.StatusOK, gin.H{
		"token": tokenString,
	})
}

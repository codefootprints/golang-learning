package middlewares

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

// eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3NzcwMjE2ODMsInVzZXJfaWQiOjEsInVzZXJuYW1lIjoiYnVkaSJ9.xNbPaodDFjrHsjGTLwbMJklTCU8d-ijs7rsGxn1GOlY

var jwtSecret = []byte("RAHASIA_NEGARA")

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Ambil header Authorization
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": "Belum login",
			})
			return
		}

		// Format biasanya: "Bearer <token>"
		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": "Format token salah",
			})
			return
		}

		tokenString := parts[1]

		// Validasi token
		token, err := jwt.Parse(tokenString, func(t *jwt.Token) (interface{}, error) {
			// Pastikan metode signingnya benar (HS256)
			if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("Metode signing salah: %v", t.Header["alg"])
			}
			return jwtSecret, nil
		})

		// Cek apakah token valid dan ambil isinya (Claims)
		if err != nil || !token.Valid {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": "Token tidak valid atau kadaluarsa",
			})
			return
		}

		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": "Gagal membedah data token",
			})
			return
		}

		// Simpaan user_id ke Context agar bisa dipanggil di Handler
		// JWT menyimpan angka sebagai float64, kita konversi ke unit
		userID := uint(claims["user_id"].(float64))
		c.Set("currentUserID", userID)

		// Lanjut ke proses berikutnya
		c.Next()
	}
}

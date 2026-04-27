package repositories

import (
	"my-project/config"
	"my-project/models"

	"gorm.io/gorm"
)

type UserRepository struct{}

// FindByUsername digunakan untuk proses Login
func (r *UserRepository) FindByUsername(username string) (models.User, error) {
	var user models.User
	err := config.DB.Where("username = ?", username).First(&user).Error
	return user, err
}

// Create digunakan untuk proses Register
func (r *UserRepository) Create(user *models.User) error {
	return config.DB.Create(user).Error
}

// GetByID untuk mengambil profil atau verifikasi poin
func (r *UserRepository) GetByID(id uint) (models.User, error) {
	var user models.User
	err := config.DB.First(&user, id).Error
	return user, err
}

// AddPoints: Contoh fungsi yang akan dipanggil oleh Service nanti
func (r *UserRepository) AddPoints(id uint, points int) error {
	return config.DB.Model(&models.User{}).Where("id = ?", id).
		Update("points", gorm.Expr("points + ?", points)).Error
}

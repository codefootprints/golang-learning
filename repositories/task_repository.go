package repositories

import (
	"my-project/config"
	"my-project/models"
)

// TaskRepository mengelola semua interaksi database untuk Task
type TaskRepository struct{}

func (r *TaskRepository) GetByID(id string, userID uint) (models.Task, error) {
	var task models.Task
	err := config.DB.Where("user_id = ?", userID).First(&task, id).Error
	return task, err
}

func (r *TaskRepository) Update(task *models.Task, data interface{}) error {
	return config.DB.Model(task).Updates(data).Error
}

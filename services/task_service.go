package services

import "my-project/repositories"

type TaskService struct {
	taskRepo repositories.TaskRepository
	userRepo repositories.UserRepository
}

func (s *TaskService) CompleteTask(taskID string, userID uint) error {
	// Update status task jadi Done via TaskRepo
	err := s.taskRepo.UpdateStatus(taskID, userID, "Done")
	if err != nil {
		return err
	}

	// Tambah point user via UserRepo
	err = s.userRepo.AddPoints(userID, 10)
	if err != nil {
		return err
	}

	return nil
}

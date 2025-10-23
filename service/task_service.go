package service

import (
	"errors"

	"github.com/lieucongduy182/go-gin-todo-api/models"
	"github.com/lieucongduy182/go-gin-todo-api/repository"
)

type TaskService interface {
	GetTask(taskID, userID uint) (*models.TaskResponse, error)
	GetUserTask(userID uint) ([]*models.TaskResponse, error)
	GetUserTaskPaginated(userID uint, page, pageSize int) (*models.PaginationResponse, error)
	SearchTasks(userID uint, keyword string) ([]*models.TaskResponse, error)
	GetTaskByStats(userID uint) (map[string]interface{}, error)
	GetTasksByPriority(userID uint, priority string) ([]models.TaskResponse, error)
	GetTasksByStatus(userID uint, completed bool) ([]models.TaskResponse, error)
	CreateTask(userID uint, req *models.CreateTaskRequest) (*models.TaskResponse, error)
	UpdateTask(taskID, userID uint, req *models.UpdateTaskRequest) (*models.TaskResponse, error)
	DeleteTask(taskID, userID uint) error
}

type taskService struct {
	userRepo repository.UserRepository
	taskRepo repository.TaskRepository
}

func isValidPriority(priority string) bool {
	validProperties := map[string]bool{
		"low":    true,
		"medium": true,
		"high":   true,
	}
	return validProperties[priority]
}

// CreateTask implements TaskService.
func (t *taskService) CreateTask(userID uint, req *models.CreateTaskRequest) (*models.TaskResponse, error) {
	task := &models.Task{
		UserID:      userID,
		Title:       req.Title,
		Description: req.Description,
		Priority:    req.Priority,
		DueDate:     req.DueDate,
		Completed:   false,
	}

	if task.Priority == "" {
		task.Priority = "medium"
	}

	if !isValidPriority(task.Priority) {
		return nil, errors.New("invalid priority values")
	}

	if err := t.taskRepo.Create(task); err != nil {
		return nil, err
	}

	response := task.ToResponse()
	return &response, nil
}

// DeleteTask implements TaskService.
func (t *taskService) DeleteTask(taskID uint, userID uint) error {
	return t.taskRepo.Delete(taskID, userID)
}

// GetTask implements TaskService.
func (t *taskService) GetTask(taskID uint, userID uint) (*models.TaskResponse, error) {
	task, err := t.taskRepo.GetById(taskID, userID)
	if err != nil {
		return nil, err
	}

	response := task.ToResponse()
	return &response, nil
}

// GetTaskByStats implements TaskService.
func (t *taskService) GetTaskByStats(userID uint) (map[string]interface{}, error) {
	return t.taskRepo.GetStats(userID)
}

// GetTasksByPriority implements TaskService.
func (t *taskService) GetTasksByPriority(userID uint, priority string) ([]models.TaskResponse, error) {
	if !isValidPriority(priority) {
		return nil, errors.New("invalid priority values")
	}

	tasks, err := t.taskRepo.GetByPriority(userID, priority)
	if err != nil {
		return nil, err
	}

	var responses []models.TaskResponse
	for _, task := range tasks {
		responses = append(responses, task.ToResponse())
	}

	return responses, nil
}

// GetTasksByStatus implements TaskService.
func (t *taskService) GetTasksByStatus(userID uint, completed bool) ([]models.TaskResponse, error) {
	tasks, err := t.taskRepo.GetByStatus(userID, completed)
	if err != nil {
		return nil, err
	}

	var responses []models.TaskResponse
	for _, task := range tasks {
		responses = append(responses, task.ToResponse())
	}

	return responses, nil
}

// GetUserTask implements TaskService.
func (t *taskService) GetUserTask(userID uint) ([]*models.TaskResponse, error) {
	tasks, err := t.taskRepo.GetByUserId(userID)
	if err != nil {
		return nil, err
	}

	var responses []*models.TaskResponse
	for _, task := range *tasks {
		t := task.ToResponse()
		responses = append(responses, &t)
	}

	return responses, nil
}

// GetUserTaskPaginated implements TaskService.
func (t *taskService) GetUserTaskPaginated(userID uint, page int, pageSize int) (*models.PaginationResponse, error) {
	tasks, total, err := t.taskRepo.GetByUserIdPaginated(userID, page, pageSize)
	if err != nil {
		return nil, err
	}

	var responses []models.TaskResponse
	for _, task := range tasks {
		responses = append(responses, task.ToResponse())
	}

	totalPages := (total + int64(pageSize) - 1) / int64(pageSize)

	return &models.PaginationResponse{
		Data:       responses,
		Total:      total,
		Page:       page,
		PageSize:   pageSize,
		TotalPages: totalPages,
	}, nil
}

// SearchTasks implements TaskService.
func (t *taskService) SearchTasks(userID uint, keyword string) ([]*models.TaskResponse, error) {
	task, err := t.taskRepo.Search(userID, keyword)
	if err != nil {
		return nil, err
	}

	var responses []*models.TaskResponse
	for _, t := range task {
		resp := t.ToResponse()
		responses = append(responses, &resp)
	}
	return responses, nil
}

// UpdateTask implements TaskService.
func (t *taskService) UpdateTask(taskID uint, userID uint, req *models.UpdateTaskRequest) (*models.TaskResponse, error) {
	task, err := t.taskRepo.GetById(taskID, userID)
	if err != nil {
		return nil, err
	}

	if req.Title != "" {
		task.Title = req.Title
	}

	if req.Description != "" {
		task.Description = req.Description
	}

	if req.Completed != nil {
		task.Completed = *req.Completed
	}

	if req.Priority != "" {
		if !isValidPriority(req.Priority) {
			return nil, errors.New("invalid priority values")
		}

		task.Priority = req.Priority
	}

	if req.DueDate != nil {
		task.DueDate = req.DueDate
	}

	if err := t.taskRepo.Update(task); err != nil {
		return nil, err
	}

	response := task.ToResponse()
	return &response, nil
}

func NewTaskService(
	userRepo repository.UserRepository,
	taskRepo repository.TaskRepository,
) TaskService {
	return &taskService{
		userRepo: userRepo,
		taskRepo: taskRepo,
	}
}

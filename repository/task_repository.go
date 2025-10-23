package repository

import (
	"errors"

	"github.com/lieucongduy182/go-gin-todo-api/models"
	"gorm.io/gorm"
)

type TaskRepository interface {
	Create(task *models.Task) error
	GetById(id, userID uint) (*models.Task, error)
	GetByUserId(userID uint) (*[]models.Task, error)
	GetByUserIdPaginated(userID uint, page, pageSize int) ([]*models.Task, int64, error)
	Search(userID uint, keyword string) ([]*models.Task, error)
	GetByStatus(userID uint, status bool) ([]*models.Task, error)
	GetByPriority(userID uint, priority string) ([]*models.Task, error)
	Update(task *models.Task) error
	Delete(id, userID uint) error
	DeleteByUserId(userID uint) error
	Count(userID uint) (int64, error)
	GetStats(userID uint) (map[string]interface{}, error)
}

// taskRepository implement The TaskRepository interface
type taskRepository struct {
	db *gorm.DB
}

// Count implements TaskRepository.
func (t *taskRepository) Count(userID uint) (int64, error) {
	var count int64
	if err := t.db.Model(&models.Task{}).Where("user_id = ?", userID).Count(&count).Error; err != nil {
		return 0, err
	}

	return count, nil
}

// Create implements TaskRepository.
func (t *taskRepository) Create(task *models.Task) error {
	if err := t.db.Save(task).Error; err != nil {
		return err
	}
	return nil
}

// Delete implements TaskRepository.
func (t *taskRepository) Delete(id uint, userID uint) error {
	if err := t.db.Where("id = ? and user_id = ?", id, userID).Delete(&models.Task{}).Error; err != nil {
		return err
	}

	return nil
}

// DeleteByUserId implements TaskRepository.
func (t *taskRepository) DeleteByUserId(userID uint) error {
	if err := t.db.Where("user_id = ?", userID).Delete(&models.Task{}).Error; err != nil {
		return err
	}

	return nil
}

// GetById implements TaskRepository.
func (t *taskRepository) GetById(id uint, userID uint) (*models.Task, error) {
	var task *models.Task
	if err := t.db.Where("id = ? and userID = ?", id, userID).First(&task).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("Task not Found")
		}

		return nil, err
	}

	return task, nil
}

// GetByPriority implements TaskRepository.
func (t *taskRepository) GetByPriority(userID uint, priority string) ([]*models.Task, error) {
	var task []*models.Task

	if err := t.db.Where("user_id = ? AND priority = ?", userID, priority).
		Order("created_at desc").
		First(&task).Error; err != nil {
		return nil, err
	}

	return task, nil
}

// GetByStatus implements TaskRepository.
func (t *taskRepository) GetByStatus(userID uint, status bool) ([]*models.Task, error) {
	var task []*models.Task

	if err := t.db.Where("user_id = ? AND completed = ?", userID, status).
		Order("created_at desc").
		First(&task).Error; err != nil {
		return nil, err
	}

	return task, nil
}

// GetByUserId implements TaskRepository.
func (t *taskRepository) GetByUserId(userID uint) (*[]models.Task, error) {
	var task *[]models.Task
	if err := t.db.Where("userID = ?", userID).
		Order("created_at desc").
		Find(&task).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("Task not Found")
		}

		return nil, err
	}

	return task, nil
}

// GetByUserIdPaginated implements TaskRepository.
func (t *taskRepository) GetByUserIdPaginated(userID uint, page int, pageSize int) ([]*models.Task, int64, error) {
	var task []*models.Task
	var total int64

	if err := t.db.Model(&models.Task{}).Where("user_id = ?", userID).
		Count(&total).Error; err != nil {
		return nil, 0, err
	}
	offset := (page - 1) * pageSize
	if err := t.db.Where("user_id = ?", userID).
		Offset(offset).
		Limit(pageSize).
		Find(&task).Error; err != nil {
		return nil, 0, err
	}

	return task, total, nil
}

// GetStats implements TaskRepository.
func (t *taskRepository) GetStats(userID uint) (map[string]interface{}, error) {
	var total, completed, pending, high, medium, low int64

	t.db.Model(&models.Task{}).Where("user_id = ?", userID).Count(&total)
	t.db.Model(&models.Task{}).Where("user_id = ? AND completed = ?", userID, true).Count(&completed)
	t.db.Model(&models.Task{}).Where("user_id = ? AND completed = ?", userID, false).Count(&pending)

	t.db.Model(&models.Task{}).Where("user_id = ? AND priority = ?", userID, "high").Count(&high)
	t.db.Model(&models.Task{}).Where("user_id = ? AND priority = ?", userID, "medium").Count(&medium)
	t.db.Model(&models.Task{}).Where("user_id = ? AND priority = ?", userID, "low").Count(&low)

	stats := map[string]interface{}{
		"total": total,
		"statuses": map[string]int64{
			"completed": completed,
			"pending":   pending,
		},
		"priorities": map[string]int64{
			"high":   high,
			"medium": medium,
			"low":    low,
		},
	}
	return stats, nil
}

// Search implements TaskRepository.
func (t *taskRepository) Search(userID uint, keyword string) ([]*models.Task, error) {
	var task []*models.Task
	searchPattern := "%" + keyword + "%"

	if err := t.db.Where("user_id = ? AND (title ILIKE ? OR description ILIKE ?)",
		userID, searchPattern, searchPattern).
		Order("created_at desc").
		Find(&task).Error; err != nil {
		return nil, err
	}

	return task, nil
}

// Update implements TaskRepository.
func (t *taskRepository) Update(task *models.Task) error {
	if err := t.db.Save(task).Error; err != nil {
		return err
	}

	return nil
}

func NewTaskRepository(db *gorm.DB) TaskRepository {
	return &taskRepository{db: db}
}

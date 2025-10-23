package repository

import (
	"errors"

	"github.com/lieucongduy182/go-gin-todo-api/models"
	"gorm.io/gorm"
)

// User Repository defines the interface
type UserRepository interface {
	Create(user *models.User) error
	GetByID(id uint) (*models.User, error)
	GetByEmail(email string) (*models.User, error)
	GetByUsername(username string) (*models.User, error)
	Update(user *models.User) error
	Delete(id uint) error
	UserExists(email string) (bool, error)
	GetAll(page, pageSize int) ([]models.User, int64, error)
}

// userRepository implement UserRepository interface
type userRepository struct {
	db *gorm.DB
}

// Create implements UserRepository.
func (u *userRepository) Create(user *models.User) error {
	if err := u.db.Create(user).Error; err != nil {
		return err
	}

	return nil
}

// Delete implements UserRepository.
func (u *userRepository) Delete(id uint) error {
	if err := u.db.Delete(id).Error; err != nil {
		return err
	}

	return nil
}

// GetAll implements UserRepository.
func (u *userRepository) GetAll(page int, pageSize int) ([]models.User, int64, error) {
	var users []models.User
	var total int64

	if err := u.db.Model(&models.User{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * pageSize

	if err := u.db.Offset(offset).Limit(pageSize).Find(&users).Error; err != nil {
		return nil, 0, err
	}

	return users, total, nil
}

// GetByEmail implements UserRepository.
func (u *userRepository) GetByEmail(email string) (*models.User, error) {
	var user *models.User

	if err := u.db.Where("email = ?", email).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("user not found")
		}
		return nil, err
	}

	return user, nil
}

// GetByID implements UserRepository.
func (u *userRepository) GetByID(id uint) (*models.User, error) {
	var user models.User

	if err := u.db.Where("id = ?", id).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("user not found")
		}
		return nil, err
	}

	return &user, nil
}

// GetByUsername implements UserRepository.
func (u *userRepository) GetByUsername(username string) (*models.User, error) {
	var user models.User

	if err := u.db.Where("username = ?", username).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("user not found")
		}
		return nil, err
	}

	return &user, nil
}

// Update implements UserRepository.
func (u *userRepository) Update(user *models.User) error {
	if err := u.db.Save(user).Error; err != nil {
		return err
	}

	return nil
}

// UserExists implements UserRepository.
func (u *userRepository) UserExists(email string) (bool, error) {
	var count int64
	if err := u.db.Where("email = ?", email).Count(&count).Error; err != nil {
		return false, err
	}

	return count > 0, nil
}

func NewUserRepository(db *gorm.DB) UserRepository {
	return &userRepository{db: db}
}

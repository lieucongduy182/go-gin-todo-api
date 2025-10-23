package service

import (
	"errors"

	"github.com/lieucongduy182/go-gin-todo-api/models"
	"github.com/lieucongduy182/go-gin-todo-api/repository"
	"github.com/lieucongduy182/go-gin-todo-api/utils"
)

// UserService interface the business logic for users
type UserService interface {
	Register(req *models.RegisterRequest) (*models.AuthResponse, error)
	Login(req *models.LoginRequest) (*models.AuthResponse, error)
	GetProfile(userID uint) (*models.UserResponse, error)
	ChangePassword(userID uint, req *models.ChangePasswordRequest) error
	UpdateProfile(userID uint, username string) error
	DeleteAccount(userID uint) error
	GetAllUsers(page, pageSize int) (*models.PaginationResponse, error)
}

type userService struct {
	userRepo repository.UserRepository
	taskRepo repository.TaskRepository
}

// ChangePassword implements UserService.
func (u *userService) ChangePassword(userID uint, req *models.ChangePasswordRequest) error {
	user, err := u.userRepo.GetByID(userID)
	if err != nil {
		return err
	}

	if err := user.CheckPassword(req.OldPassword); err != nil {
		return errors.New("current password incorrect")
	}

	if err := user.HashPassword(req.NewPassword); err != nil {
		return errors.New("failed to hash password")
	}

	return u.userRepo.Update(user)
}

// DeleteAccount implements UserService.
func (u *userService) DeleteAccount(userID uint) error {
	if err := u.taskRepo.DeleteByUserId(userID); err != nil {
		return err
	}

	return u.userRepo.Delete(userID)
}

// GetAllUsers implements UserService.
func (u *userService) GetAllUsers(page int, pageSize int) (*models.PaginationResponse, error) {
	users, total, err := u.userRepo.GetAll(page, pageSize)
	if err != nil {
		return nil, err
	}

	// convert to response DTOs
	var responses []models.UserResponse
	for _, user := range users {
		responses = append(responses, user.ToResponse())
	}

	totalPages := (total + int64(pageSize) - 1) / int64(pageSize)

	return &models.PaginationResponse{
		Data:       responses,
		Total:      total,
		TotalPages: totalPages,
		Page:       page,
		PageSize:   pageSize,
	}, nil
}

// GetProfile implements UserService.
func (u *userService) GetProfile(userID uint) (*models.UserResponse, error) {
	user, err := u.userRepo.GetByID(userID)
	if err != nil {
		return nil, err
	}

	response := user.ToResponse()
	return &response, nil
}

// Login implements UserService.
func (u *userService) Login(req *models.LoginRequest) (*models.AuthResponse, error) {
	user, err := u.userRepo.GetByEmail(req.Email)
	if err != nil {
		return nil, errors.New("invalid email or password")
	}

	if err := user.CheckPassword(req.Password); err != nil {
		return nil, errors.New("invalid email or password")
	}

	token, err := utils.GenerateToken(user.ID, user.Email)
	if err != nil {
		return nil, errors.New("Failed to generate token")
	}

	return &models.AuthResponse{
		Token:     token,
		User:      user.ToResponse(),
		ExpiresIn: 86400,
		TokenType: "Bearer",
	}, nil
}

// Register implements UserService.
func (u *userService) Register(req *models.RegisterRequest) (*models.AuthResponse, error) {
	exists, err := u.userRepo.UserExists(req.Email)
	if err != nil {
		return nil, err
	}

	if exists {
		return nil, errors.New("User with this email existing")
	}

	user := &models.User{
		Username: req.Username,
		Email:    req.Email,
	}

	if err := user.HashPassword(req.Password); err != nil {
		return nil, errors.New("Failed to hash password")
	}

	if err := u.userRepo.Create(user); err != nil {
		return nil, errors.New("Failed to create user")
	}

	token, err := utils.GenerateToken(user.ID, user.Email)
	if err != nil {
		return nil, errors.New("Failed to generate token")
	}

	return &models.AuthResponse{
		Token:     token,
		User:      user.ToResponse(),
		ExpiresIn: 86400, // 24 hours
		TokenType: "Bearer",
	}, nil
}

// UpdateProfile implements UserService.
func (u *userService) UpdateProfile(userID uint, username string) error {
	user, err := u.userRepo.GetByID(userID)
	if err != nil {
		return err
	}

	user.Username = username
	return u.userRepo.Update(user)
}

func NewUserService(
	userRepo repository.UserRepository,
	taskRepo repository.TaskRepository,
) UserService {
	return &userService{
		userRepo: userRepo,
		taskRepo: taskRepo,
	}
}

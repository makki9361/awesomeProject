package service

import (
	"awesomeProject/internal/models"
	"awesomeProject/internal/repository"
	"fmt"

	"github.com/go-playground/validator/v10"
)

type UserService struct {
	repo      *repository.UserRepository
	validator *validator.Validate
}

func NewUserService(repo *repository.UserRepository) *UserService {
	return &UserService{
		repo:      repo,
		validator: validator.New(),
	}
}

func (s *UserService) CreateUser(req *models.CreateUserRequest) (*models.User, error) {
	if err := s.validator.Struct(req); err != nil {
		return nil, fmt.Errorf("validation failed: %w", err)
	}

	user := &models.User{
		Name: req.Name,
		Role: req.Role,
	}

	if err := s.repo.Create(user); err != nil {
		return nil, err
	}

	return user, nil
}

func (s *UserService) GetUser(id int) (*models.User, error) {
	return s.repo.GetByID(id)
}

func (s *UserService) UpdateUser(id int, req *models.UpdateUserRequest) error {
	user, err := s.repo.GetByID(id)
	if err != nil {
		return err
	}
	if user == nil {
		return fmt.Errorf("user not found")
	}

	if req.Name != "" {
		user.Name = req.Name
	}
	if req.Role != "" {
		user.Role = req.Role
	}

	return s.repo.Update(id, user)
}

func (s *UserService) DeleteUser(id int) error {
	return s.repo.Delete(id)
}

func (s *UserService) ListUsers(limit, offset int) ([]models.User, error) {
	return s.repo.List(limit, offset)
}

package service

import (
	"awesomeProject/internal/models"
	"awesomeProject/internal/repository"
	"fmt"

	"github.com/go-playground/validator/v10"
)

type CategoryService struct {
	repo      *repository.CategoryRepository
	validator *validator.Validate
}

func NewCategoryService(repo *repository.CategoryRepository) *CategoryService {
	return &CategoryService{
		repo:      repo,
		validator: validator.New(),
	}
}

func (s *CategoryService) CreateCategory(req *models.CreateCategoryRequest) (*models.RuleCategory, error) {
	if err := s.validator.Struct(req); err != nil {
		return nil, fmt.Errorf("validation failed: %w", err)
	}

	category := &models.RuleCategory{
		Name: req.Name,
	}

	if err := s.repo.Create(category); err != nil {
		return nil, err
	}

	return category, nil
}

func (s *CategoryService) GetCategory(id int) (*models.RuleCategory, error) {
	return s.repo.GetByID(id)
}

func (s *CategoryService) UpdateCategory(id int, req *models.UpdateCategoryRequest) error {
	category, err := s.repo.GetByID(id)
	if err != nil {
		return err
	}
	if category == nil {
		return fmt.Errorf("category not found")
	}

	category.Name = req.Name
	return s.repo.Update(id, category)
}

func (s *CategoryService) DeleteCategory(id int) error {
	return s.repo.Delete(id)
}

func (s *CategoryService) ListCategories(limit, offset int) ([]models.RuleCategory, error) {
	return s.repo.List(limit, offset)
}

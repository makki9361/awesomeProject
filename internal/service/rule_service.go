package service

import (
	"awesomeProject/internal/models"
	"awesomeProject/internal/repository"
	"fmt"

	"github.com/go-playground/validator/v10"
)

type RuleService struct {
	repo      *repository.RuleRepository
	validator *validator.Validate
}

func NewRuleService(repo *repository.RuleRepository) *RuleService {
	return &RuleService{
		repo:      repo,
		validator: validator.New(),
	}
}

func (s *RuleService) CreateRule(req *models.CreateRuleRequest) (*models.Rule, error) {
	if err := s.validator.Struct(req); err != nil {
		return nil, fmt.Errorf("validation failed: %w", err)
	}

	rule := &models.Rule{
		Title:      req.Title,
		Content:    req.Content,
		CategoryID: req.CategoryID,
		Status:     req.Status,
		CreatedBy:  req.CreatedBy,
	}

	if err := s.repo.Create(rule); err != nil {
		return nil, err
	}

	return rule, nil
}

func (s *RuleService) GetRule(id int) (*models.RuleWithCategory, error) {
	return s.repo.GetByID(id)
}

func (s *RuleService) UpdateRule(id int, req *models.UpdateRuleRequest) error {
	existingRule, err := s.repo.GetByID(id)
	if err != nil {
		return err
	}
	if existingRule == nil {
		return fmt.Errorf("rule not found")
	}

	updatedRule := &models.Rule{
		ID:         id,
		Title:      existingRule.Title,
		Content:    existingRule.Content,
		CategoryID: existingRule.CategoryID,
		Status:     existingRule.Status,
		Version:    existingRule.Version,
	}

	versionChanged := false
	if req.Title != nil {
		updatedRule.Title = *req.Title
		versionChanged = true
	}
	if req.Content != nil {
		updatedRule.Content = *req.Content
	}
	if req.CategoryID != nil {
		updatedRule.CategoryID = *req.CategoryID
	}
	if req.Status != nil {
		updatedRule.Status = *req.Status
	}

	if versionChanged {
		updatedRule.Version = existingRule.Version + 1
	}

	return s.repo.Update(id, updatedRule)
}

func (s *RuleService) DeleteRule(id int) error {
	return s.repo.Delete(id)
}

func (s *RuleService) ListRules(categoryID *int, status *string, createdBy *int, search string, page, pageSize int) ([]models.RuleWithCategory, error) {
	if page < 1 {
		page = 1
	}
	if pageSize < 1 {
		pageSize = 10
	}
	if pageSize > 100 {
		pageSize = 100
	}

	offset := (page - 1) * pageSize

	filter := repository.RuleFilter{
		CategoryID: categoryID,
		Status:     status,
		CreatedBy:  createdBy,
		Search:     search,
		Limit:      pageSize,
		Offset:     offset,
	}

	rules, err := s.repo.List(filter)
	if err != nil {
		return nil, err
	}

	if rules == nil {
		return []models.RuleWithCategory{}, nil
	}

	return rules, nil
}

func (s *RuleService) PublishRule(id int, req *models.PublishRuleRequest) error {
	existingRule, err := s.repo.GetByID(id)
	if err != nil {
		return err
	}
	if existingRule == nil {
		return fmt.Errorf("rule not found")
	}

	updatedRule := &models.Rule{
		ID:         id,
		Title:      existingRule.Title,
		Content:    existingRule.Content,
		CategoryID: existingRule.CategoryID,
		Status:     req.Status,
		Version:    existingRule.Version,
	}

	return s.repo.Update(id, updatedRule)
}

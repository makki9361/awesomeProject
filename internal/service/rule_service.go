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
	fmt.Printf("UpdateRule called with id=%d, req=%+v\n", id, req)
	fmt.Printf("Status is nil? %v, value: %v\n", req.Status == nil, req.Status)
	fmt.Printf("CategoryID is nil? %v\n", req.CategoryID == nil)

	existingRule, err := s.repo.GetByID(id)
	if err != nil {
		return err
	}
	if existingRule == nil {
		return fmt.Errorf("rule not found")
	}

	fmt.Printf("Existing rule category_id: %d\n", existingRule.CategoryID)

	updatedRule := &models.Rule{
		ID:         id,
		Title:      existingRule.Title,
		Content:    existingRule.Content,
		CategoryID: existingRule.CategoryID,
		Status:     existingRule.Status,
	}

	if req.Title != nil {
		updatedRule.Title = *req.Title
	}
	if req.Content != nil {
		updatedRule.Content = *req.Content
	}
	if req.CategoryID != nil {
		fmt.Printf("Updating category_id from %d to %d\n", existingRule.CategoryID, *req.CategoryID)
		updatedRule.CategoryID = *req.CategoryID
	}
	if req.Status != nil {
		fmt.Printf("Updating status from %s to %s\n", existingRule.Status, *req.Status)
		updatedRule.Status = *req.Status
	}

	fmt.Printf("Final updatedRule: %+v\n", updatedRule)

	return s.repo.Update(id, updatedRule)
}

func (s *RuleService) DeleteRule(id int) error {
	return s.repo.Delete(id)
}

func (s *RuleService) ListRules(categoryID *int, status *string, createdBy *int, page, pageSize int) ([]models.RuleWithCategory, error) {
	if page < 1 {
		page = 1
	}
	if pageSize < 1 {
		pageSize = 10
	}

	offset := (page - 1) * pageSize

	filter := repository.RuleFilter{
		CategoryID: categoryID,
		Status:     status,
		CreatedBy:  createdBy,
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

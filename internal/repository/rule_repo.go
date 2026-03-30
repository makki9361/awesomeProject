package repository

import (
	"awesomeProject/internal/models"
	"database/sql"
	"fmt"
	"strings"

	"github.com/jmoiron/sqlx"
)

type RuleRepository struct {
	db *sqlx.DB
}

func NewRuleRepository(db *sqlx.DB) *RuleRepository {
	return &RuleRepository{db: db}
}

func (r *RuleRepository) Create(rule *models.Rule) error {
	query := `INSERT INTO rules (title, content, category_id, status, version, created_by, created_at, updated_at) 
              VALUES ($1, $2, $3, $4, 1, $5, NOW(), NOW()) RETURNING id`

	err := r.db.QueryRowx(query, rule.Title, rule.Content, rule.CategoryID, rule.Status, rule.CreatedBy).Scan(&rule.ID)
	if err != nil {
		return fmt.Errorf("failed to create rule: %w", err)
	}
	return nil
}

func (r *RuleRepository) GetByID(id int) (*models.RuleWithCategory, error) {
	var rule models.RuleWithCategory
	query := `SELECT r.*, rc.name as category_name 
              FROM rules r 
              JOIN rule_categories rc ON r.category_id = rc.id 
              WHERE r.id = $1`

	err := r.db.Get(&rule, query, id)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to get rule: %w", err)
	}
	return &rule, nil
}

func (r *RuleRepository) Update(id int, rule *models.Rule) error {
	updates := []string{}
	args := []interface{}{}
	argIndex := 1

	if rule.Title != "" {
		updates = append(updates, fmt.Sprintf("title = $%d", argIndex))
		args = append(args, rule.Title)
		argIndex++
	}

	if rule.Content != "" {
		updates = append(updates, fmt.Sprintf("content = $%d", argIndex))
		args = append(args, rule.Content)
		argIndex++
	}

	if rule.CategoryID != 0 {
		updates = append(updates, fmt.Sprintf("category_id = $%d", argIndex))
		args = append(args, rule.CategoryID)
		argIndex++
	}

	if rule.Status != "" {
		updates = append(updates, fmt.Sprintf("status = $%d", argIndex))
		args = append(args, rule.Status)
		argIndex++
	}

	if rule.Version != 0 {
		updates = append(updates, fmt.Sprintf("version = $%d", argIndex))
		args = append(args, rule.Version)
		argIndex++
	}

	if len(updates) == 0 {
		return nil
	}

	updates = append(updates, "updated_at = NOW()")

	query := fmt.Sprintf("UPDATE rules SET %s WHERE id = $%d",
		strings.Join(updates, ", "),
		argIndex)
	args = append(args, id)

	_, err := r.db.Exec(query, args...)
	if err != nil {
		return fmt.Errorf("failed to update rule: %w", err)
	}

	return nil
}
func (r *RuleRepository) Delete(id int) error {
	query := `DELETE FROM rules WHERE id = $1`

	result, err := r.db.Exec(query, id)
	if err != nil {
		return fmt.Errorf("failed to delete rule: %w", err)
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rows == 0 {
		return sql.ErrNoRows
	}

	return nil
}

type RuleFilter struct {
	CategoryID *int
	Status     *string
	CreatedBy  *int
	Search     string
	Limit      int
	Offset     int
}

func (r *RuleRepository) List(filter RuleFilter) ([]models.RuleWithCategory, error) {
	var rules []models.RuleWithCategory

	query := `SELECT r.*, rc.name as category_name 
              FROM rules r 
              JOIN rule_categories rc ON r.category_id = rc.id 
              WHERE 1=1`

	args := []interface{}{}
	argIndex := 1

	if filter.CategoryID != nil {
		query += fmt.Sprintf(" AND r.category_id = $%d", argIndex)
		args = append(args, *filter.CategoryID)
		argIndex++
	}

	if filter.Status != nil && *filter.Status != "" {
		query += fmt.Sprintf(" AND r.status = $%d", argIndex)
		args = append(args, *filter.Status)
		argIndex++
	}

	if filter.CreatedBy != nil {
		query += fmt.Sprintf(" AND r.created_by = $%d", argIndex)
		args = append(args, *filter.CreatedBy)
		argIndex++
	}

	if filter.Search != "" {
		query += fmt.Sprintf(" AND (r.title ILIKE $%d OR r.content ILIKE $%d)", argIndex, argIndex)
		args = append(args, "%"+filter.Search+"%")
		argIndex++
	}

	query += fmt.Sprintf(" ORDER BY r.created_at DESC LIMIT $%d OFFSET $%d", argIndex, argIndex+1)
	args = append(args, filter.Limit, filter.Offset)

	err := r.db.Select(&rules, query, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to list rules: %w", err)
	}

	return rules, nil
}

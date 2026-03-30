package repository

import (
	"awesomeProject/internal/models"
	"database/sql"
	"fmt"

	"github.com/jmoiron/sqlx"
)

type CategoryRepository struct {
	db *sqlx.DB
}

func NewCategoryRepository(db *sqlx.DB) *CategoryRepository {
	return &CategoryRepository{db: db}
}

func (r *CategoryRepository) Create(category *models.RuleCategory) error {
	query := `INSERT INTO rule_categories (name, created_at, updated_at) 
              VALUES ($1, NOW(), NOW()) RETURNING id`

	err := r.db.QueryRowx(query, category.Name).Scan(&category.ID)
	if err != nil {
		return fmt.Errorf("failed to create category: %w", err)
	}
	return nil
}

func (r *CategoryRepository) GetByID(id int) (*models.RuleCategory, error) {
	var category models.RuleCategory
	query := `SELECT id, name, created_at, updated_at FROM rule_categories WHERE id = $1`

	err := r.db.Get(&category, query, id)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to get category: %w", err)
	}
	return &category, nil
}

func (r *CategoryRepository) Update(id int, category *models.RuleCategory) error {
	query := `UPDATE rule_categories SET name = $1, updated_at = NOW() WHERE id = $2`

	result, err := r.db.Exec(query, category.Name, id)
	if err != nil {
		return fmt.Errorf("failed to update category: %w", err)
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

func (r *CategoryRepository) Delete(id int) error {
	query := `DELETE FROM rule_categories WHERE id = $1`

	result, err := r.db.Exec(query, id)
	if err != nil {
		return fmt.Errorf("failed to delete category: %w", err)
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

func (r *CategoryRepository) List(limit, offset int) ([]models.RuleCategory, error) {
	var categories []models.RuleCategory
	query := `SELECT id, name, created_at, updated_at FROM rule_categories 
              ORDER BY id LIMIT $1 OFFSET $2`

	err := r.db.Select(&categories, query, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to list categories: %w", err)
	}

	return categories, nil
}

func (r *CategoryRepository) HasRules(categoryID int) (bool, error) {
	var count int
	query := `SELECT COUNT(*) FROM rules WHERE category_id = $1`

	err := r.db.Get(&count, query, categoryID)
	if err != nil {
		return false, fmt.Errorf("failed to check rules: %w", err)
	}

	return count > 0, nil
}

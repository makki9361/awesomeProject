package models

import (
	"time"
)

type Rule struct {
	ID         int       `db:"id" json:"id"`
	Title      string    `db:"title" json:"title" validate:"required,min=1,max=200"`
	Content    string    `db:"content" json:"content" validate:"required"`
	CategoryID int       `db:"category_id" json:"category_id" validate:"required"`
	Status     string    `db:"status" json:"status" validate:"required,oneof=draft published archived"`
	Version    int       `db:"version" json:"version"`
	CreatedBy  int       `db:"created_by" json:"created_by" validate:"required"`
	CreatedAt  time.Time `db:"created_at" json:"created_at"`
	UpdatedAt  time.Time `db:"updated_at" json:"updated_at"`
}

type CreateRuleRequest struct {
	Title      string `json:"title" validate:"required,min=1,max=200"`
	Content    string `json:"content" validate:"required"`
	CategoryID int    `json:"category_id" validate:"required"`
	Status     string `json:"status" validate:"required,oneof=draft published archived"`
	CreatedBy  int    `json:"created_by" validate:"required"`
}

type UpdateRuleRequest struct {
	Title      *string `json:"title" validate:"omitempty,min=1,max=200"`
	Content    *string `json:"content" validate:"omitempty"`
	CategoryID *int    `json:"category_id" validate:"omitempty"`
	Status     *string `json:"status" validate:"omitempty,oneof=draft published archived"`
}

type RuleWithCategory struct {
	Rule
	CategoryName string `db:"category_name" json:"category_name"`
}

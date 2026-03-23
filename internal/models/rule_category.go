package models

import (
	"time"
)

type RuleCategory struct {
	ID        int       `db:"id" json:"id"`
	Name      string    `db:"name" json:"name" validate:"required,min=1,max=100"`
	CreatedAt time.Time `db:"created_at" json:"created_at"`
	UpdatedAt time.Time `db:"updated_at" json:"updated_at"`
}

type CreateCategoryRequest struct {
	Name string `json:"name" validate:"required,min=1,max=100"`
}

type UpdateCategoryRequest struct {
	Name string `json:"name" validate:"required,min=1,max=100"`
}

package models

import (
	"time"
)

type User struct {
	ID        int       `db:"id" json:"id"`
	Name      string    `db:"name" json:"name" validate:"required,min=2,max=100"`
	Role      string    `db:"role" json:"role" validate:"required,oneof=admin employee"`
	CreatedAt time.Time `db:"created_at" json:"created_at"`
	UpdatedAt time.Time `db:"updated_at" json:"updated_at"`
}

type CreateUserRequest struct {
	Name string `json:"name" validate:"required,min=2,max=100"`
	Role string `json:"role" validate:"required,oneof=admin employee"`
}

type UpdateUserRequest struct {
	Name string `json:"name" validate:"omitempty,min=2,max=100"`
	Role string `json:"role" validate:"omitempty,oneof=admin employee"`
}

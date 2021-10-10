package domain

import (
	"context"
	"time"
)

type User struct {
	ID        int64     `json:"ID"`
	Fullname  string    `json:"fullname"`
	Username  string    `json:"username"`
	Email     string    `json:"email" validate:"required,email"`
	Password  string    `json:"passwd"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type UserUsecase interface {
	Fetch(ctx context.Context, cursor string, num int64) ([]User, string, error)
}

type UserRepository interface {
	Fetch(ctx context.Context, cursor string, num int64) (res []User, nextCursor int64, err error)
}

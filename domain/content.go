package domain

import (
	"context"
	"time"
)

// Content ...
type Content struct {
	ID        int64     `json:"id"`
	Title     string    `json:"title" validate:"required"`
	Content   string    `json:"content" validate:"required"`
	Author    Author    `json:"author"`
	UpdatedAt time.Time `json:"updated_at"`
	CreatedAt time.Time `json:"created_at"`
}

// ContentUsecase represent the content's usecases
type ContentUsecase interface {
	Fetch(ctx context.Context, cursor string, num int64) ([]Content, string, error)
	GetByID(ctx context.Context, id int64) (Content, error)
	Update(ctx context.Context, ar *Content) error
	GetByTitle(ctx context.Context, title string) (Content, error)
	Store(context.Context, *Content) error
	Delete(ctx context.Context, id int64) error
}

// ContentRepository represent the content's repository contract
type ContentRepository interface {
	Fetch(ctx context.Context, cursor string, num int64) (res []Content, nextCursor string, err error)
	GetByID(ctx context.Context, id int64) (Content, error)
	GetByTitle(ctx context.Context, title string) (Content, error)
	Update(ctx context.Context, ar *Content) error
	Store(ctx context.Context, a *Content) error
	Delete(ctx context.Context, id int64) error
}

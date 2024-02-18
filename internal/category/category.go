package category

import (
	"context"
	"time"
)

type Service interface {
	// GetCategoryByID returns a category with the given category ID.
	GetCategoryByID(ctx context.Context, id int64) (Category, error)

	// GetCategories returns list of categories.
	GetCategories(ctx context.Context) ([]Category, error)
}

type Category struct {
	ID          int64
	Name        string
	Description string
	CreateTime  time.Time
	UpdateTime  time.Time
}

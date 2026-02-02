package repository

import (
	"context"

	"kasir-api/internal/category/model"
)

// Repository defines the interface for category data access
type Repository interface {
	List(ctx context.Context) ([]model.Category, error)
	Create(ctx context.Context, c model.Category) (model.Category, error)
	GetByID(ctx context.Context, id int64) (model.Category, error)
	Update(ctx context.Context, id int64, c model.Category) (model.Category, error)
	Delete(ctx context.Context, id int64) error
}

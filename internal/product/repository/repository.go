package repository

import (
	"context"

	"kasir-api/internal/product/model"
)

// Repository defines the interface for product data access
type Repository interface {
	List(ctx context.Context) ([]model.Product, error)
	Search(ctx context.Context, name string) ([]model.Product, error)
	Create(ctx context.Context, p model.Product) (model.Product, error)
	GetByID(ctx context.Context, id int64) (model.Product, error)
	Update(ctx context.Context, id int64, p model.Product) (model.Product, error)
	Delete(ctx context.Context, id int64) error
}

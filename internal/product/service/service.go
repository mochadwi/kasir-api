package service

import (
	"context"
	"fmt"

	"kasir-api/internal/product/model"
	"kasir-api/internal/product/repository"
)

// Service defines the interface for product business logic
type Service interface {
	ListProducts(ctx context.Context) ([]model.Product, error)
	SearchProducts(ctx context.Context, name string) ([]model.Product, error)
	CreateProduct(ctx context.Context, p model.Product) (model.Product, error)
	GetProduct(ctx context.Context, id int64) (model.Product, error)
	UpdateProduct(ctx context.Context, id int64, p model.Product) (model.Product, error)
	DeleteProduct(ctx context.Context, id int64) error
}

// productService implements Service interface
type productService struct {
	repo repository.Repository
}

// New creates a new product service
func New(repo repository.Repository) Service {
	return &productService{repo: repo}
}

// ListProducts returns all products
func (s *productService) ListProducts(ctx context.Context) ([]model.Product, error) {
	return s.repo.List(ctx)
}

// SearchProducts returns products matching the name
func (s *productService) SearchProducts(ctx context.Context, name string) ([]model.Product, error) {
	return s.repo.Search(ctx, name)
}

// CreateProduct creates a new product with validation
func (s *productService) CreateProduct(ctx context.Context, p model.Product) (model.Product, error) {
	if err := validateProduct(p); err != nil {
		return model.Product{}, err
	}
	return s.repo.Create(ctx, p)
}

// GetProduct returns a product by ID
func (s *productService) GetProduct(ctx context.Context, id int64) (model.Product, error) {
	return s.repo.GetByID(ctx, id)
}

// UpdateProduct updates an existing product with validation
func (s *productService) UpdateProduct(ctx context.Context, id int64, p model.Product) (model.Product, error) {
	if err := validateProduct(p); err != nil {
		return model.Product{}, err
	}
	return s.repo.Update(ctx, id, p)
}

// DeleteProduct removes a product by ID
func (s *productService) DeleteProduct(ctx context.Context, id int64) error {
	return s.repo.Delete(ctx, id)
}

func validateProduct(p model.Product) error {
	if p.Nama == "" {
		return fmt.Errorf("%w: nama is required", model.ErrInvalidProduct)
	}
	if p.Harga < 0 {
		return fmt.Errorf("%w: harga cannot be negative", model.ErrInvalidProduct)
	}
	if p.Stok < 0 {
		return fmt.Errorf("%w: stok cannot be negative", model.ErrInvalidProduct)
	}
	if p.CategoryID != nil && *p.CategoryID <= 0 {
		return fmt.Errorf("%w: category_id must be positive", model.ErrInvalidProduct)
	}
	return nil
}

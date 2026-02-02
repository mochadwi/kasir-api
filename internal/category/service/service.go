package service

import (
	"context"

	"kasir-api/internal/category/model"
	"kasir-api/internal/category/repository"
)

// Service defines the interface for category business logic
type Service interface {
	ListCategories(ctx context.Context) ([]model.Category, error)
	CreateCategory(ctx context.Context, c model.Category) (model.Category, error)
	GetCategory(ctx context.Context, id int64) (model.Category, error)
	UpdateCategory(ctx context.Context, id int64, c model.Category) (model.Category, error)
	DeleteCategory(ctx context.Context, id int64) error
}

// categoryService implements Service interface
type categoryService struct {
	repo repository.Repository
}

// New creates a new category service
func New(repo repository.Repository) Service {
	return &categoryService{repo: repo}
}

// ListCategories returns all categories
func (s *categoryService) ListCategories(ctx context.Context) ([]model.Category, error) {
	return s.repo.List(ctx)
}

// CreateCategory creates a new category with validation
func (s *categoryService) CreateCategory(ctx context.Context, c model.Category) (model.Category, error) {
	if c.Name == "" {
		return model.Category{}, model.ErrInvalidCategory
	}
	return s.repo.Create(ctx, c)
}

// GetCategory returns a category by ID
func (s *categoryService) GetCategory(ctx context.Context, id int64) (model.Category, error) {
	return s.repo.GetByID(ctx, id)
}

// UpdateCategory updates an existing category with validation
func (s *categoryService) UpdateCategory(ctx context.Context, id int64, c model.Category) (model.Category, error) {
	if c.Name == "" {
		return model.Category{}, model.ErrInvalidCategory
	}
	return s.repo.Update(ctx, id, c)
}

// DeleteCategory removes a category by ID
func (s *categoryService) DeleteCategory(ctx context.Context, id int64) error {
	return s.repo.Delete(ctx, id)
}

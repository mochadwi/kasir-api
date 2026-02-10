package service

import (
	"context"
	"fmt"

	"kasir-api/internal/transaction/model"
	"kasir-api/internal/transaction/repository"
)

// Service defines the interface for transaction business logic
type Service interface {
	CreateTransaction(ctx context.Context, items []model.CheckoutItem) (*model.Transaction, error)
}

// transactionService implements Service interface
type transactionService struct {
	repo repository.Repository
}

// New creates a new transaction service
func New(repo repository.Repository) Service {
	return &transactionService{repo: repo}
}

// CreateTransaction creates a new transaction with validation
func (s *transactionService) CreateTransaction(ctx context.Context, items []model.CheckoutItem) (*model.Transaction, error) {
	if err := validateCheckoutItems(items); err != nil {
		return nil, err
	}
	return s.repo.CreateTransaction(ctx, items)
}

func validateCheckoutItems(items []model.CheckoutItem) error {
	if len(items) == 0 {
		return fmt.Errorf("%w: at least one item required", model.ErrInvalidTransaction)
	}
	for i, item := range items {
		if item.ProductID <= 0 {
			return fmt.Errorf("%w: item %d has invalid product_id", model.ErrInvalidTransaction, i)
		}
		if item.Quantity <= 0 {
			return fmt.Errorf("%w: item %d has invalid quantity", model.ErrInvalidTransaction, i)
		}
	}
	return nil
}

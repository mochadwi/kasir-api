package repository

import (
	"context"

	"kasir-api/internal/transaction/model"
)

// Repository defines the interface for transaction data access
type Repository interface {
	CreateTransaction(ctx context.Context, items []model.CheckoutItem) (*model.Transaction, error)
}

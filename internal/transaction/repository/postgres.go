package repository

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"kasir-api/internal/transaction/model"
)

// PostgresRepository implements Repository interface using PostgreSQL
type PostgresRepository struct {
	db *sql.DB
}

// NewPostgres creates a new PostgreSQL repository
func NewPostgres(db *sql.DB) Repository {
	return &PostgresRepository{db: db}
}

// CreateTransaction creates a new transaction with details and updates stock
func (r *PostgresRepository) CreateTransaction(ctx context.Context, items []model.CheckoutItem) (*model.Transaction, error) {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback()

	var totalAmount int
	var details []model.TransactionDetail

	for _, item := range items {
		var nama string
		var harga, stok int
		err := tx.QueryRowContext(ctx,
			"SELECT nama, harga, stok FROM products WHERE id = $1 FOR UPDATE",
			item.ProductID).Scan(&nama, &harga, &stok)
		if err != nil {
			if err == sql.ErrNoRows {
				return nil, fmt.Errorf("%w: product %d", model.ErrProductNotFound, item.ProductID)
			}
			return nil, fmt.Errorf("failed to get product %d: %w", item.ProductID, err)
		}

		if stok < item.Quantity {
			return nil, fmt.Errorf("%w: product %s has %d stock, requested %d",
				model.ErrInsufficientStock, nama, stok, item.Quantity)
		}

		_, err = tx.ExecContext(ctx,
			"UPDATE products SET stok = stok - $1 WHERE id = $2",
			item.Quantity, item.ProductID)
		if err != nil {
			return nil, fmt.Errorf("failed to update stock for product %d: %w", item.ProductID, err)
		}

		subtotal := harga * item.Quantity
		totalAmount += subtotal

		details = append(details, model.TransactionDetail{
			ProductID:   item.ProductID,
			ProductName: nama,
			Quantity:    item.Quantity,
			Subtotal:    subtotal,
		})
	}

	var transactionID int64
	err = tx.QueryRowContext(ctx,
		"INSERT INTO transactions (total_amount) VALUES ($1) RETURNING id",
		totalAmount).Scan(&transactionID)
	if err != nil {
		return nil, fmt.Errorf("failed to create transaction: %w", err)
	}

	for i := range details {
		details[i].TransactionID = transactionID
		err = tx.QueryRowContext(ctx,
			"INSERT INTO transaction_details (transaction_id, product_id, quantity, subtotal) VALUES ($1, $2, $3, $4) RETURNING id",
			transactionID, details[i].ProductID, details[i].Quantity, details[i].Subtotal).Scan(&details[i].ID)
		if err != nil {
			return nil, fmt.Errorf("failed to create transaction detail: %w", err)
		}
	}

	if err := tx.Commit(); err != nil {
		return nil, fmt.Errorf("failed to commit transaction: %w", err)
	}

	return &model.Transaction{
		ID:          transactionID,
		TotalAmount: totalAmount,
		CreatedAt:   time.Now(),
		Details:     details,
	}, nil
}

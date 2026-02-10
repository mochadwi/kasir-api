package model

import (
	"errors"
	"time"
)

// Transaction represents a sales transaction
type Transaction struct {
	ID          int64               `json:"id"`
	TotalAmount int                 `json:"total_amount"`
	CreatedAt   time.Time           `json:"created_at"`
	Details     []TransactionDetail `json:"details,omitempty"`
}

// TransactionDetail represents a single item in a transaction
type TransactionDetail struct {
	ID            int64  `json:"id"`
	TransactionID int64  `json:"transaction_id"`
	ProductID     int64  `json:"product_id"`
	ProductName   string `json:"product_name"`
	Quantity      int    `json:"quantity"`
	Subtotal      int    `json:"subtotal"`
}

// CheckoutItem represents a single item in a checkout request
type CheckoutItem struct {
	ProductID int64 `json:"product_id"`
	Quantity  int   `json:"quantity"`
}

// CheckoutRequest represents the request body for checkout
type CheckoutRequest struct {
	Items []CheckoutItem `json:"items"`
}

// Domain errors for transaction operations
var (
	ErrTransactionNotFound = errors.New("transaction not found")
	ErrInvalidTransaction  = errors.New("invalid transaction data")
	ErrInsufficientStock   = errors.New("insufficient stock")
	ErrProductNotFound     = errors.New("product not found")
)

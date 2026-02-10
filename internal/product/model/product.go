package model

import "errors"

// Product represents a product in the system
type Product struct {
	ID         int64  `json:"id"`
	Nama       string `json:"nama"`
	Harga      int    `json:"harga"`
	Stok       int    `json:"stok"`
	CategoryID *int64 `json:"category_id,omitempty"`
}

// Domain errors for product operations
var (
	ErrProductNotFound = errors.New("product not found")
	ErrInvalidProduct  = errors.New("invalid product data")
)

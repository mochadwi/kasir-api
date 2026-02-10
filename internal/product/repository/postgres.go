package repository

import (
	"context"
	"database/sql"
	"fmt"

	"kasir-api/internal/product/model"
)

// PostgresRepository implements Repository interface using PostgreSQL
type PostgresRepository struct {
	db *sql.DB
}

// NewPostgres creates a new PostgreSQL repository
func NewPostgres(db *sql.DB) Repository {
	return &PostgresRepository{db: db}
}

// List returns all products
func (r *PostgresRepository) List(ctx context.Context) ([]model.Product, error) {
	rows, err := r.db.QueryContext(ctx, "SELECT id, nama, harga, stok, category_id FROM products ORDER BY id")
	if err != nil {
		return nil, fmt.Errorf("failed to list products: %w", err)
	}
	defer rows.Close()

	var products []model.Product
	for rows.Next() {
		var p model.Product
		if err := rows.Scan(&p.ID, &p.Nama, &p.Harga, &p.Stok, &p.CategoryID); err != nil {
			return nil, fmt.Errorf("failed to scan product: %w", err)
		}
		products = append(products, p)
	}

	return products, rows.Err()
}

// Create inserts a new product and returns it with ID
func (r *PostgresRepository) Create(ctx context.Context, p model.Product) (model.Product, error) {
	var id int64
	err := r.db.QueryRowContext(ctx,
		"INSERT INTO products (nama, harga, stok, category_id) VALUES ($1, $2, $3, $4) RETURNING id",
		p.Nama, p.Harga, p.Stok, p.CategoryID).Scan(&id)
	if err != nil {
		return model.Product{}, fmt.Errorf("failed to create product: %w", err)
	}

	p.ID = id
	return p, nil
}

// GetByID returns a product by ID
func (r *PostgresRepository) GetByID(ctx context.Context, id int64) (model.Product, error) {
	var p model.Product
	err := r.db.QueryRowContext(ctx,
		"SELECT id, nama, harga, stok, category_id FROM products WHERE id = $1",
		id).Scan(&p.ID, &p.Nama, &p.Harga, &p.Stok, &p.CategoryID)
	if err != nil {
		if err == sql.ErrNoRows {
			return model.Product{}, model.ErrProductNotFound
		}
		return model.Product{}, fmt.Errorf("failed to get product: %w", err)
	}
	return p, nil
}

// Update modifies an existing product
func (r *PostgresRepository) Update(ctx context.Context, id int64, p model.Product) (model.Product, error) {
	result, err := r.db.ExecContext(ctx,
		"UPDATE products SET nama = $1, harga = $2, stok = $3, category_id = $4 WHERE id = $5",
		p.Nama, p.Harga, p.Stok, p.CategoryID, id)
	if err != nil {
		return model.Product{}, fmt.Errorf("failed to update product: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return model.Product{}, fmt.Errorf("failed to check rows affected: %w", err)
	}
	if rowsAffected == 0 {
		return model.Product{}, model.ErrProductNotFound
	}

	p.ID = id
	return p, nil
}

// Delete removes a product by ID
func (r *PostgresRepository) Delete(ctx context.Context, id int64) error {
	result, err := r.db.ExecContext(ctx, "DELETE FROM products WHERE id = $1", id)
	if err != nil {
		return fmt.Errorf("failed to delete product: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to check rows affected: %w", err)
	}
	if rowsAffected == 0 {
		return model.ErrProductNotFound
	}

	return nil
}

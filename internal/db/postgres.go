package db

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	_ "github.com/jackc/pgx/v5/stdlib"
)

// Open creates a new PostgreSQL database connection
func Open(url string) (*sql.DB, error) {
	db, err := sql.Open("pgx", url)
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	db.SetMaxOpenConns(10)
	db.SetMaxIdleConns(5)
	db.SetConnMaxLifetime(5 * time.Minute)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := db.PingContext(ctx); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	return db, nil
}

// Migrate creates database tables if they don't exist
func Migrate(db *sql.DB) error {
	categoriesTable := `
	CREATE TABLE IF NOT EXISTS categories (
		id SERIAL PRIMARY KEY,
		name VARCHAR(100) NOT NULL,
		description TEXT,
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
	);`

	productsTable := `
	CREATE TABLE IF NOT EXISTS products (
		id SERIAL PRIMARY KEY,
		nama VARCHAR(200) NOT NULL,
		harga INTEGER NOT NULL DEFAULT 0,
		stok INTEGER NOT NULL DEFAULT 0,
		category_id INTEGER REFERENCES categories(id),
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
	);`

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if _, err := db.ExecContext(ctx, categoriesTable); err != nil {
		return fmt.Errorf("failed to create categories table: %w", err)
	}

	if _, err := db.ExecContext(ctx, productsTable); err != nil {
		return fmt.Errorf("failed to create products table: %w", err)
	}

	// Add category_id column if it doesn't exist (migration for existing tables)
	addCategoryID := `
	ALTER TABLE products 
	ADD COLUMN IF NOT EXISTS category_id INTEGER REFERENCES categories(id);`

	if _, err := db.ExecContext(ctx, addCategoryID); err != nil {
		return fmt.Errorf("failed to add category_id column: %w", err)
	}

	transactionsTable := `
	CREATE TABLE IF NOT EXISTS transactions (
		id SERIAL PRIMARY KEY,
		total_amount INTEGER NOT NULL,
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
	);`

	if _, err := db.ExecContext(ctx, transactionsTable); err != nil {
		return fmt.Errorf("failed to create transactions table: %w", err)
	}

	transactionDetailsTable := `
	CREATE TABLE IF NOT EXISTS transaction_details (
		id SERIAL PRIMARY KEY,
		transaction_id INTEGER REFERENCES transactions(id) ON DELETE CASCADE,
		product_id INTEGER REFERENCES products(id),
		quantity INTEGER NOT NULL,
		subtotal INTEGER NOT NULL
	);`

	if _, err := db.ExecContext(ctx, transactionDetailsTable); err != nil {
		return fmt.Errorf("failed to create transaction_details table: %w", err)
	}

	// Insert sample data
	insertCategories := `
	INSERT INTO categories (name, description) VALUES
		('Minuman', 'Minuman dingin dan panas'),
		('Makanan', 'Makanan ringan dan berat'),
		('Snack', 'Camilan dan kudapan')
	ON CONFLICT DO NOTHING;`

	insertProducts := `
	INSERT INTO products (nama, harga, stok, category_id) VALUES
		('Es Teh Manis', 5000, 100, 1),
		('Nasi Goreng', 25000, 50, 2),
		('Keripik Kentang', 8000, 200, 3)
	ON CONFLICT DO NOTHING;`

	db.ExecContext(ctx, insertCategories)
	db.ExecContext(ctx, insertProducts)

	return nil
}

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

	// Insert sample data
	insertCategories := `
	INSERT INTO categories (name, description) VALUES
		('Minuman', 'Minuman dingin dan panas'),
		('Makanan', 'Makanan ringan dan berat'),
		('Snack', 'Camilan dan kudapan')
	ON CONFLICT DO NOTHING;`

	insertProducts := `
	INSERT INTO products (nama, harga, stok) VALUES
		('Es Teh Manis', 5000, 100),
		('Nasi Goreng', 25000, 50),
		('Keripik Kentang', 8000, 200)
	ON CONFLICT DO NOTHING;`

	db.ExecContext(ctx, insertCategories)
	db.ExecContext(ctx, insertProducts)

	return nil
}

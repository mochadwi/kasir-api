package repository

import (
	"context"
	"database/sql"
	"fmt"

	"kasir-api/internal/category/model"
)

// PostgresRepository implements Repository interface using PostgreSQL
type PostgresRepository struct {
	db *sql.DB
}

// NewPostgres creates a new PostgreSQL repository
func NewPostgres(db *sql.DB) Repository {
	return &PostgresRepository{db: db}
}

// List returns all categories
func (r *PostgresRepository) List(ctx context.Context) ([]model.Category, error) {
	rows, err := r.db.QueryContext(ctx, "SELECT id, name, description FROM categories ORDER BY id")
	if err != nil {
		return nil, fmt.Errorf("failed to list categories: %w", err)
	}
	defer rows.Close()

	var categories []model.Category
	for rows.Next() {
		var c model.Category
		if err := rows.Scan(&c.ID, &c.Name, &c.Description); err != nil {
			return nil, fmt.Errorf("failed to scan category: %w", err)
		}
		categories = append(categories, c)
	}

	return categories, rows.Err()
}

// Create inserts a new category and returns it with ID
func (r *PostgresRepository) Create(ctx context.Context, c model.Category) (model.Category, error) {
	var id int64
	err := r.db.QueryRowContext(ctx,
		"INSERT INTO categories (name, description) VALUES ($1, $2) RETURNING id",
		c.Name, c.Description).Scan(&id)
	if err != nil {
		return model.Category{}, fmt.Errorf("failed to create category: %w", err)
	}

	c.ID = id
	return c, nil
}

// GetByID returns a category by ID
func (r *PostgresRepository) GetByID(ctx context.Context, id int64) (model.Category, error) {
	var c model.Category
	err := r.db.QueryRowContext(ctx,
		"SELECT id, name, description FROM categories WHERE id = $1",
		id).Scan(&c.ID, &c.Name, &c.Description)
	if err != nil {
		if err == sql.ErrNoRows {
			return model.Category{}, model.ErrCategoryNotFound
		}
		return model.Category{}, fmt.Errorf("failed to get category: %w", err)
	}
	return c, nil
}

// Update modifies an existing category
func (r *PostgresRepository) Update(ctx context.Context, id int64, c model.Category) (model.Category, error) {
	result, err := r.db.ExecContext(ctx,
		"UPDATE categories SET name = $1, description = $2 WHERE id = $3",
		c.Name, c.Description, id)
	if err != nil {
		return model.Category{}, fmt.Errorf("failed to update category: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return model.Category{}, fmt.Errorf("failed to check rows affected: %w", err)
	}
	if rowsAffected == 0 {
		return model.Category{}, model.ErrCategoryNotFound
	}

	c.ID = id
	return c, nil
}

// Delete removes a category by ID
func (r *PostgresRepository) Delete(ctx context.Context, id int64) error {
	result, err := r.db.ExecContext(ctx, "DELETE FROM categories WHERE id = $1", id)
	if err != nil {
		return fmt.Errorf("failed to delete category: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to check rows affected: %w", err)
	}
	if rowsAffected == 0 {
		return model.ErrCategoryNotFound
	}

	return nil
}

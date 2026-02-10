package repository

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"kasir-api/internal/report/model"
)

type PostgresRepository struct {
	db *sql.DB
}

func NewPostgres(db *sql.DB) Repository {
	return &PostgresRepository{db: db}
}

func (r *PostgresRepository) GetSalesReport(ctx context.Context, startDate, endDate *time.Time) (*model.SalesReport, error) {
	report := &model.SalesReport{}

	err := r.db.QueryRowContext(ctx,
		"SELECT COALESCE(SUM(total_amount), 0) FROM transactions WHERE created_at BETWEEN $1 AND $2",
		startDate, endDate).Scan(&report.TotalRevenue)
	if err != nil {
		return nil, fmt.Errorf("failed to get total revenue: %w", err)
	}

	err = r.db.QueryRowContext(ctx,
		"SELECT COUNT(*) FROM transactions WHERE created_at BETWEEN $1 AND $2",
		startDate, endDate).Scan(&report.TotalTransactions)
	if err != nil {
		return nil, fmt.Errorf("failed to get total transactions: %w", err)
	}

	var produkTerlaris sql.NullString
	var qtyTerjual sql.NullInt64
	err = r.db.QueryRowContext(ctx,
		`SELECT p.nama, SUM(td.quantity) as qty 
		 FROM transaction_details td 
		 JOIN products p ON td.product_id = p.id 
		 JOIN transactions t ON td.transaction_id = t.id 
		 WHERE t.created_at BETWEEN $1 AND $2 
		 GROUP BY p.nama 
		 ORDER BY qty DESC 
		 LIMIT 1`,
		startDate, endDate).Scan(&produkTerlaris, &qtyTerjual)
	if err != nil && err != sql.ErrNoRows {
		return nil, fmt.Errorf("failed to get produk terlaris: %w", err)
	}

	if produkTerlaris.Valid {
		report.ProdukTerlaris = produkTerlaris.String
		report.QtyTerjual = int(qtyTerjual.Int64)
	}

	return report, nil
}

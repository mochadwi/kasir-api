package repository

import (
	"context"
	"time"

	"kasir-api/internal/report/model"
)

type Repository interface {
	GetSalesReport(ctx context.Context, startDate, endDate *time.Time) (*model.SalesReport, error)
}

package service

import (
	"context"
	"time"

	"kasir-api/internal/report/model"
	"kasir-api/internal/report/repository"
)

type Service interface {
	GetSalesReport(ctx context.Context, startDate, endDate *string) (*model.SalesReport, error)
	GetTodayReport(ctx context.Context) (*model.SalesReport, error)
}

type reportService struct {
	repo repository.Repository
}

func New(repo repository.Repository) Service {
	return &reportService{repo: repo}
}

func (s *reportService) GetSalesReport(ctx context.Context, startDate, endDate *string) (*model.SalesReport, error) {
	now := time.Now()
	start := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
	end := start.Add(24 * time.Hour).Add(-time.Nanosecond)

	if startDate != nil {
		parsed, err := time.Parse("2006-01-02", *startDate)
		if err != nil {
			return nil, err
		}
		start = time.Date(parsed.Year(), parsed.Month(), parsed.Day(), 0, 0, 0, 0, parsed.Location())
	}

	if endDate != nil {
		parsed, err := time.Parse("2006-01-02", *endDate)
		if err != nil {
			return nil, err
		}
		end = time.Date(parsed.Year(), parsed.Month(), parsed.Day(), 23, 59, 59, 999999999, parsed.Location())
	}

	return s.repo.GetSalesReport(ctx, &start, &end)
}

func (s *reportService) GetTodayReport(ctx context.Context) (*model.SalesReport, error) {
	now := time.Now()
	start := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
	end := start.Add(24 * time.Hour).Add(-time.Nanosecond)
	return s.repo.GetSalesReport(ctx, &start, &end)
}

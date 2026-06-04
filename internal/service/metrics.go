package service

import (
	"context"

	"github.com/churilovmn1/workout-tracker/internal/models"
)

type metricsRepository interface {
	Create(ctx context.Context, m *models.BodyMetric) (int, error)
	ListByUser(ctx context.Context, userID int) ([]models.BodyMetric, error)
	Delete(ctx context.Context, id, userID int) error
}

// MetricsService handles body-metric business logic.
type MetricsService struct {
	repo metricsRepository
}

// NewMetricsService creates a new MetricsService.
func NewMetricsService(repo metricsRepository) *MetricsService {
	return &MetricsService{repo: repo}
}

// Create adds a new body metric snapshot.
func (s *MetricsService) Create(ctx context.Context, m *models.BodyMetric) (int, error) {
	return s.repo.Create(ctx, m)
}

// ListByUser returns all metrics for a user.
func (s *MetricsService) ListByUser(ctx context.Context, userID int) ([]models.BodyMetric, error) {
	return s.repo.ListByUser(ctx, userID)
}

// Delete removes a metric owned by the user.
func (s *MetricsService) Delete(ctx context.Context, id, userID int) error {
	return s.repo.Delete(ctx, id, userID)
}

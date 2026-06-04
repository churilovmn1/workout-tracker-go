package repository

import (
	"context"
	"fmt"

	"github.com/churilovmn1/workout-tracker/internal/models"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

// MetricsRepository handles database operations for body metrics.
type MetricsRepository struct {
	pool *pgxpool.Pool
}

// NewMetricsRepository creates a new MetricsRepository.
func NewMetricsRepository(pool *pgxpool.Pool) *MetricsRepository {
	return &MetricsRepository{pool: pool}
}

// Create inserts a new body metric snapshot and returns its ID.
func (r *MetricsRepository) Create(ctx context.Context, m *models.BodyMetric) (int, error) {
	var id int
	err := r.pool.QueryRow(ctx,
		`INSERT INTO body_metrics
		   (user_id, weight_kg, body_fat_percent, chest_cm, waist_cm, hips_cm, bicep_cm, measured_at)
		 VALUES ($1,$2,$3,$4,$5,$6,$7,$8)
		 RETURNING id`,
		m.UserID, m.WeightKg, m.BodyFatPercent, m.ChestCm, m.WaistCm, m.HipsCm, m.BicepCm, m.MeasuredAt,
	).Scan(&id)
	if err != nil {
		return 0, fmt.Errorf("create body metric: %w", err)
	}
	return id, nil
}

// ListByUser returns all metrics for a user ordered oldest-first (for charts).
func (r *MetricsRepository) ListByUser(ctx context.Context, userID int) ([]models.BodyMetric, error) {
	rows, err := r.pool.Query(ctx,
		`SELECT id, user_id, weight_kg, body_fat_percent, chest_cm, waist_cm, hips_cm, bicep_cm,
		        measured_at, created_at
		 FROM body_metrics
		 WHERE user_id = $1
		 ORDER BY measured_at ASC, id ASC`,
		userID,
	)
	if err != nil {
		return nil, fmt.Errorf("list body metrics: %w", err)
	}
	defer rows.Close()

	return pgx.CollectRows(rows, func(row pgx.CollectableRow) (models.BodyMetric, error) {
		var m models.BodyMetric
		return m, row.Scan(&m.ID, &m.UserID, &m.WeightKg, &m.BodyFatPercent,
			&m.ChestCm, &m.WaistCm, &m.HipsCm, &m.BicepCm, &m.MeasuredAt, &m.CreatedAt)
	})
}

// Delete removes a metric owned by the user.
func (r *MetricsRepository) Delete(ctx context.Context, id, userID int) error {
	tag, err := r.pool.Exec(ctx,
		`DELETE FROM body_metrics WHERE id = $1 AND user_id = $2`, id, userID)
	if err != nil {
		return fmt.Errorf("delete body metric: %w", err)
	}
	if tag.RowsAffected() == 0 {
		return fmt.Errorf("body metric not found")
	}
	return nil
}

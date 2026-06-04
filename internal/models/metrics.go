package models

import "time"

// BodyMetric is a single body-measurement snapshot for a user.
type BodyMetric struct {
	ID             int       `json:"id"`
	UserID         int       `json:"user_id"`
	WeightKg       *float64  `json:"weight_kg"`
	BodyFatPercent *float64  `json:"body_fat_percent"`
	ChestCm        *float64  `json:"chest_cm"`
	WaistCm        *float64  `json:"waist_cm"`
	HipsCm         *float64  `json:"hips_cm"`
	BicepCm        *float64  `json:"bicep_cm"`
	MeasuredAt     time.Time `json:"measured_at"`
	CreatedAt      time.Time `json:"created_at"`
}

package handler

import (
	"net/http"
	"strconv"
	"time"

	"github.com/churilovmn1/workout-tracker/internal/models"
	"github.com/churilovmn1/workout-tracker/internal/service"
	"github.com/go-chi/chi/v5"
)

// MetricsHandler handles body-metric endpoints.
type MetricsHandler struct {
	metricsService *service.MetricsService
}

// NewMetricsHandler creates a new MetricsHandler.
func NewMetricsHandler(metricsService *service.MetricsService) *MetricsHandler {
	return &MetricsHandler{metricsService: metricsService}
}

type metricRequest struct {
	WeightKg       *float64 `json:"weight_kg"`
	BodyFatPercent *float64 `json:"body_fat_percent"`
	ChestCm        *float64 `json:"chest_cm"`
	WaistCm        *float64 `json:"waist_cm"`
	HipsCm         *float64 `json:"hips_cm"`
	BicepCm        *float64 `json:"bicep_cm"`
	MeasuredAt     string   `json:"measured_at"` // "2006-01-02"
}

// List returns all body metrics for the authenticated user.
func (h *MetricsHandler) List(w http.ResponseWriter, r *http.Request) {
	metrics, err := h.metricsService.ListByUser(r.Context(), getUserID(r))
	if err != nil {
		writeError(w, http.StatusInternalServerError, "failed to list metrics")
		return
	}
	writeJSON(w, http.StatusOK, metrics)
}

// Create adds a new body metric snapshot.
func (h *MetricsHandler) Create(w http.ResponseWriter, r *http.Request) {
	var req metricRequest
	if err := decodeJSON(r, &req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	measuredAt := time.Now()
	if req.MeasuredAt != "" {
		if t, err := time.Parse("2006-01-02", req.MeasuredAt); err == nil {
			measuredAt = t
		}
	}

	m := &models.BodyMetric{
		UserID:         getUserID(r),
		WeightKg:       req.WeightKg,
		BodyFatPercent: req.BodyFatPercent,
		ChestCm:        req.ChestCm,
		WaistCm:        req.WaistCm,
		HipsCm:         req.HipsCm,
		BicepCm:        req.BicepCm,
		MeasuredAt:     measuredAt,
	}

	id, err := h.metricsService.Create(r.Context(), m)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "failed to create metric")
		return
	}
	m.ID = id
	writeJSON(w, http.StatusCreated, m)
}

// Delete removes a body metric owned by the user.
func (h *MetricsHandler) Delete(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid metric id")
		return
	}
	if err := h.metricsService.Delete(r.Context(), id, getUserID(r)); err != nil {
		writeError(w, http.StatusNotFound, "metric not found")
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

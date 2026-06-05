package handler

import (
	"errors"
	"net/http"
	"strconv"
	"time"

	"github.com/churilovmn1/workout-tracker/internal/broker"
	"github.com/churilovmn1/workout-tracker/internal/models"
	"github.com/churilovmn1/workout-tracker/internal/service"
	"github.com/go-chi/chi/v5"
)

// AdminHandler handles trainer/admin endpoints.
type AdminHandler struct {
	adminService   *service.AdminService
	workoutService *service.WorkoutService
	metricsService *service.MetricsService
	publisher      broker.Publisher
}

// NewAdminHandler creates a new AdminHandler.
func NewAdminHandler(
	adminService *service.AdminService,
	workoutService *service.WorkoutService,
	metricsService *service.MetricsService,
	publisher broker.Publisher,
) *AdminHandler {
	return &AdminHandler{
		adminService:   adminService,
		workoutService: workoutService,
		metricsService: metricsService,
		publisher:      publisher,
	}
}

type commentRequest struct {
	Comment string `json:"comment"`
}

type scheduleCreateRequest struct {
	ClientID        int    `json:"client_id"`
	Title           string `json:"title"`
	ScheduledAt     string `json:"scheduled_at"` // "2026-01-05T10:00"
	DurationMinutes int    `json:"duration_minutes"`
	Notes           string `json:"notes"`
}

type scheduleUpdateRequest struct {
	Status          string `json:"status"`
	Title           string `json:"title"`
	ScheduledAt     string `json:"scheduled_at"`
	DurationMinutes int    `json:"duration_minutes"`
	Notes           string `json:"notes"`
	ClientID        int    `json:"client_id"`
}

// mondayOf returns the Monday of the week containing t (UTC).
func mondayOf(t time.Time) time.Time {
	t = t.UTC().Truncate(24 * time.Hour)
	wd := t.Weekday()
	if wd == time.Sunday {
		wd = 7
	}
	return t.AddDate(0, 0, -int(wd)+1)
}

// ── Users ──────────────────────────────────────────────────────────────────

// ListUsers returns all registered users.
//
// @Summary      List all users
// @Tags         admin
// @Produce      json
// @Success      200  {array}   models.User
// @Failure      500  {object}  errorResponse
// @Security     BearerAuth
// @Router       /admin/users [get]
func (h *AdminHandler) ListUsers(w http.ResponseWriter, r *http.Request) {
	users, err := h.adminService.ListUsers(r.Context())
	if err != nil {
		writeError(w, http.StatusInternalServerError, "failed to list users")
		return
	}
	writeJSON(w, http.StatusOK, users)
}

// ListUserWorkouts returns all workouts for a specific user.
//
// @Summary      List client workouts
// @Tags         admin
// @Produce      json
// @Param        id   path      int  true  "Client user ID"
// @Success      200  {array}   models.Workout
// @Failure      400  {object}  errorResponse
// @Security     BearerAuth
// @Router       /admin/users/{id}/workouts [get]
func (h *AdminHandler) ListUserWorkouts(w http.ResponseWriter, r *http.Request) {
	userID, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid user id")
		return
	}

	workouts, err := h.adminService.ListUserWorkouts(r.Context(), userID)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "failed to list workouts")
		return
	}
	writeJSON(w, http.StatusOK, workouts)
}

// SetComment sets a trainer comment on a workout.
//
// @Summary      Set trainer comment
// @Tags         admin
// @Accept       json
// @Produce      json
// @Param        id    path      int             true  "Workout ID"
// @Param        body  body      commentRequest  true  "Comment"
// @Success      200   {object}  map[string]bool
// @Failure      400   {object}  errorResponse
// @Security     BearerAuth
// @Router       /admin/workouts/{id}/comment [put]
func (h *AdminHandler) SetComment(w http.ResponseWriter, r *http.Request) {
	workoutID, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid workout id")
		return
	}

	var req commentRequest
	if err := decodeJSON(r, &req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	if err := h.adminService.SetTrainerComment(r.Context(), workoutID, req.Comment); err != nil {
		writeError(w, http.StatusInternalServerError, "failed to set comment")
		return
	}

	_ = h.publisher.Publish(r.Context(), broker.NewEvent(broker.EventWorkoutCommented, broker.Payload{
		WorkoutID: workoutID,
		Comment:   req.Comment,
	}))

	writeJSON(w, http.StatusOK, map[string]bool{"ok": true})
}

// CreateWorkoutForUser creates a workout on behalf of a user.
//
// @Summary      Create workout for client
// @Tags         admin
// @Accept       json
// @Produce      json
// @Param        id    path      int             true  "Client user ID"
// @Param        body  body      workoutRequest  true  "Workout data"
// @Success      201   {object}  models.Workout
// @Failure      400   {object}  errorResponse
// @Security     BearerAuth
// @Router       /admin/users/{id}/workouts [post]
func (h *AdminHandler) CreateWorkoutForUser(w http.ResponseWriter, r *http.Request) {
	userID, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid user id")
		return
	}

	var req workoutRequest
	if err := decodeJSON(r, &req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	if req.Title == "" {
		writeError(w, http.StatusBadRequest, "title is required")
		return
	}

	date, err := time.Parse("2006-01-02", req.Date)
	if err != nil {
		date = time.Now()
	}

	workout := &models.Workout{
		UserID:          userID,
		Title:           req.Title,
		Date:            date,
		DurationMinutes: req.DurationMinutes,
		Notes:           req.Notes,
	}

	for _, e := range req.Exercises {
		workout.Exercises = append(workout.Exercises, models.WorkoutExercise{
			ExerciseID: e.ExerciseID,
			Sets:       e.Sets,
			Reps:       e.Reps,
			WeightKg:   e.WeightKg,
		})
	}

	id, err := h.adminService.CreateWorkoutForUser(r.Context(), workout)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "failed to create workout")
		return
	}

	workout.ID = id
	writeJSON(w, http.StatusCreated, workout)
}

// ── Schedule ───────────────────────────────────────────────────────────────

// ListSchedule returns schedule entries for the authenticated trainer's week.
//
// @Summary      List trainer schedule
// @Tags         admin
// @Produce      json
// @Param        week  query     string  false  "Week start date (YYYY-MM-DD). Defaults to current week."
// @Success      200   {array}   models.ScheduleEntry
// @Security     BearerAuth
// @Router       /admin/schedule [get]
func (h *AdminHandler) ListSchedule(w http.ResponseWriter, r *http.Request) {
	weekStr := r.URL.Query().Get("week")
	var base time.Time
	var parseErr error
	if weekStr != "" {
		base, parseErr = time.Parse("2006-01-02", weekStr)
	}
	if weekStr == "" || parseErr != nil {
		base = time.Now()
	}

	start := mondayOf(base)
	end := start.AddDate(0, 0, 7)

	entries, err := h.adminService.ListScheduleWeek(r.Context(), getUserID(r), start, end)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "failed to list schedule")
		return
	}
	writeJSON(w, http.StatusOK, entries)
}

// CreateSchedule adds a new schedule entry.
//
// @Summary      Create schedule entry
// @Tags         admin
// @Accept       json
// @Produce      json
// @Param        body  body      scheduleCreateRequest  true  "Schedule data"
// @Success      201   {object}  models.ScheduleEntry
// @Failure      400   {object}  errorResponse
// @Security     BearerAuth
// @Router       /admin/schedule [post]
func (h *AdminHandler) CreateSchedule(w http.ResponseWriter, r *http.Request) {
	var req scheduleCreateRequest
	if err := decodeJSON(r, &req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	if req.Title == "" || req.ClientID == 0 || req.ScheduledAt == "" {
		writeError(w, http.StatusBadRequest, "title, client_id and scheduled_at are required")
		return
	}

	at, err := time.Parse("2006-01-02T15:04", req.ScheduledAt)
	if err != nil {
		writeError(w, http.StatusBadRequest, "scheduled_at must be in format 2006-01-02T15:04")
		return
	}

	dur := req.DurationMinutes
	if dur <= 0 {
		dur = 60
	}

	entry := &models.ScheduleEntry{
		TrainerID:       getUserID(r),
		ClientID:        req.ClientID,
		Title:           req.Title,
		ScheduledAt:     at,
		DurationMinutes: dur,
		Notes:           req.Notes,
	}

	id, err := h.adminService.CreateScheduleEntry(r.Context(), entry)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "failed to create schedule entry")
		return
	}

	entry.ID = id
	writeJSON(w, http.StatusCreated, entry)
}

// UpdateSchedule modifies status or fields of a schedule entry.
//
// @Summary      Update schedule entry
// @Tags         admin
// @Accept       json
// @Produce      json
// @Param        id    path      int                    true  "Schedule entry ID"
// @Param        body  body      scheduleUpdateRequest  true  "Updated fields"
// @Success      200   {object}  models.ScheduleEntry
// @Failure      400   {object}  errorResponse
// @Failure      403   {object}  errorResponse
// @Failure      404   {object}  errorResponse
// @Security     BearerAuth
// @Router       /admin/schedule/{id} [put]
func (h *AdminHandler) UpdateSchedule(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid schedule id")
		return
	}

	trainerID := getUserID(r)

	existing, err := h.adminService.GetScheduleEntry(r.Context(), id)
	if err != nil {
		writeError(w, http.StatusNotFound, "schedule entry not found")
		return
	}
	if existing.TrainerID != trainerID {
		writeError(w, http.StatusForbidden, "access denied")
		return
	}

	var req scheduleUpdateRequest
	if err := decodeJSON(r, &req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	if req.Status != "" {
		existing.Status = req.Status
	}
	if req.Title != "" {
		existing.Title = req.Title
	}
	if req.DurationMinutes > 0 {
		existing.DurationMinutes = req.DurationMinutes
	}
	if req.Notes != "" {
		existing.Notes = req.Notes
	}
	if req.ClientID > 0 {
		existing.ClientID = req.ClientID
	}
	if req.ScheduledAt != "" {
		at, err := time.Parse("2006-01-02T15:04", req.ScheduledAt)
		if err == nil {
			existing.ScheduledAt = at
		}
	}

	if err := h.adminService.UpdateScheduleEntry(r.Context(), existing); err != nil {
		writeError(w, http.StatusInternalServerError, "failed to update schedule entry")
		return
	}

	writeJSON(w, http.StatusOK, existing)
}

// DeleteSchedule removes a schedule entry.
//
// @Summary      Delete schedule entry
// @Tags         admin
// @Param        id   path  int  true  "Schedule entry ID"
// @Success      204
// @Failure      400  {object}  errorResponse
// @Failure      403  {object}  errorResponse
// @Security     BearerAuth
// @Router       /admin/schedule/{id} [delete]
func (h *AdminHandler) DeleteSchedule(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid schedule id")
		return
	}

	if err := h.adminService.DeleteScheduleEntry(r.Context(), id, getUserID(r)); err != nil {
		if errors.Is(err, service.ErrForbidden) {
			writeError(w, http.StatusForbidden, "access denied")
			return
		}
		writeError(w, http.StatusInternalServerError, "failed to delete schedule entry")
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// ── Client analytics (trainer view) ───────────────────────────────────────

// GetClientMetrics returns body metrics for a specific client.
//
// @Summary      Client body metrics
// @Tags         admin
// @Produce      json
// @Param        id   path      int  true  "Client user ID"
// @Success      200  {array}   models.BodyMetric
// @Failure      400  {object}  errorResponse
// @Security     BearerAuth
// @Router       /admin/users/{id}/metrics [get]
func (h *AdminHandler) GetClientMetrics(w http.ResponseWriter, r *http.Request) {
	clientID, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid user id")
		return
	}
	metrics, err := h.metricsService.ListByUser(r.Context(), clientID)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "failed to get client metrics")
		return
	}
	writeJSON(w, http.StatusOK, metrics)
}

// GetClientExerciseProgress returns exercise weight history for a specific client.
//
// @Summary      Client exercise progress
// @Tags         admin
// @Produce      json
// @Param        id           path      int  true  "Client user ID"
// @Param        exercise_id  query     int  true  "Exercise ID"
// @Success      200          {array}   models.ExerciseProgress
// @Failure      400          {object}  errorResponse
// @Security     BearerAuth
// @Router       /admin/users/{id}/exercise-progress [get]
func (h *AdminHandler) GetClientExerciseProgress(w http.ResponseWriter, r *http.Request) {
	clientID, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid user id")
		return
	}
	exerciseID, err := strconv.Atoi(r.URL.Query().Get("exercise_id"))
	if err != nil || exerciseID <= 0 {
		writeError(w, http.StatusBadRequest, "exercise_id is required")
		return
	}
	progress, err := h.workoutService.GetExerciseProgress(r.Context(), clientID, exerciseID)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "failed to get client exercise progress")
		return
	}
	writeJSON(w, http.StatusOK, progress)
}

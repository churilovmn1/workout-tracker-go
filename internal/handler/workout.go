package handler

import (
	"errors"
	"net/http"
	"strconv"
	"time"

	"github.com/churilovmn1/workout-tracker/internal/models"
	"github.com/churilovmn1/workout-tracker/internal/service"
	"github.com/go-chi/chi/v5"
)

// WorkoutHandler обрабатывает запросы к тренировкам.
type WorkoutHandler struct {
	workoutService *service.WorkoutService
}

// NewWorkoutHandler создаёт WorkoutHandler.
func NewWorkoutHandler(workoutService *service.WorkoutService) *WorkoutHandler {
	return &WorkoutHandler{workoutService: workoutService}
}

type workoutExerciseRequest struct {
	ExerciseID int     `json:"exercise_id"`
	Sets       int     `json:"sets"`
	Reps       int     `json:"reps"`
	WeightKg   float64 `json:"weight_kg"`
}

type workoutRequest struct {
	Title           string                   `json:"title"`
	Date            string                   `json:"date"`
	DurationMinutes int                      `json:"duration_minutes"`
	Notes           string                   `json:"notes"`
	Exercises       []workoutExerciseRequest `json:"exercises"`
}

type volumeResponse struct {
	WeeklyVolume float64 `json:"weekly_volume"`
}

// List возвращает все тренировки авторизованного пользователя.
//
// @Summary      List workouts
// @Tags         workouts
// @Produce      json
// @Success      200  {array}   models.Workout
// @Failure      500  {object}  errorResponse
// @Security     BearerAuth
// @Router       /workouts [get]
func (h *WorkoutHandler) List(w http.ResponseWriter, r *http.Request) {
	workouts, err := h.workoutService.ListByUser(r.Context(), getUserID(r))
	if err != nil {
		writeError(w, http.StatusInternalServerError, "failed to list workouts")
		return
	}

	writeJSON(w, http.StatusOK, workouts)
}

// GetByID возвращает тренировку по ID.
//
// @Summary      Get workout
// @Tags         workouts
// @Produce      json
// @Param        id   path      int  true  "Workout ID"
// @Success      200  {object}  models.Workout
// @Failure      403  {object}  errorResponse
// @Failure      404  {object}  errorResponse
// @Security     BearerAuth
// @Router       /workouts/{id} [get]
func (h *WorkoutHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid workout id")
		return
	}

	// Сервис проверяет, что тренировка принадлежит запрашивающему пользователю.
	// Если нет — возвращает service.ErrForbidden, который маппируется в 403.
	workout, err := h.workoutService.GetByID(r.Context(), id, getUserID(r))
	if err != nil {
		if errors.Is(err, service.ErrForbidden) {
			writeError(w, http.StatusForbidden, "access denied")
			return
		}
		writeError(w, http.StatusNotFound, "workout not found")
		return
	}

	writeJSON(w, http.StatusOK, workout)
}

// Create создаёт новую тренировку с упражнениями.
//
// @Summary      Create workout
// @Tags         workouts
// @Accept       json
// @Produce      json
// @Param        body  body      workoutRequest  true  "Workout data"
// @Success      201   {object}  models.Workout
// @Failure      400   {object}  errorResponse
// @Security     BearerAuth
// @Router       /workouts [post]
func (h *WorkoutHandler) Create(w http.ResponseWriter, r *http.Request) {
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
		UserID:          getUserID(r),
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

	// WorkoutRepository.Create использует транзакцию: сначала вставляет workouts,
	// затем все workout_exercises — атомарно. При ошибке любого INSERT — rollback.
	id, err := h.workoutService.Create(r.Context(), workout)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "failed to create workout")
		return
	}

	workout.ID = id
	writeJSON(w, http.StatusCreated, workout)
}

// Update изменяет существующую тренировку.
//
// @Summary      Update workout
// @Tags         workouts
// @Accept       json
// @Produce      json
// @Param        id    path      int             true  "Workout ID"
// @Param        body  body      workoutRequest  true  "Workout data"
// @Success      200   {object}  models.Workout
// @Failure      400   {object}  errorResponse
// @Failure      403   {object}  errorResponse
// @Security     BearerAuth
// @Router       /workouts/{id} [put]
func (h *WorkoutHandler) Update(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid workout id")
		return
	}

	var req workoutRequest
	if err := decodeJSON(r, &req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	date, err := time.Parse("2006-01-02", req.Date)
	if err != nil {
		date = time.Now()
	}

	workout := &models.Workout{
		ID:              id,
		UserID:          getUserID(r),
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

	if err := h.workoutService.Update(r.Context(), workout); err != nil {
		if errors.Is(err, service.ErrForbidden) {
			writeError(w, http.StatusForbidden, "access denied")
			return
		}
		writeError(w, http.StatusInternalServerError, "failed to update workout")
		return
	}

	writeJSON(w, http.StatusOK, workout)
}

// Delete удаляет тренировку.
//
// @Summary      Delete workout
// @Tags         workouts
// @Param        id   path  int  true  "Workout ID"
// @Success      204
// @Failure      400  {object}  errorResponse
// @Security     BearerAuth
// @Router       /workouts/{id} [delete]
func (h *WorkoutHandler) Delete(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid workout id")
		return
	}

	if err := h.workoutService.Delete(r.Context(), id, getUserID(r)); err != nil {
		writeError(w, http.StatusInternalServerError, "failed to delete workout")
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// Copy создаёт копию существующей тренировки.
//
// @Summary      Copy workout
// @Tags         workouts
// @Produce      json
// @Param        id   path      int  true  "Source workout ID"
// @Success      201  {object}  map[string]int
// @Failure      403  {object}  errorResponse
// @Security     BearerAuth
// @Router       /workouts/{id}/copy [post]
func (h *WorkoutHandler) Copy(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid workout id")
		return
	}

	newID, err := h.workoutService.CopyWorkout(r.Context(), id, getUserID(r))
	if err != nil {
		if errors.Is(err, service.ErrForbidden) {
			writeError(w, http.StatusForbidden, "access denied")
			return
		}
		writeError(w, http.StatusInternalServerError, "failed to copy workout")
		return
	}

	writeJSON(w, http.StatusCreated, map[string]int{"id": newID})
}

// PersonalRecords возвращает лучший вес по каждому упражнению.
//
// @Summary      Personal records
// @Tags         stats
// @Produce      json
// @Success      200  {array}   models.WorkoutExercise
// @Security     BearerAuth
// @Router       /stats/pr [get]
func (h *WorkoutHandler) PersonalRecords(w http.ResponseWriter, r *http.Request) {
	records, err := h.workoutService.GetPersonalRecords(r.Context(), getUserID(r))
	if err != nil {
		writeError(w, http.StatusInternalServerError, "failed to get personal records")
		return
	}

	writeJSON(w, http.StatusOK, records)
}

// ExerciseProgress возвращает историю максимального веса по упражнению (для графика).
//
// @Summary      Exercise progress history
// @Tags         stats
// @Produce      json
// @Param        exercise_id  query     int  true  "Exercise ID"
// @Success      200          {array}   models.ExerciseProgress
// @Failure      400          {object}  errorResponse
// @Security     BearerAuth
// @Router       /stats/exercise-progress [get]
func (h *WorkoutHandler) ExerciseProgress(w http.ResponseWriter, r *http.Request) {
	exerciseID, err := strconv.Atoi(r.URL.Query().Get("exercise_id"))
	if err != nil || exerciseID <= 0 {
		writeError(w, http.StatusBadRequest, "exercise_id is required")
		return
	}

	progress, err := h.workoutService.GetExerciseProgress(r.Context(), getUserID(r), exerciseID)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "failed to get exercise progress")
		return
	}

	writeJSON(w, http.StatusOK, progress)
}

// WeeklyVolume возвращает суммарный объём за последние 7 дней.
//
// @Summary      Weekly volume
// @Tags         stats
// @Produce      json
// @Success      200  {object}  volumeResponse
// @Security     BearerAuth
// @Router       /stats/volume [get]
func (h *WorkoutHandler) WeeklyVolume(w http.ResponseWriter, r *http.Request) {
	volume, err := h.workoutService.GetWeeklyVolume(r.Context(), getUserID(r))
	if err != nil {
		writeError(w, http.StatusInternalServerError, "failed to get weekly volume")
		return
	}

	writeJSON(w, http.StatusOK, volumeResponse{WeeklyVolume: volume})
}

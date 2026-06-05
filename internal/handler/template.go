package handler

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/churilovmn1/workout-tracker/internal/models"
	"github.com/churilovmn1/workout-tracker/internal/service"
	"github.com/go-chi/chi/v5"
)

// TemplateHandler handles workout template endpoints.
type TemplateHandler struct {
	templateService *service.TemplateService
	workoutService  *service.WorkoutService
}

// NewTemplateHandler creates a new TemplateHandler.
func NewTemplateHandler(ts *service.TemplateService, ws *service.WorkoutService) *TemplateHandler {
	return &TemplateHandler{templateService: ts, workoutService: ws}
}

type templateExerciseRequest struct {
	ExerciseID int     `json:"exercise_id"`
	Sets       int     `json:"sets"`
	Reps       int     `json:"reps"`
	WeightKg   float64 `json:"weight_kg"`
}

type templateRequest struct {
	Name      string                    `json:"name"`
	IsPublic  bool                      `json:"is_public"`
	Exercises []templateExerciseRequest `json:"exercises"`
}

// List returns templates available to the authenticated user.
//
// @Summary      List templates
// @Tags         templates
// @Produce      json
// @Success      200  {array}   models.WorkoutTemplate
// @Security     BearerAuth
// @Router       /templates [get]
func (h *TemplateHandler) List(w http.ResponseWriter, r *http.Request) {
	templates, err := h.templateService.ListByUser(r.Context(), getUserID(r))
	if err != nil {
		writeError(w, http.StatusInternalServerError, "failed to list templates")
		return
	}

	writeJSON(w, http.StatusOK, templates)
}

// GetByID returns a template by ID if the user has access.
//
// @Summary      Get template
// @Tags         templates
// @Produce      json
// @Param        id   path      int  true  "Template ID"
// @Success      200  {object}  models.WorkoutTemplate
// @Failure      403  {object}  errorResponse
// @Failure      404  {object}  errorResponse
// @Security     BearerAuth
// @Router       /templates/{id} [get]
func (h *TemplateHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid template id")
		return
	}

	tmpl, err := h.templateService.GetByID(r.Context(), id, getUserID(r))
	if err != nil {
		if errors.Is(err, service.ErrForbidden) {
			writeError(w, http.StatusForbidden, "access denied")
			return
		}
		writeError(w, http.StatusNotFound, "template not found")
		return
	}

	writeJSON(w, http.StatusOK, tmpl)
}

// Create adds a new template.
//
// @Summary      Create template
// @Tags         templates
// @Accept       json
// @Produce      json
// @Param        body  body      templateRequest        true  "Template data"
// @Success      201   {object}  models.WorkoutTemplate
// @Failure      400   {object}  errorResponse
// @Security     BearerAuth
// @Router       /templates [post]
func (h *TemplateHandler) Create(w http.ResponseWriter, r *http.Request) {
	var req templateRequest
	if err := decodeJSON(r, &req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	if req.Name == "" {
		writeError(w, http.StatusBadRequest, "name is required")
		return
	}

	tmpl := &models.WorkoutTemplate{
		UserID:   getUserID(r),
		Name:     req.Name,
		IsPublic: req.IsPublic,
	}

	for _, e := range req.Exercises {
		tmpl.Exercises = append(tmpl.Exercises, models.TemplateExercise{
			ExerciseID: e.ExerciseID,
			Sets:       e.Sets,
			Reps:       e.Reps,
			WeightKg:   e.WeightKg,
		})
	}

	id, err := h.templateService.Create(r.Context(), tmpl)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "failed to create template")
		return
	}

	tmpl.ID = id
	writeJSON(w, http.StatusCreated, tmpl)
}

// Update modifies a template.
//
// @Summary      Update template
// @Tags         templates
// @Accept       json
// @Produce      json
// @Param        id    path      int             true  "Template ID"
// @Param        body  body      templateRequest true  "Template data"
// @Success      200   {object}  models.WorkoutTemplate
// @Security     BearerAuth
// @Router       /templates/{id} [put]
func (h *TemplateHandler) Update(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid template id")
		return
	}

	var req templateRequest
	if err := decodeJSON(r, &req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	tmpl := &models.WorkoutTemplate{
		ID:       id,
		UserID:   getUserID(r),
		Name:     req.Name,
		IsPublic: req.IsPublic,
	}

	for _, e := range req.Exercises {
		tmpl.Exercises = append(tmpl.Exercises, models.TemplateExercise{
			ExerciseID: e.ExerciseID,
			Sets:       e.Sets,
			Reps:       e.Reps,
			WeightKg:   e.WeightKg,
		})
	}

	if err := h.templateService.Update(r.Context(), tmpl); err != nil {
		writeError(w, http.StatusInternalServerError, "failed to update template")
		return
	}

	writeJSON(w, http.StatusOK, tmpl)
}

// Delete removes a template.
//
// @Summary      Delete template
// @Tags         templates
// @Param        id   path  int  true  "Template ID"
// @Success      204
// @Security     BearerAuth
// @Router       /templates/{id} [delete]
func (h *TemplateHandler) Delete(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid template id")
		return
	}

	if err := h.templateService.Delete(r.Context(), id, getUserID(r)); err != nil {
		writeError(w, http.StatusInternalServerError, "failed to delete template")
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// Start creates a workout for the authenticated user based on a template.
//
// @Summary      Start workout from template
// @Tags         templates
// @Produce      json
// @Param        id   path      int  true  "Template ID"
// @Success      201  {object}  models.Workout
// @Failure      403  {object}  errorResponse
// @Failure      404  {object}  errorResponse
// @Security     BearerAuth
// @Router       /templates/{id}/start [post]
func (h *TemplateHandler) Start(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid template id")
		return
	}

	userID := getUserID(r)
	tmpl, err := h.templateService.GetByID(r.Context(), id, userID)
	if err != nil {
		if errors.Is(err, service.ErrForbidden) {
			writeError(w, http.StatusForbidden, "access denied")
			return
		}
		writeError(w, http.StatusNotFound, "template not found")
		return
	}

	workout := h.templateService.CreateWorkoutFromTemplate(tmpl, userID)

	workoutID, err := h.workoutService.Create(r.Context(), workout)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "failed to create workout")
		return
	}

	workout.ID = workoutID
	writeJSON(w, http.StatusCreated, workout)
}

package repository

import (
	"context"
	"fmt"
	"strings"

	"github.com/churilovmn1/workout-tracker/internal/models"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

// exerciseFilter accumulates WHERE conditions and their positional args.
type exerciseFilter struct {
	conds []string
	args  []any
}

// ExerciseFilterOption configures an exercise query filter. Options are applied
// in order by buildExerciseFilter using the variadic ...Option pattern.
type ExerciseFilterOption func(*exerciseFilter)

// WithMuscleGroup filters by exact muscle group. A no-op when group is empty.
func WithMuscleGroup(group string) ExerciseFilterOption {
	return func(f *exerciseFilter) {
		if group == "" {
			return
		}
		f.args = append(f.args, group)
		f.conds = append(f.conds, fmt.Sprintf("muscle_group = $%d", len(f.args)))
	}
}

// WithSearch matches term as a case-insensitive substring against any of the
// given fields (OR-combined). A no-op when term or fields are empty. Field names
// are caller-supplied constants, never user input.
func WithSearch(term string, fields ...string) ExerciseFilterOption {
	return func(f *exerciseFilter) {
		if term == "" || len(fields) == 0 {
			return
		}
		f.args = append(f.args, "%"+term+"%")
		placeholder := fmt.Sprintf("$%d", len(f.args))
		parts := make([]string, len(fields))
		for i, field := range fields {
			parts[i] = field + " ILIKE " + placeholder
		}
		f.conds = append(f.conds, "("+strings.Join(parts, " OR ")+")")
	}
}

// buildExerciseFilter applies the options and returns a WHERE clause (empty when
// no conditions) together with the positional args.
func buildExerciseFilter(opts ...ExerciseFilterOption) (string, []any) {
	f := &exerciseFilter{}
	for _, opt := range opts {
		opt(f)
	}
	if len(f.conds) == 0 {
		return "", nil
	}
	return " WHERE " + strings.Join(f.conds, " AND "), f.args
}

// ExerciseRepository handles database operations for exercises.
type ExerciseRepository struct {
	pool *pgxpool.Pool
}

// NewExerciseRepository creates a new ExerciseRepository.
func NewExerciseRepository(pool *pgxpool.Pool) *ExerciseRepository {
	return &ExerciseRepository{pool: pool}
}

// Create inserts a new exercise and returns its ID.
func (r *ExerciseRepository) Create(ctx context.Context, ex *models.Exercise) (int, error) {
	var id int
	err := r.pool.QueryRow(ctx,
		`INSERT INTO exercises (name, muscle_group, description)
		 VALUES ($1, $2, $3)
		 RETURNING id`,
		ex.Name, ex.MuscleGroup, ex.Description,
	).Scan(&id)
	if err != nil {
		return 0, fmt.Errorf("create exercise: %w", err)
	}
	return id, nil
}

// GetByID returns an exercise by ID.
func (r *ExerciseRepository) GetByID(ctx context.Context, id int) (*models.Exercise, error) {
	ex := &models.Exercise{}
	err := r.pool.QueryRow(ctx,
		`SELECT id, name, muscle_group, description
		 FROM exercises WHERE id = $1`, id,
	).Scan(&ex.ID, &ex.Name, &ex.MuscleGroup, &ex.Description)
	if err != nil {
		return nil, fmt.Errorf("get exercise by id: %w", err)
	}
	return ex, nil
}

// List returns exercises filtered by muscle group and/or a search term matched
// against both the name and description fields.
func (r *ExerciseRepository) List(ctx context.Context, muscleGroup, search string) ([]models.Exercise, error) {
	where, args := buildExerciseFilter(
		WithMuscleGroup(muscleGroup),
		WithSearch(search, "name", "description"),
	)
	query := `SELECT id, name, muscle_group, description FROM exercises` +
		where + ` ORDER BY muscle_group, name`

	rows, err := r.pool.Query(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("list exercises: %w", err)
	}
	defer rows.Close()

	return pgx.CollectRows(rows, func(row pgx.CollectableRow) (models.Exercise, error) {
		var ex models.Exercise
		err := row.Scan(&ex.ID, &ex.Name, &ex.MuscleGroup, &ex.Description)
		return ex, err
	})
}

// Update modifies an existing exercise.
func (r *ExerciseRepository) Update(ctx context.Context, ex *models.Exercise) error {
	_, err := r.pool.Exec(ctx,
		`UPDATE exercises SET name = $1, muscle_group = $2, description = $3
		 WHERE id = $4`,
		ex.Name, ex.MuscleGroup, ex.Description, ex.ID,
	)
	if err != nil {
		return fmt.Errorf("update exercise: %w", err)
	}
	return nil
}

// Delete removes an exercise by ID.
func (r *ExerciseRepository) Delete(ctx context.Context, id int) error {
	_, err := r.pool.Exec(ctx, `DELETE FROM exercises WHERE id = $1`, id)
	if err != nil {
		return fmt.Errorf("delete exercise: %w", err)
	}
	return nil
}

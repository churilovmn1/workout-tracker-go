package repository

import (
	"context"
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/churilovmn1/workout-tracker/internal/models"
	"github.com/jackc/pgx/v5/pgxpool"
)

// benchPool returns a connection pool for benchmarks, skipping the benchmark
// when no database is reachable (CI without Postgres, local runs without docker).
// It honors TEST_DATABASE_URL, falling back to DATABASE_URL.
func benchPool(b *testing.B) *pgxpool.Pool {
	b.Helper()

	url := os.Getenv("TEST_DATABASE_URL")
	if url == "" {
		url = os.Getenv("DATABASE_URL")
	}
	if url == "" {
		b.Skip("set TEST_DATABASE_URL or DATABASE_URL to run repository benchmarks")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	pool, err := pgxpool.New(ctx, url)
	if err != nil {
		b.Skipf("cannot create pool: %v", err)
	}
	if err := pool.Ping(ctx); err != nil {
		pool.Close()
		b.Skipf("database unreachable: %v", err)
	}
	b.Cleanup(pool.Close)
	return pool
}

func benchExerciseName(i int) string {
	return fmt.Sprintf("bench-ex-%d-%d", time.Now().UnixNano(), i)
}

func BenchmarkExerciseRepository_Create(b *testing.B) {
	pool := benchPool(b)
	repo := NewExerciseRepository(pool)
	ctx := context.Background()

	ids := make([]int, 0, b.N)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		id, err := repo.Create(ctx, &models.Exercise{
			Name:        benchExerciseName(i),
			MuscleGroup: "Грудь",
			Description: "benchmark fixture",
		})
		if err != nil {
			b.Fatalf("create: %v", err)
		}
		ids = append(ids, id)
	}
	b.StopTimer()

	for _, id := range ids {
		_ = repo.Delete(ctx, id)
	}
}

func BenchmarkExerciseRepository_GetByID(b *testing.B) {
	pool := benchPool(b)
	repo := NewExerciseRepository(pool)
	ctx := context.Background()

	id, err := repo.Create(ctx, &models.Exercise{
		Name:        benchExerciseName(0),
		MuscleGroup: "Спина",
		Description: "benchmark fixture",
	})
	if err != nil {
		b.Fatalf("setup create: %v", err)
	}
	defer func() { _ = repo.Delete(ctx, id) }()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		if _, err := repo.GetByID(ctx, id); err != nil {
			b.Fatalf("get: %v", err)
		}
	}
}

func BenchmarkExerciseRepository_List(b *testing.B) {
	pool := benchPool(b)
	repo := NewExerciseRepository(pool)
	ctx := context.Background()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		if _, err := repo.List(ctx, "", ""); err != nil {
			b.Fatalf("list: %v", err)
		}
	}
}

func BenchmarkUserRepository_GetByID(b *testing.B) {
	pool := benchPool(b)
	repo := NewUserRepository(pool)
	ctx := context.Background()

	suffix := time.Now().UnixNano()
	id, err := repo.Create(ctx, &models.User{
		Login:        fmt.Sprintf("bench-user-%d", suffix),
		Email:        fmt.Sprintf("bench-%d@example.com", suffix),
		PasswordHash: "-",
		Role:         models.RoleUser,
	})
	if err != nil {
		b.Fatalf("setup create: %v", err)
	}
	defer func() { _ = repo.Delete(ctx, id) }()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		if _, err := repo.GetByID(ctx, id); err != nil {
			b.Fatalf("get: %v", err)
		}
	}
}

func BenchmarkWorkoutRepository_Create(b *testing.B) {
	pool := benchPool(b)
	userRepo := NewUserRepository(pool)
	workoutRepo := NewWorkoutRepository(pool)
	ctx := context.Background()

	suffix := time.Now().UnixNano()
	userID, err := userRepo.Create(ctx, &models.User{
		Login:        fmt.Sprintf("bench-wuser-%d", suffix),
		Email:        fmt.Sprintf("bench-w-%d@example.com", suffix),
		PasswordHash: "-",
		Role:         models.RoleUser,
	})
	if err != nil {
		b.Fatalf("setup user: %v", err)
	}
	defer func() { _ = userRepo.Delete(ctx, userID) }()

	ids := make([]int, 0, b.N)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		id, err := workoutRepo.Create(ctx, &models.Workout{
			UserID:          userID,
			Title:           "bench workout",
			Date:            time.Now(),
			DurationMinutes: 60,
		})
		if err != nil {
			b.Fatalf("create: %v", err)
		}
		ids = append(ids, id)
	}
	b.StopTimer()

	for _, id := range ids {
		_ = workoutRepo.Delete(ctx, id, userID)
	}
}

func BenchmarkWorkoutRepository_ListByUser(b *testing.B) {
	pool := benchPool(b)
	userRepo := NewUserRepository(pool)
	workoutRepo := NewWorkoutRepository(pool)
	ctx := context.Background()

	suffix := time.Now().UnixNano()
	userID, err := userRepo.Create(ctx, &models.User{
		Login:        fmt.Sprintf("bench-luser-%d", suffix),
		Email:        fmt.Sprintf("bench-l-%d@example.com", suffix),
		PasswordHash: "-",
		Role:         models.RoleUser,
	})
	if err != nil {
		b.Fatalf("setup user: %v", err)
	}
	defer func() { _ = userRepo.Delete(ctx, userID) }()

	for i := 0; i < 10; i++ {
		if _, err := workoutRepo.Create(ctx, &models.Workout{
			UserID:          userID,
			Title:           "bench workout",
			Date:            time.Now(),
			DurationMinutes: 45,
		}); err != nil {
			b.Fatalf("setup workout: %v", err)
		}
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		if _, err := workoutRepo.ListByUser(ctx, userID); err != nil {
			b.Fatalf("list: %v", err)
		}
	}
}

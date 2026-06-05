// @title           Workout Tracker API
// @version         1.0
// @description     REST API for tracking workouts, exercises, body metrics and trainer schedule.
// @host            localhost:8080
// @BasePath        /api
// @securityDefinitions.apikey BearerAuth
// @in              header
// @name            Authorization
// @description     JWT token. Format: "Bearer <token>"
package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/churilovmn1/workout-tracker/config"
	_ "github.com/churilovmn1/workout-tracker/docs"
	"github.com/churilovmn1/workout-tracker/internal/handler"
	"github.com/churilovmn1/workout-tracker/internal/repository"
	"github.com/churilovmn1/workout-tracker/internal/service"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("failed to load config: %v", err)
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Единственное место, где создаётся пул соединений с БД.
	// Все репозитории получают *pgxpool.Pool — он потокобезопасен и переиспользует соединения.
	pool, err := repository.NewPostgresPool(ctx, cfg.DatabaseURL)
	if err != nil {
		log.Fatalf("failed to connect to database: %v", err)
	}
	defer pool.Close()

	// ── Слой репозиториев ────────────────────────────────────────────────────
	// Каждый репозиторий оборачивает пул и предоставляет типизированные методы.
	// Конкретные типы удовлетворяют узким интерфейсам, объявленным в service/.
	userRepo := repository.NewUserRepository(pool)
	exerciseRepo := repository.NewExerciseRepository(pool)
	workoutRepo := repository.NewWorkoutRepository(pool)
	templateRepo := repository.NewTemplateRepository(pool)
	scheduleRepo := repository.NewScheduleRepository(pool)
	metricsRepo := repository.NewMetricsRepository(pool)

	// ── Слой сервисов ────────────────────────────────────────────────────────
	// Зависимости передаются явно через конструкторы — без глобального состояния
	// и DI-фреймворка. Это упрощает тестирование: в тестах вместо реального
	// репозитория передаётся мок (см. internal/service/*_test.go).
	authService := service.NewAuthService(userRepo, cfg.JWTSecret)
	exerciseService := service.NewExerciseService(exerciseRepo)
	workoutService := service.NewWorkoutService(workoutRepo)
	templateService := service.NewTemplateService(templateRepo)
	adminService := service.NewAdminService(userRepo, workoutRepo, scheduleRepo)
	metricsService := service.NewMetricsService(metricsRepo)

	// ── HTTP-роутер ──────────────────────────────────────────────────────────
	// NewRouter собирает chi-роутер: навешивает middleware (rate limiter, logger,
	// recoverer) и регистрирует все маршруты с нужными middleware-цепочками.
	router := handler.NewRouter(
		authService, exerciseService, workoutService,
		templateService, adminService, metricsService,
		"web",
	)

	srv := &http.Server{
		Addr:         ":" + cfg.Port,
		Handler:      router,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	go func() {
		log.Printf("server starting on :%s", cfg.Port)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("server error: %v", err)
		}
	}()

	// ── Graceful shutdown ────────────────────────────────────────────────────
	// Ждём SIGINT / SIGTERM, затем даём серверу 5 секунд завершить текущие запросы.
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("shutting down server...")
	cancel()

	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer shutdownCancel()

	if err := srv.Shutdown(shutdownCtx); err != nil {
		log.Fatalf("server forced shutdown: %v", err)
	}
	log.Println("server stopped")
}

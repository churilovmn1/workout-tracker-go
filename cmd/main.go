package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/churilovmn1/workout-tracker/bot"
	"github.com/churilovmn1/workout-tracker/config"
	"github.com/churilovmn1/workout-tracker/internal/broker"
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

	pool, err := repository.NewPostgresPool(ctx, cfg.DatabaseURL)
	if err != nil {
		log.Fatalf("failed to connect to database: %v", err)
	}
	defer pool.Close()

	userRepo := repository.NewUserRepository(pool)
	exerciseRepo := repository.NewExerciseRepository(pool)
	workoutRepo := repository.NewWorkoutRepository(pool)
	templateRepo := repository.NewTemplateRepository(pool)
	scheduleRepo := repository.NewScheduleRepository(pool)
	metricsRepo := repository.NewMetricsRepository(pool)

	authService := service.NewAuthService(userRepo, cfg.JWTSecret)
	exerciseService := service.NewExerciseService(exerciseRepo)
	workoutService := service.NewWorkoutService(workoutRepo)
	templateService := service.NewTemplateService(templateRepo)
	adminService := service.NewAdminService(userRepo, workoutRepo, scheduleRepo)
	metricsService := service.NewMetricsService(metricsRepo)

	var tgBot *bot.Bot
	if cfg.BotToken != "" {
		b, err := bot.New(cfg.BotToken, userRepo, workoutService, exerciseService, templateService)
		if err != nil {
			log.Printf("failed to create telegram bot: %v", err)
		} else {
			tgBot = b
			go tgBot.Start(ctx)
		}
	}

	// Message broker (optional). When REDIS_URL is set, events are queued in
	// Redis and a worker delivers Telegram notifications; otherwise publishing
	// is a no-op.
	var publisher broker.Publisher = broker.NoopPublisher{}
	if cfg.RedisURL != "" {
		rb, err := broker.Connect(ctx, cfg.RedisURL)
		if err != nil {
			log.Printf("failed to connect to redis: %v", err)
		} else {
			defer rb.Close()
			publisher = rb

			var notifier broker.Notifier
			if tgBot != nil {
				notifier = tgBot
			}
			worker := broker.NewWorker(rb, notifier, userRepo, workoutRepo, scheduleRepo)
			go worker.Run(ctx)
		}
	}

	router := handler.NewRouter(authService, exerciseService, workoutService, templateService, adminService, metricsService, publisher, "web")

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

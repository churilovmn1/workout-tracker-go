package handler

import (
	"net/http/pprof"
	"time"

	"github.com/churilovmn1/workout-tracker/internal/broker"
	"github.com/churilovmn1/workout-tracker/internal/service"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

// NewRouter sets up all HTTP routes.
func NewRouter(
	authService *service.AuthService,
	exerciseService *service.ExerciseService,
	workoutService *service.WorkoutService,
	templateService *service.TemplateService,
	adminService *service.AdminService,
	metricsService *service.MetricsService,
	publisher broker.Publisher,
	webDir string,
) chi.Router {
	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(NewRateLimiter(100, time.Minute).Middleware)

	authHandler := NewAuthHandler(authService)
	exerciseHandler := NewExerciseHandler(exerciseService)
	workoutHandler := NewWorkoutHandler(workoutService, publisher)
	templateHandler := NewTemplateHandler(templateService, workoutService)
	adminHandler := NewAdminHandler(adminService, publisher)
	metricsHandler := NewMetricsHandler(metricsService)
	webHandler := NewWebHandler(webDir)

	r.Get("/", webHandler.Index)
	r.Handle("/static/*", webHandler.StaticHandler())

	// pprof profiler — admin only.
	r.Route("/debug/pprof", func(r chi.Router) {
		r.Use(AuthMiddleware(authService))
		r.Use(AdminOnly)
		r.HandleFunc("/", pprof.Index)
		r.HandleFunc("/cmdline", pprof.Cmdline)
		r.HandleFunc("/profile", pprof.Profile)
		r.HandleFunc("/symbol", pprof.Symbol)
		r.HandleFunc("/trace", pprof.Trace)
		r.HandleFunc("/{profile}", pprof.Index) // heap, goroutine, allocs, ...
	})

	r.Route("/api", func(r chi.Router) {
		r.Route("/auth", func(r chi.Router) {
			r.Post("/register", authHandler.Register)
			r.Post("/login", authHandler.Login)
		})

		r.Group(func(r chi.Router) {
			r.Use(AuthMiddleware(authService))

			r.Route("/exercises", func(r chi.Router) {
				r.Get("/", exerciseHandler.List)
				r.Get("/{id}", exerciseHandler.GetByID)

				r.Group(func(r chi.Router) {
					r.Use(AdminOnly)
					r.Post("/", exerciseHandler.Create)
					r.Put("/{id}", exerciseHandler.Update)
					r.Delete("/{id}", exerciseHandler.Delete)
				})
			})

			r.Route("/workouts", func(r chi.Router) {
				r.Get("/", workoutHandler.List)
				r.Post("/", workoutHandler.Create)
				r.Get("/{id}", workoutHandler.GetByID)
				r.Put("/{id}", workoutHandler.Update)
				r.Delete("/{id}", workoutHandler.Delete)
				r.Post("/{id}/copy", workoutHandler.Copy)
			})

			r.Route("/templates", func(r chi.Router) {
				r.Get("/", templateHandler.List)
				r.Post("/", templateHandler.Create)
				r.Get("/{id}", templateHandler.GetByID)
				r.Put("/{id}", templateHandler.Update)
				r.Delete("/{id}", templateHandler.Delete)
				r.Post("/{id}/start", templateHandler.Start)
			})

			r.Route("/stats", func(r chi.Router) {
				r.Get("/pr", workoutHandler.PersonalRecords)
				r.Get("/volume", workoutHandler.WeeklyVolume)
				r.Get("/exercise-progress", workoutHandler.ExerciseProgress)
			})

			r.Route("/metrics", func(r chi.Router) {
				r.Get("/", metricsHandler.List)
				r.Post("/", metricsHandler.Create)
				r.Delete("/{id}", metricsHandler.Delete)
			})

			r.Route("/admin", func(r chi.Router) {
				r.Use(AdminOnly)
				r.Get("/users", adminHandler.ListUsers)
				r.Get("/users/{id}/workouts", adminHandler.ListUserWorkouts)
				r.Post("/users/{id}/workouts", adminHandler.CreateWorkoutForUser)
				r.Put("/workouts/{id}/comment", adminHandler.SetComment)
				r.Get("/schedule", adminHandler.ListSchedule)
				r.Post("/schedule", adminHandler.CreateSchedule)
				r.Put("/schedule/{id}", adminHandler.UpdateSchedule)
				r.Delete("/schedule/{id}", adminHandler.DeleteSchedule)
			})
		})
	})

	return r
}

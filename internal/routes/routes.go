package routes

import (
	"github.com/go-chi/chi/v5"
	"github.com/kevin120202/habit-tracker/internal/app"
)

func SetupRoutes(app *app.Application) *chi.Mux {
	r := chi.NewRouter()

	r.Get("/health", app.HealthCheck)
	r.Get("/habits/{id}", app.HabitHandler.HandleGetHabitByID)
	r.Post("/habits", app.HabitHandler.HandleCreateHabit)

	return r
}

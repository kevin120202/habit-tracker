package routes

import (
	"github.com/go-chi/chi/v5"
	"github.com/kevin120202/habit-tracker/internal/app"
)

func SetupRoutes(app *app.Application) *chi.Mux {
	r := chi.NewRouter()

	r.Get("/health", app.HealthCheck)

	r.Get("/habits", app.HabitHandler.HandleGetHabits)
	r.Get("/habits/{id}", app.HabitHandler.HandleGetHabitByID)
	r.Post("/habits", app.HabitHandler.HandleCreateHabit)
	r.Put("/habits/{id}", app.HabitHandler.HandleUpdateHabitByID)
	r.Delete("/habits/{id}", app.HabitHandler.HandleDeleteHabitByID)

	r.Post("/tags", app.TagHandler.HandleCreateTag)
	r.Get("/tags", app.TagHandler.HandleGetTags)

	return r
}

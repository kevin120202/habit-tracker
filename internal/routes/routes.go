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
	r.Post("/habits/{id}/log", app.HabitHandler.HandleLogHabitCompletions)
	r.Post("/habits/{id}/complete", app.HabitHandler.HandleCompleteHabit)
	r.Get("/habits/tags/{id}", app.HabitHandler.HandleGetHabitsByTag)
	r.Post("/habits/{id}/tags", app.HabitHandler.HandleCreateTagToHabit)
	r.Delete("/habits/{id}/tags/{tagID}", app.HabitHandler.HandleDeleteTagFromHabit)

	r.Post("/tags", app.TagHandler.HandleCreateTag)
	r.Get("/tags", app.TagHandler.HandleGetTags)
	r.Get("/tags/{id}", app.TagHandler.HandleGetTagByID)
	r.Put("/tags/{id}", app.TagHandler.HandleUpdateTagByID)
	r.Delete("/tags/{id}", app.TagHandler.HandleDeleteTagByID)

	return r
}

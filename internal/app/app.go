package app

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/kevin120202/habit-tracker/internal/api"
	"github.com/kevin120202/habit-tracker/internal/store"
	"github.com/kevin120202/habit-tracker/migrations"
)

type Application struct {
	Logger       *log.Logger
	HabitHandler *api.HabitHandler
	TagHandler   *api.TagHandler
	DB           *sql.DB
}

func NewApplication() (*Application, error) {
	pgDB, err := store.Open()
	if err != nil {
		return nil, err
	}

	err = store.MigrateFS(pgDB, migrations.FS, ".")
	if err != nil {
		panic(err)
	}

	logger := log.New(os.Stdout, "", log.Ldate|(log.Ltime))

	habitStore := store.NewPostgresHabitStore(pgDB)
	tagStore := store.NewPostgresTagStore(pgDB)

	habitHandler := api.NewHabitHandler(habitStore, logger)
	tagHandler := api.NewTagHandler(tagStore, logger)

	app := &Application{
		Logger:       logger,
		HabitHandler: habitHandler,
		TagHandler:   tagHandler,
		DB:           pgDB,
	}

	return app, nil
}

func (a *Application) HealthCheck(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Status is available\n")
}

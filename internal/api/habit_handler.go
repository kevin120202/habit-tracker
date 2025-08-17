package api

import (
	"fmt"
	"net/http"

	"github.com/kevin120202/habit-tracker/internal/utils"
)

type HabitHandler struct{}

func NewHabitHandler() *HabitHandler {
	return &HabitHandler{}
}

func (hh *HabitHandler) HandleGetHabitByID(w http.ResponseWriter, r *http.Request) {
	habitID, err := utils.ReadIDParam(r)
	if err != nil {
		http.NotFound(w, r)
		return
	}

	fmt.Fprintf(w, "this is the habit id %d\n", habitID)
}

func (hh *HabitHandler) HandleCreateHabit(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "create a habit\n")
}

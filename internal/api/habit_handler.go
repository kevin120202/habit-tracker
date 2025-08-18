package api

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/kevin120202/habit-tracker/internal/store"
	"github.com/kevin120202/habit-tracker/internal/utils"
)

type HabitHandler struct {
	habitStore store.HabitStore
	logger     *log.Logger
}

func NewHabitHandler(habitStore store.HabitStore, logger *log.Logger) *HabitHandler {
	return &HabitHandler{
		habitStore: habitStore,
		logger:     logger,
	}
}

func (hh *HabitHandler) HandleGetHabitByID(w http.ResponseWriter, r *http.Request) {
	habitID, err := utils.ReadIDParam(r)
	if err != nil {
		hh.logger.Printf("ERROR: readIDParam: %v", err)
		return
	}

	habit, err := hh.habitStore.GetHabitByID(habitID)
	if err != nil {
		hh.logger.Printf("ERROR: getHabitByID: %v", err)
		utils.WriteJSON(w, http.StatusInternalServerError, utils.Envelope{"error": "internal server error"})
		return
	}

	utils.WriteJSON(w, http.StatusOK, utils.Envelope{"habit": habit})
}

func (hh *HabitHandler) HandleCreateHabit(w http.ResponseWriter, r *http.Request) {
	var habit store.Habit

	err := json.NewDecoder(r.Body).Decode(&habit)
	if err != nil {
		hh.logger.Printf("ERROR: decodingCreateWorkout: %v", err)
		utils.WriteJSON(w, http.StatusBadRequest, utils.Envelope{"error": "invalid request sent"})
		return
	}

	createdHabit, err := hh.habitStore.CreateHabit(&habit)
	if err != nil {
		hh.logger.Printf("ERROR: createHabit: %v", err)
		utils.WriteJSON(w, http.StatusInternalServerError, utils.Envelope{"error": "failed to create habit"})
		return
	}

	utils.WriteJSON(w, http.StatusCreated, utils.Envelope{"habit": createdHabit})
}

func (hh *HabitHandler) HandleGetHabits(w http.ResponseWriter, r *http.Request) {
	habits, err := hh.habitStore.GetHabits()
	if err != nil {
		hh.logger.Printf("ERROR: getHabits: %v", err)
		utils.WriteJSON(w, http.StatusInternalServerError, utils.Envelope{"error": "failed to retrieve habits"})
		return
	}

	utils.WriteJSON(w, http.StatusOK, utils.Envelope{"habits": habits})
}

func (hh *HabitHandler) HandleUpdateHabitByID(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "update a habit\n")
}

func (hh *HabitHandler) HandleDeleteHabitByID(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "delete a habit\n")
}

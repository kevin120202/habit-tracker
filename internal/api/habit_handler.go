package api

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"

	"github.com/google/uuid"
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
	habitID, err := utils.ReadIDParam(r)
	if err != nil {
		hh.logger.Printf("ERROR: readIDParam: %v", err)
		utils.WriteJSON(w, http.StatusBadRequest, utils.Envelope{"error": "invalid habit id"})
		return
	}

	existingHabit, err := hh.habitStore.GetHabitByID(habitID)
	if err != nil {
		hh.logger.Printf("ERROR: getHabitByID: %v", err)
		utils.WriteJSON(w, http.StatusInternalServerError, utils.Envelope{"error": "internal server error"})
		return
	}

	if existingHabit == nil {
		http.NotFound(w, r)
		return
	}

	var updateHabitRequest struct {
		Name        *string `json:"name"`
		Description *string `json:"description"`
		Frequency   *string `json:"frequency"`
		TargetCount *int    `json:"target_count"`
		IsActive    *bool   `json:"is_active"`
	}

	err = json.NewDecoder(r.Body).Decode(&updateHabitRequest)
	if err != nil {
		hh.logger.Printf("ERROR: decodingUpdateRequest: %v", err)
		utils.WriteJSON(w, http.StatusBadRequest, utils.Envelope{"error": "invalid request payload"})
		return
	}

	if updateHabitRequest.Name != nil {
		existingHabit.Name = *updateHabitRequest.Name
	}
	if updateHabitRequest.Description != nil {
		existingHabit.Description = *updateHabitRequest.Description
	}
	if updateHabitRequest.Frequency != nil {
		existingHabit.Frequency = *updateHabitRequest.Frequency
	}
	if updateHabitRequest.TargetCount != nil {
		existingHabit.TargetCount = *updateHabitRequest.TargetCount
	}
	if updateHabitRequest.IsActive != nil {
		existingHabit.IsActive = *updateHabitRequest.IsActive
	}

	err = hh.habitStore.UpdateHabit(existingHabit)
	if err != nil {
		hh.logger.Printf("ERROR: updatingHabit: %v", err)
		utils.WriteJSON(w, http.StatusInternalServerError, utils.Envelope{"error": "internal server error"})
		return
	}

	utils.WriteJSON(w, http.StatusOK, utils.Envelope{"habit": existingHabit})
}

func (hh *HabitHandler) HandleDeleteHabitByID(w http.ResponseWriter, r *http.Request) {
	habitID, err := utils.ReadIDParam(r)
	if err != nil {
		hh.logger.Printf("ERROR: readIDParam: %v", err)
		return
	}

	err = hh.habitStore.DeleteHabit(habitID)
	if err == sql.ErrNoRows {
		utils.WriteJSON(w, http.StatusNotFound, utils.Envelope{"message": "habit not found"})
		return
	}

	if err != nil {
		utils.WriteJSON(w, http.StatusInternalServerError, utils.Envelope{"message": "error deleting habit"})
		return
	}

	utils.WriteJSON(w, http.StatusOK, utils.Envelope{"message": "habit deleted successfully"})
}

func (hh *HabitHandler) HandleLogHabitCompletions(w http.ResponseWriter, r *http.Request) {
	habitID, err := utils.ReadIDParam(r)
	if err != nil {
		hh.logger.Printf("ERROR: readIDParam: %v", err)
		return
	}

	var habitEntry store.HabitEntry
	habitEntry.HabitID = habitID

	err = json.NewDecoder(r.Body).Decode(&habitEntry)
	if err != nil {
		hh.logger.Printf("ERROR: decodingCreateHabitEntry: %v", err)
		utils.WriteJSON(w, http.StatusBadRequest, utils.Envelope{"error": "invalid request sent"})
		return
	}

	// Verify habit belongs to user

	createdHabitEntry, err := hh.habitStore.LogHabit(&habitEntry)
	if err != nil {
		hh.logger.Printf("ERROR: logHabit: %v", err)
		utils.WriteJSON(w, http.StatusInternalServerError, utils.Envelope{"error": "failed to create habit entry"})
		return
	}

	utils.WriteJSON(w, http.StatusCreated, utils.Envelope{"habitEntry": createdHabitEntry})
}

func (hh *HabitHandler) HandleCompleteHabit(w http.ResponseWriter, r *http.Request) {
	habitID, err := utils.ReadIDParam(r)
	if err != nil {
		hh.logger.Printf("ERROR: readIDParam: %v", err)
		return
	}

	existingHabit, err := hh.habitStore.GetHabitByID(habitID)
	if err != nil {
		hh.logger.Printf("ERROR: getHabitByID: %v", err)
		utils.WriteJSON(w, http.StatusInternalServerError, utils.Envelope{"error": "internal server error"})
		return
	}

	if existingHabit == nil {
		utils.WriteJSON(w, http.StatusInternalServerError, utils.Envelope{"error": "habit not found"})
		return
	}

	if !existingHabit.IsActive {
		utils.WriteJSON(w, http.StatusBadRequest, utils.Envelope{"error": "cannot complete an inactive habit"})
		return
	}

	var completedHabitEntry store.HabitEntry
	completedHabitEntry.HabitID = habitID

	err = json.NewDecoder(r.Body).Decode(&completedHabitEntry)
	if err != nil {
		hh.logger.Printf("ERROR: decodingCreateHabitEntry: %v", err)
		utils.WriteJSON(w, http.StatusBadRequest, utils.Envelope{"error": "invalid request sent"})
		return
	}

	// Verify habit belongs to user

	createdCompletedHabitEntry, err := hh.habitStore.LogHabit(&completedHabitEntry)
	if err != nil {
		hh.logger.Printf("ERROR: logHabit: %v", err)
		utils.WriteJSON(w, http.StatusInternalServerError, utils.Envelope{"error": "failed to create complete habit entry"})
		return
	}

	utils.WriteJSON(w, http.StatusCreated, utils.Envelope{"completedHabitEntry": createdCompletedHabitEntry, "message": "Habit completed successfully"})
}

func (hh *HabitHandler) HandleCreateTagToHabit(w http.ResponseWriter, r *http.Request) {
	habitID, err := utils.ReadIDParam(r)
	if err != nil {
		hh.logger.Printf("ERROR: readHabitIDParam: %v", err)
		utils.WriteJSON(w, http.StatusBadRequest, utils.Envelope{"error": "invalid habit id"})
		return
	}

	var requestBody struct {
		TagID uuid.UUID `json:"tag_id"`
	}

	err = json.NewDecoder(r.Body).Decode(&requestBody)
	if err != nil {
		hh.logger.Printf("ERROR: decodingAddTagRequest: %v", err)
		utils.WriteJSON(w, http.StatusBadRequest, utils.Envelope{"error": "invalid request payload"})
		return
	}

	// Verify habit exists
	existingHabit, err := hh.habitStore.GetHabitByID(habitID)
	if err != nil {
		hh.logger.Printf("ERROR: getHabitByID: %v", err)
		utils.WriteJSON(w, http.StatusInternalServerError, utils.Envelope{"error": "internal server error"})
		return
	}

	if existingHabit == nil {
		utils.WriteJSON(w, http.StatusNotFound, utils.Envelope{"error": "habit not found"})
		return
	}

	err = hh.habitStore.AddTagToHabit(habitID, requestBody.TagID)
	if err != nil {
		hh.logger.Printf("ERROR: addTagToHabit: %v", err)
		utils.WriteJSON(w, http.StatusInternalServerError, utils.Envelope{"error": "failed to add tag to habit"})
		return
	}

	utils.WriteJSON(w, http.StatusCreated, utils.Envelope{"message": "tag added to habit successfully"})
}

func (hh *HabitHandler) HandleDeleteTagFromHabit(w http.ResponseWriter, r *http.Request) {
	habitID, err := utils.ReadIDParam(r)
	if err != nil {
		hh.logger.Printf("ERROR: readHabitIDParam: %v", err)
		utils.WriteJSON(w, http.StatusBadRequest, utils.Envelope{"error": "invalid habit id"})
		return
	}

	tagID, err := utils.ReadTagIDParam(r)
	if err != nil {
		hh.logger.Printf("ERROR: readTagIDParam: %v", err)
		utils.WriteJSON(w, http.StatusBadRequest, utils.Envelope{"error": "invalid tag id"})
		return
	}

	// Verify habit exists
	existingHabit, err := hh.habitStore.GetHabitByID(habitID)
	if err != nil {
		hh.logger.Printf("ERROR: getHabitByID: %v", err)
		utils.WriteJSON(w, http.StatusInternalServerError, utils.Envelope{"error": "internal server error"})
		return
	}

	if existingHabit == nil {
		utils.WriteJSON(w, http.StatusNotFound, utils.Envelope{"error": "habit not found"})
		return
	}

	err = hh.habitStore.RemoveTagFromHabit(habitID, tagID)
	if err == sql.ErrNoRows {
		utils.WriteJSON(w, http.StatusNotFound, utils.Envelope{"error": "tag not found on habit"})
		return
	}

	if err != nil {
		hh.logger.Printf("ERROR: removeTagFromHabit: %v", err)
		utils.WriteJSON(w, http.StatusInternalServerError, utils.Envelope{"error": "failed to remove tag from habit"})
		return
	}

	utils.WriteJSON(w, http.StatusOK, utils.Envelope{"message": "tag removed from habit successfully"})
}

func (hh *HabitHandler) HandleGetHabitsByTag(w http.ResponseWriter, r *http.Request) {
	tagID, err := utils.ReadIDParam(r)
	if err != nil {
		hh.logger.Printf("ERROR: readTagIDParam: %v", err)
		utils.WriteJSON(w, http.StatusBadRequest, utils.Envelope{"error": "invalid tag id"})
		return
	}

	habits, err := hh.habitStore.GetHabitsByTag(tagID)
	if err != nil {
		hh.logger.Printf("ERROR: getHabitsByTag: %v", err)
		utils.WriteJSON(w, http.StatusInternalServerError, utils.Envelope{"error": "failed to retrieve habits by tag"})
		return
	}

	utils.WriteJSON(w, http.StatusOK, utils.Envelope{"habits": habits})
}

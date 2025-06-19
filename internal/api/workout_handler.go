package api

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"

	"github.com/melkeydev/femProject/internal/store"
	"github.com/melkeydev/femProject/internal/utils"
)

type WorkoutHandler struct {
	workoutStore store.WorkoutStore
	logger *log.Logger
}

func NewWorkoutHandler(ws store.WorkoutStore, logger *log.Logger) *WorkoutHandler {
	return &WorkoutHandler{
		workoutStore: ws,
		logger: logger,
	}
}

func (wh *WorkoutHandler) HandleGetWorkoutByID(w http.ResponseWriter, r *http.Request) {
	workoutID, err := utils.ReadIDParam(r)
	if err != nil {
		wh.logger.Printf("ERROR: readIDParams %v\n", err)
		utils.WriteJSON(w, http.StatusBadRequest, utils.Envelope{
			"error": "invalid workout id",
		})
		return
	}

	workout, err := wh.workoutStore.GetWorkoutByID(workoutID)
	if err != nil {
		wh.logger.Printf("ERROR: GetWorkoutByID %v\n", err)
		utils.WriteJSON(w, http.StatusNotFound, utils.Envelope{
			"error": "Not found",
		})
		return
	}

	utils.WriteJSON(w, http.StatusOK, utils.Envelope{
		"data": workout,
	})
}

func (wh *WorkoutHandler) HandleCreateWorkout(w http.ResponseWriter, r *http.Request) {
	var workout store.Workout

	err := json.NewDecoder(r.Body).Decode(&workout)
	if err != nil {
		wh.logger.Printf("ERROR: decodeJSON %v\n", err)
		utils.WriteJSON(w, http.StatusBadRequest, utils.Envelope{
			"error": err.Error(),
		})
		return
	}

	createdWorkout, err := wh.workoutStore.CreateWorkout(&workout)
	if err != nil {
		wh.logger.Printf("ERROR: CreateWorkout %v\n", err)
		utils.WriteJSON(w, http.StatusInternalServerError, utils.Envelope{
			"error": err.Error(),
		})
		return
	}

	utils.WriteJSON(w, http.StatusOK, utils.Envelope{"data": createdWorkout})
}

func (wh *WorkoutHandler) HandleUpdateWorkoutByID(w http.ResponseWriter, r *http.Request) {
	workoutID, err := utils.ReadIDParam(r)
	if err != nil {
		wh.logger.Printf("ERROR: ReadIDParam %v\n", err)
		utils.WriteJSON(w, http.StatusBadRequest, utils.Envelope{"error": err.Error()})
		return
	}

	workout, err := wh.workoutStore.GetWorkoutByID(workoutID)
	if err != nil {
		wh.logger.Printf("ERROR: GetWorkoutByID %v\n", err)
		utils.WriteJSON(w, http.StatusBadRequest, utils.Envelope{
			"error": err.Error(),
		})
		return
	}

	if workout == nil {
		wh.logger.Printf("ERROR: WorkoutNotFound %v\n", err)
		utils.WriteJSON(w, http.StatusNotFound, utils.Envelope{"error": "Workout Not Found"})
		return
	}

	var updateWorkout struct {
		ID              *int                 `json:"id"`
		Title           *string              `json:"title"`
		Description     *string              `json:"description"`
		DurationMinutes *int                 `json:"duration_minutes"`
		CaloriesBurned  *int                 `json:"calories_burned"`
		Entries         []store.WorkoutEntry `json:"entries"`
	}

	err = json.NewDecoder(r.Body).Decode(&updateWorkout)
	if err != nil {
		wh.logger.Printf("ERROR: Workout Body Decode Err %v\n", err)
		utils.WriteJSON(w, http.StatusBadRequest, utils.Envelope{"error": "Error Bad Body"})
		return
	}

	if updateWorkout.Title != nil {
		workout.Title = *updateWorkout.Title
	}
	if updateWorkout.Description != nil {
		workout.Description = *updateWorkout.Description
	}
	if updateWorkout.DurationMinutes != nil {
		workout.DurationMinutes = *updateWorkout.DurationMinutes
	}
	if updateWorkout.CaloriesBurned != nil {
		workout.CaloriesBurned = *updateWorkout.CaloriesBurned
	}
	if updateWorkout.Entries != nil {
		workout.Entries = updateWorkout.Entries
	}
	err = wh.workoutStore.UpdateWorkout(workout)
	if err != nil {
		wh.logger.Printf("ERROR: Workout Update %v\n", err)
		utils.WriteJSON(w, http.StatusInternalServerError, utils.Envelope{"error": err.Error()})
		return
	}
	utils.WriteJSON(w, http.StatusOK, utils.Envelope{"data": workout})
}

func (wh *WorkoutHandler) DeleteWorkoutByID(w http.ResponseWriter, r *http.Request) {
	workoutID, err := utils.ReadIDParam(r)
	if err != nil {
		wh.logger.Printf("ERROR: Read ID Params %w\n", err)
		utils.WriteJSON(w, http.StatusBadRequest, utils.Envelope{"error": err.Error()})
		return
	}

	err = wh.workoutStore.DeleteWorkout(workoutID)

	if err == sql.ErrNoRows {
		wh.logger.Printf("ERROR: Delete Workout %v\n", err)
		utils.WriteJSON(w, http.StatusNotFound, utils.Envelope{"error": "workour not found"})
		return
	}

	if err != nil {
		wh.logger.Printf("ERROR: DeleteWorkout %v\n", err)
		utils.WriteJSON(w, http.StatusInternalServerError, utils.Envelope{"error": err.Error()})
		return
	}

	utils.WriteJSON(w, http.StatusOK, utils.Envelope{"data": "deleted successfully"})
}

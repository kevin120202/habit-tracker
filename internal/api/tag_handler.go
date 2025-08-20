package api

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"

	"github.com/kevin120202/habit-tracker/internal/store"
	"github.com/kevin120202/habit-tracker/internal/utils"
)

type TagHandler struct {
	tagStore store.TagStore
	logger   *log.Logger
}

func NewTagHandler(tagStore store.TagStore, logger *log.Logger) *TagHandler {
	return &TagHandler{
		tagStore: tagStore,
		logger:   logger,
	}
}

func (th *TagHandler) HandleCreateTag(w http.ResponseWriter, r *http.Request) {
	var tag store.Tag

	err := json.NewDecoder(r.Body).Decode(&tag)
	if err != nil {
		th.logger.Printf("ERROR: decodingCreateTag: %v", err)
		utils.WriteJSON(w, http.StatusBadRequest, utils.Envelope{"error": "invalid request sent"})
		return
	}

	createdTag, err := th.tagStore.CreateTag(&tag)
	if err != nil {
		th.logger.Printf("ERROR: createTag: %v", err)

		if err.Error() == "tag with this name already exists" {
			utils.WriteJSON(w, http.StatusConflict, utils.Envelope{"error": "tag with this name already exists"})
			return
		}

		utils.WriteJSON(w, http.StatusInternalServerError, utils.Envelope{"error": "failed to create tag"})
		return
	}

	utils.WriteJSON(w, http.StatusCreated, utils.Envelope{"tag": createdTag})
}

func (th *TagHandler) HandleGetTagByID(w http.ResponseWriter, r *http.Request) {
	tagID, err := utils.ReadIDParam(r)
	if err != nil {
		th.logger.Printf("ERROR: readIDParam: %v", err)
		return
	}

	tag, err := th.tagStore.GetTagByID(tagID)
	if err != nil {
		th.logger.Printf("ERROR: getTagByID: %v", err)
		utils.WriteJSON(w, http.StatusInternalServerError, utils.Envelope{"error": "internal server error"})
		return
	}

	utils.WriteJSON(w, http.StatusOK, utils.Envelope{"tag": tag})
}

func (th *TagHandler) HandleGetTags(w http.ResponseWriter, r *http.Request) {
	tags, err := th.tagStore.GetTags()
	if err != nil {
		th.logger.Printf("ERROR: getTags: %v", err)
		utils.WriteJSON(w, http.StatusInternalServerError, utils.Envelope{"error": "failed to retrieve tags"})
		return
	}

	utils.WriteJSON(w, http.StatusOK, utils.Envelope{"tags": tags})
}

func (th *TagHandler) HandleUpdateTagByID(w http.ResponseWriter, r *http.Request) {
	tagID, err := utils.ReadIDParam(r)
	if err != nil {
		th.logger.Printf("ERROR: readIDParam: %v", err)
		utils.WriteJSON(w, http.StatusBadRequest, utils.Envelope{"error": "invalid tag id"})
		return
	}

	existingTag, err := th.tagStore.GetTagByID(tagID)
	if err != nil {
		th.logger.Printf("ERROR: getTagByID: %v", err)
		utils.WriteJSON(w, http.StatusInternalServerError, utils.Envelope{"error": "internal server error"})
		return
	}

	if existingTag == nil {
		utils.WriteJSON(w, http.StatusNotFound, utils.Envelope{"error": "tag not found"})
		return
	}

	var updateTagRequest struct {
		Name  *string `json:"name"`
		Color *string `json:"color"`
	}

	err = json.NewDecoder(r.Body).Decode(&updateTagRequest)
	if err != nil {
		th.logger.Printf("ERROR: decodingUpdateRequest: %v", err)
		utils.WriteJSON(w, http.StatusBadRequest, utils.Envelope{"error": "invalid request payload"})
		return
	}

	if updateTagRequest.Name != nil {
		existingTag.Name = *updateTagRequest.Name
	}
	if updateTagRequest.Color != nil {
		existingTag.Color = *updateTagRequest.Color
	}

	err = th.tagStore.UpdateTag(existingTag)
	if err != nil {
		th.logger.Printf("ERROR: updatingTag: %v", err)
		utils.WriteJSON(w, http.StatusInternalServerError, utils.Envelope{"error": "internal server error"})
		return
	}

	utils.WriteJSON(w, http.StatusOK, utils.Envelope{"tag": existingTag})
}

func (th *TagHandler) HandleDeleteTagByID(w http.ResponseWriter, r *http.Request) {
	tagID, err := utils.ReadIDParam(r)
	if err != nil {
		th.logger.Printf("ERROR: readIDParam: %v", err)
		return
	}

	err = th.tagStore.DeleteTag(tagID)
	if err == sql.ErrNoRows {
		utils.WriteJSON(w, http.StatusNotFound, utils.Envelope{"message": "tag not found"})
		return
	}

	if err != nil {
		utils.WriteJSON(w, http.StatusInternalServerError, utils.Envelope{"message": "error deleting tag"})
		return
	}

	utils.WriteJSON(w, http.StatusOK, utils.Envelope{"message": "tag deleted successfully"})
}

package api

import (
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

func (th *TagHandler) HandleGetTags(w http.ResponseWriter, r *http.Request) {
	tags, err := th.tagStore.GetTags()
	if err != nil {
		th.logger.Printf("ERROR: getTags: %v", err)
		utils.WriteJSON(w, http.StatusInternalServerError, utils.Envelope{"error": "failed to retrieve tags"})
		return
	}

	utils.WriteJSON(w, http.StatusOK, utils.Envelope{"tags": tags})
}

package topics

import (
	"net/http"

	"github.com/rahulsoibam/koubru-prod-api/utils"
)

// List all topics
func (a App) List(w http.ResponseWriter, r *http.Request) {
	result, err := dbList(a.DB)
	if err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	utils.RespondWithJSON(w, http.StatusOK, &result)
}

// Create a topic
func (a App) Create(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Create a topic"))
}

// Get details of a topic
func (a App) Get(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	topicID := ctx.Value("topic_id").(int64)
	result, err := dbGet(a.DB, topicID)
	if err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	utils.RespondWithJSON(w, http.StatusOK, &result)
}

// Patch a topic
func (a App) Patch(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Update a topic"))
}

// Delete a topic
func (a App) Delete(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Anonymize a topic"))
}

// Followers of a topic
func (a App) Followers(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	topicID := ctx.Value("topic_id").(int64)
	utils.RespondWithJSON(w, http.StatusOK, topicID)
}

// Follow a topic
func (a App) Follow(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Follow a topic"))
}

// Unfollow a topic
func (a App) Unfollow(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Unfollow a topic"))
}

// Report a topic
func (a App) Report(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Report a topic"))
}

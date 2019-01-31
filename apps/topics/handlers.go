package topics

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/go-chi/chi"
	"github.com/lib/pq"

	"github.com/rahulsoibam/koubru-prod-api/middleware"
	"github.com/rahulsoibam/koubru-prod-api/utils"
)

// List all topics
func (a *App) List(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	userID, ok := ctx.Value(middleware.UserCtxKeys(0)).(int64)
	perPage := r.FormValue("per_page")
	page := r.FormValue("page")
	sort := r.FormValue("sort")
	order := r.FormValue("order")

	limit, err := strconv.Atoi(perPage)
	if err != nil || limit <= 0 {
		limit = 30
	}

	var offset = 0
	pg, err := strconv.Atoi(page)
	if err != nil || pg <= 1 {
		offset = 0
	} else {
		offset = (pg - 1) * limit
	}

	var orderBy string
	switch sort {
	case "",
		"created":
		orderBy = "created_on"
	default:
		utils.RespondWithError(w, http.StatusBadRequest, "sort value invalid")
		return
	}

	switch order {
	case "":
		order = "desc"
	case "asc":
	case "desc":
	default:
		utils.RespondWithError(w, http.StatusBadRequest, "order value invalid")
		return
	}

	var topics *[]Topic
	// Optional authentication
	if ok {
		topics, err = a.dbAuthenticatedListTopics(userID, limit, offset, orderBy, order)
	} else {
		topics, err = a.dbListTopics(limit, offset, orderBy, order)
	}
	if err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	utils.RespondWithJSON(w, http.StatusOK, &topics)
}

// Create a topic
func (a *App) Create(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	userID, ok := ctx.Value(middleware.UserCtxKeys(0)).(int64)
	if !ok {
		utils.RespondWithError(w, http.StatusUnauthorized, "Problem with the user id associated with this token")
		return
	}
	var nt *NewTopic
	err := json.NewDecoder(r.Body).Decode(&nt)
	if err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, err.Error())
		return
	}
	defer r.Body.Close()
	// Validate the topic
	if err := nt.Validate(); err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	topic, err := a.dbCreateTopic(nt, userID)
	if err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	utils.RespondWithJSON(w, http.StatusOK, &topic)
}

// Get details of a topic
func (a *App) Get(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	topicID, err := strconv.ParseInt(chi.URLParam(r, "topic_id"), 10, 64)
	if err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	userID, ok := ctx.Value(middleware.UserCtxKeys(0)).(int64)
	var topic *Topic
	if ok {
		topic, err = a.dbAuthenticatedGetTopicByID(userID, topicID)
	} else {
		topic, err = a.dbGetTopicByID(topicID)
	}
	if err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	utils.RespondWithJSON(w, http.StatusOK, &topic)
}

// Patch a topic
func (a *App) Patch(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Update a topic"))
}

// Delete a topic
func (a *App) Delete(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Anonymize a topic"))
}

// Followers of a topic
func (a *App) Followers(w http.ResponseWriter, r *http.Request) {

}

// Follow a topic
func (a *App) Follow(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	id := chi.URLParam(r, "topic_id")
	followerID, ok := ctx.Value(middleware.UserCtxKeys(0)).(int64)
	if !ok {
		utils.RespondWithError(w, http.StatusUnauthorized, "Invalid user. Please try authenticating again")
		return
	}
	_, err := a.DB.Exec("INSERT INTO Topic_Follower (topic_id, user_id) VALUES ($1, $2)", id, followerID)
	if err != nil {
		if e, ok := err.(*pq.Error); ok {
			if e.Code == "23505" {
				utils.RespondWithError(w, http.StatusBadRequest, "You are already following this topic")
				return
			}
			utils.RespondWithError(w, http.StatusInternalServerError, e.Detail)
			return
		}
		utils.RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	utils.RespondWithMessage(w, http.StatusOK, "Successfully followed topic")
}

// Unfollow a topic
func (a *App) Unfollow(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	id := chi.URLParam(r, "topic_id")
	followerID, ok := ctx.Value(middleware.UserCtxKeys(0)).(int64)
	if !ok {
		utils.RespondWithError(w, http.StatusUnauthorized, "Invalid user. Please try authticating again")
		return
	}
	response, err := a.DB.Exec("DELETE FROM Topic_Follower WHERE topic_id=$1 AND user_id=$2", id, followerID)
	if err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	count, err := response.RowsAffected()
	if err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	if count == 0 {
		utils.RespondWithError(w, http.StatusBadRequest, "You do not follow this topic")
		return
	}
	utils.RespondWithMessage(w, http.StatusOK, "Topic unfollowed")
}

// Report a topic
func (a *App) Report(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Report a topic"))
}

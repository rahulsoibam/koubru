package topics

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/lib/pq"

	"github.com/rahulsoibam/koubru/errs"
	"github.com/rahulsoibam/koubru/middleware"
	"github.com/rahulsoibam/koubru/types"
	"github.com/rahulsoibam/koubru/utils"
)

// List all topics
func (a *App) List(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	topics := []types.SearchTopic{}
	var err error
	topics, err = a.ListQuery(ctx)

	if err != nil {
		log.Println(err)
		utils.RespondWithError(w, http.StatusInternalServerError, errs.InternalServerError)
		return
	}
	utils.RespondWithJSON(w, http.StatusOK, topics)
}

// Create a topic
func (a *App) Create(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	userID, auth := ctx.Value(middleware.AuthKeys("user_id")).(int64)
	if !auth {
		a.Log.Infoln(errs.Unauthorized)
		utils.RespondWithError(w, http.StatusUnauthorized, errs.Unauthorized)
		return
	}
	var t types.NewTopic
	err := json.NewDecoder(r.Body).Decode(&t)
	if err != nil {
		a.Log.Infoln(err, r.Body)
		utils.RespondWithError(w, http.StatusBadRequest, errs.BadRequest)
		return
	}
	defer r.Body.Close()
	// Validate the topic
	if err := t.Validate(); err != nil {
		a.Log.Infoln(err)
		utils.RespondWithError(w, http.StatusBadRequest, errs.BadRequest)
		return
	}

	topic, err := a.AuthCreateQuery(userID, t)
	if err != nil {
		log.Println(err)
		utils.RespondWithError(w, http.StatusInternalServerError, errs.InternalServerError)
		return
	}
	utils.RespondWithJSON(w, http.StatusOK, topic)
}

// Get details of a topic
func (a *App) Get(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	userID, auth := ctx.Value(middleware.AuthKeys("user_id")).(int64)
	topicID := ctx.Value(middleware.TopicKeys("topic_id")).(int64)

	topic := types.Topic{}
	var err error
	if auth {
		topic, err = a.AuthGetQuery(userID, topicID)
	} else {
		topic, err = a.GetQuery(topicID)
	}
	if err != nil {
		log.Println(err)
		utils.RespondWithError(w, http.StatusInternalServerError, errs.InternalServerError)
		return
	}
	utils.RespondWithJSON(w, http.StatusOK, topic)
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
	ctx := r.Context()
	userID, auth := ctx.Value(middleware.AuthKeys("user_id")).(int64)
	topicID := ctx.Value(middleware.TopicKeys("topic_id")).(int64)

	followers := []types.User_{}
	var err error
	if auth {
		followers, err = a.AuthFollowersQuery(userID, topicID)
	} else {
		followers, err = a.FollowersQuery(topicID)
	}
	if err != nil {
		log.Println(err)
		utils.RespondWithError(w, http.StatusInternalServerError, errs.InternalServerError)
		return
	}
	utils.RespondWithJSON(w, http.StatusOK, followers)
}

// Follow a topic
func (a *App) Follow(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	topicID := ctx.Value(middleware.TopicKeys("topic_id")).(int64)
	followerID, auth := ctx.Value(middleware.AuthKeys("user_id")).(int64)

	if !auth {
		log.Println(errs.UnintendedExecution)
		utils.RespondWithError(w, http.StatusUnauthorized, errs.Unauthorized)
		return
	}
	_, err := a.DB.Exec("INSERT INTO Topic_Follower (topic_id, follower_id) VALUES ($1, $2)", topicID, followerID)
	if err != nil {
		if e, ok := err.(*pq.Error); ok {
			if e.Code == "23505" {
				a.Log.Infoln(err)
				utils.RespondWithError(w, http.StatusBadRequest, errs.TopicFollowAlreadyFollowing)
				return
			}
		}
		log.Println(err)
		utils.RespondWithError(w, http.StatusInternalServerError, errs.InternalServerError)
		return
	}
	utils.RespondWithMessage(w, http.StatusOK, "Followed")
}

// Unfollow a topic
func (a *App) Unfollow(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	topicID := ctx.Value(middleware.TopicKeys("topic_id")).(int64)
	followerID, auth := ctx.Value(middleware.AuthKeys("user_id")).(int64)
	if !auth {
		log.Println(errs.UnintendedExecution)
		utils.RespondWithError(w, http.StatusUnauthorized, errs.Unauthorized)
		return
	}
	response, err := a.DB.Exec("DELETE FROM Topic_Follower WHERE topic_id=$1 AND follower_id=$2", topicID, followerID)
	if err != nil {
		log.Println(err)
		utils.RespondWithError(w, http.StatusInternalServerError, errs.InternalServerError)
		return
	}
	count, err := response.RowsAffected()
	if err != nil {
		log.Println(err)
		utils.RespondWithError(w, http.StatusInternalServerError, errs.InternalServerError)
		return
	}
	if count == 0 {
		utils.RespondWithError(w, http.StatusBadRequest, errs.CategoryUnfollowNotFollowing)
		return
	}
	utils.RespondWithMessage(w, http.StatusOK, "Unfollowed")
}

// Report a topic
func (a *App) Report(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Report a topic"))
}

func (a *App) Opinions(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	topicID := ctx.Value(middleware.TopicKeys("topic_id")).(int64)
	userID, auth := ctx.Value(middleware.AuthKeys("user_id")).(int64)

	opinions := []types.Opinion{}
	var err error
	if auth {
		opinions, err = a.AuthOpinionsQuery(userID, topicID)
	} else {
		opinions, err = a.OpinionsQuery(topicID)
	}
	if err != nil {
		log.Println(err)
		utils.RespondWithError(w, http.StatusInternalServerError, errs.InternalServerError)
		return
	}

	utils.RespondWithJSON(w, http.StatusOK, opinions)
}

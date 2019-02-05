package users

import (
	"database/sql"
	"net/http"

	"github.com/lib/pq"
	"github.com/rahulsoibam/koubru-prod-api/errs"
	"github.com/rahulsoibam/koubru-prod-api/types"

	"github.com/rahulsoibam/koubru-prod-api/middleware"
	"github.com/rahulsoibam/koubru-prod-api/utils"
)

// Get details of authenticated user
func (a *App) Get(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	userID, auth := ctx.Value(middleware.AuthKeys("user_id")).(int64)
	usernameID := ctx.Value(middleware.UsernameIDKeys("username_id")).(int64)
	// username := ctx.Value(middleware.UsernameIDKeys("username")).(string)

	user := types.User{}
	var err error
	if auth {
		user, err = a.AuthGetQuery(userID, usernameID)
	} else {
		user, err = a.GetQuery(usernameID)
	}
	if err != nil {
		if err == sql.ErrNoRows {
			a.Log.Infoln(err)
			utils.RespondWithError(w, http.StatusNotFound, errs.UserNotFound)
			return
		}
		a.Log.Errorln(err)
		utils.RespondWithError(w, http.StatusInternalServerError, errs.InternalServerError)
		return
	}
	utils.RespondWithJSON(w, http.StatusOK, user)
}

// Followers to list followers to authenticated user
func (a *App) Followers(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	userID, auth := ctx.Value(middleware.AuthKeys("user_id")).(int64)
	usernameID := ctx.Value(middleware.UsernameIDKeys("username_id")).(int64)
	// username := ctx.Value(middleware.UsernameIDKeys("username")).(string)

	followers := []types.Follower{}
	var err error
	if auth {
		followers, err = a.AuthFollowersQuery(userID, usernameID)
	} else {
		followers, err = a.FollowersQuery(usernameID)
	}

	if err != nil {
		a.Log.Errorln(err)
		utils.RespondWithError(w, http.StatusInternalServerError, errs.InternalServerError)
		return
	}
	utils.RespondWithJSON(w, http.StatusOK, followers)
}

// Following to list users whom the authenticated user is following
func (a *App) Following(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	userID, auth := ctx.Value(middleware.AuthKeys("user_id")).(int64)
	usernameID := ctx.Value(middleware.UsernameIDKeys("username_id")).(int64)
	// username := ctx.Value(middleware.UsernameIDKeys("username")).(string)

	following := []types.Following{}
	var err error
	if auth {
		following, err = a.AuthFollowingQuery(userID, usernameID)
	} else {
		following, err = a.FollowingQuery(usernameID)
	}

	if err != nil {
		a.Log.Errorln(err)
		utils.RespondWithError(w, http.StatusInternalServerError, errs.InternalServerError)
		return
	}
	utils.RespondWithJSON(w, http.StatusOK, following)
}

// Topics of authenticated user
func (a *App) Topics(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	userID, auth := ctx.Value(middleware.AuthKeys("user_id")).(int64)
	usernameID := ctx.Value(middleware.UsernameIDKeys("username_id")).(int64)
	// username := ctx.Value(middleware.UsernameIDKeys("username")).(string)

	topics := []types.Topic_{}
	var err error
	if auth {
		topics, err = a.AuthTopicsQuery(userID, usernameID)
	} else {
		topics, err = a.TopicsQuery(usernameID)
	}
	if err != nil {
		a.Log.Errorln(err)
		utils.RespondWithError(w, http.StatusInternalServerError, errs.InternalServerError)
		return
	}
	utils.RespondWithJSON(w, http.StatusOK, topics)
}

// Opinions of authenticated user
func (a *App) Opinions(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	userID, auth := ctx.Value(middleware.AuthKeys("user_id")).(int64)
	usernameID := ctx.Value(middleware.UsernameIDKeys("username_id")).(int64)
	// username := ctx.Value(middleware.UsernameIDKeys("username")).(string)

	opinions := []types.Opinion{}
	var err error
	if auth {
		opinions, err = a.AuthOpinionsQuery(userID, usernameID)
	} else {
		opinions, err = a.OpinionsQuery(usernameID)
	}
	if err != nil {
		a.Log.Errorln(err)
		utils.RespondWithError(w, http.StatusInternalServerError, errs.InternalServerError)
		return
	}
	utils.RespondWithJSON(w, http.StatusOK, opinions)
}

func (a *App) Follow(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	followerID, auth := ctx.Value(middleware.AuthKeys("user_id")).(int64)
	usernameID := ctx.Value(middleware.UsernameIDKeys("username_id")).(int64)
	username := ctx.Value(middleware.UsernameIDKeys("username")).(string)
	if !auth {
		a.Log.Errorln(ctx)
		utils.RespondWithError(w, http.StatusUnauthorized, errs.Unauthorized)
		return
	}

	_, err := a.DB.Exec("INSERT INTO User_Follower (user_id, follower_id) VALUES ($1, $2)", usernameID, followerID)
	if err != nil {
		if e, ok := err.(*pq.Error); ok {
			if e.Code == "23505" {
				a.Log.Infoln(e)
				utils.RespondWithError(w, http.StatusBadRequest, errs.UserFollowAlreadyFollowing)
				return
			}
		}
		a.Log.Errorln(err)
		utils.RespondWithError(w, http.StatusInternalServerError, errs.InternalServerError)
		return
	}
	utils.RespondWithMessage(w, http.StatusOK, "@"+username+" followed.")
}

func (a *App) Unfollow(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	followerID, auth := ctx.Value(middleware.AuthKeys("user_id")).(int64)
	usernameID := ctx.Value(middleware.UsernameIDKeys("username_id")).(int64)
	username := ctx.Value(middleware.UsernameIDKeys("username")).(string)
	if !auth {
		a.Log.Errorln(ctx)
		utils.RespondWithError(w, http.StatusUnauthorized, errs.Unauthorized)
		return
	}

	response, err := a.DB.Exec("DELETE FROM User_Follower WHERE user_id=$1 AND follower_id=$2", usernameID, followerID)
	if err != nil {
		a.Log.Errorln(err)
		utils.RespondWithError(w, http.StatusInternalServerError, errs.InternalServerError)
		return
	}
	count, err := response.RowsAffected()
	if err != nil {
		a.Log.Errorln(err)
		utils.RespondWithError(w, http.StatusInternalServerError, errs.InternalServerError)
		return
	}
	if count == 0 {
		utils.RespondWithError(w, http.StatusBadRequest, errs.UserUnfollowNotFollowing)
		return
	}
	utils.RespondWithMessage(w, http.StatusOK, "@"+username+" unfollowed.")
}

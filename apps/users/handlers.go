package users

import (
	"database/sql"
	"log"
	"net/http"

	"github.com/lib/pq"
	"github.com/rahulsoibam/koubru/errs"
	"github.com/rahulsoibam/koubru/types"

	"github.com/rahulsoibam/koubru/middleware"
	"github.com/rahulsoibam/koubru/utils"
)

func (a *App) List(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	users := []types.SearchUser{}
	var err error
	users, err = a.ListQuery(ctx)

	if err != nil {
		log.Println(err)
		utils.RespondWithError(w, http.StatusInternalServerError, errs.InternalServerError)
		return
	}
	utils.RespondWithJSON(w, http.StatusOK, users)
}

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
		log.Println(err)
		if err == sql.ErrNoRows {
			log.Println(err, errs.UnintendedExecution)
			utils.RespondWithError(w, http.StatusNotFound, errs.UserNotFound)
			return
		}
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

	followers := []types.UserForFollowList{}
	var err error
	if auth {
		followers, err = a.AuthFollowersQuery(userID, usernameID)
	} else {
		followers, err = a.FollowersQuery(usernameID)
	}

	if err != nil {
		log.Println(err)
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

	following := []types.UserForFollowList{}
	var err error
	if auth {
		following, err = a.AuthFollowingQuery(userID, usernameID)
	} else {
		following, err = a.FollowingQuery(usernameID)
	}

	if err != nil {
		log.Println(err)
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

	topics := []types.TopicForList{}
	var err error
	if auth {
		topics, err = a.AuthTopicsQuery(userID, usernameID)
	} else {
		topics, err = a.TopicsQuery(usernameID)
	}
	if err != nil {
		log.Println(err)
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
		log.Println(err)
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
		log.Println(ctx)
		utils.RespondWithError(w, http.StatusUnauthorized, errs.Unauthorized)
		return
	}

	_, err := a.DB.Exec("INSERT INTO User_Follower (user_id, follower_id) VALUES ($1, $2)", usernameID, followerID)
	if err != nil {
		log.Println(err)
		if e, ok := err.(*pq.Error); ok {
			if e.Code == "23505" {
				log.Println(e)
				utils.RespondWithError(w, http.StatusBadRequest, errs.UserFollowAlreadyFollowing)
				return
			}
		}
		log.Println(err)
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
		log.Println(ctx)
		utils.RespondWithError(w, http.StatusUnauthorized, errs.Unauthorized)
		return
	}

	response, err := a.DB.Exec("DELETE FROM User_Follower WHERE user_id=$1 AND follower_id=$2", usernameID, followerID)
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
		utils.RespondWithError(w, http.StatusBadRequest, errs.UserUnfollowNotFollowing)
		return
	}
	utils.RespondWithMessage(w, http.StatusOK, "@"+username+" unfollowed.")
}

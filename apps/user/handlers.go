package user

import (
	"database/sql"
	"log"
	"net/http"

	"github.com/rahulsoibam/koubru/errs"
	"github.com/rahulsoibam/koubru/types"

	"github.com/rahulsoibam/koubru/middleware"
	"github.com/rahulsoibam/koubru/utils"
)

// Get details of authenticated user
func (a *App) Get(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	userID, auth := ctx.Value(middleware.AuthKeys("user_id")).(int64)
	if !auth {
		log.Println(ctx, errs.UnintendedExecution)
		utils.RespondWithError(w, http.StatusUnauthorized, errs.Unauthorized)
		return
	}

	user := types.User{}
	var err error

	user, err = a.AuthGetQuery(userID)
	if err != nil {
		log.Println(err)
		if err == sql.ErrNoRows {
			utils.RespondWithError(w, http.StatusNotFound, errs.UserNotFound)
			return
		}
		utils.RespondWithError(w, http.StatusInternalServerError, err)
		return
	}
	utils.RespondWithJSON(w, http.StatusOK, user)
}

// Followers to list followers to authenticated user
func (a *App) Followers(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	userID, auth := ctx.Value(middleware.AuthKeys("user_id")).(int64)
	if !auth {
		log.Println(ctx)
		utils.RespondWithError(w, http.StatusUnauthorized, errs.Unauthorized)
		return
	}

	followers := []types.UserForFollowList{}
	var err error
	followers, err = a.AuthFollowersQuery(userID)
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
	if !auth {
		log.Println(ctx)
		utils.RespondWithError(w, http.StatusUnauthorized, errs.Unauthorized)
		return
	}

	following := []types.UserForFollowList{}
	var err error
	following, err = a.AuthFollowingQuery(userID)
	if err != nil {
		log.Println(err)
		utils.RespondWithError(w, http.StatusInternalServerError, errs.InternalServerError)
		return
	}
	utils.RespondWithJSON(w, http.StatusOK, following)
}

// Opinions of authenticated user
func (a *App) Opinions(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	userID, auth := ctx.Value(middleware.AuthKeys("user_id")).(int64)
	if !auth {
		log.Println(ctx)
		utils.RespondWithError(w, http.StatusUnauthorized, errs.Unauthorized)
		return
	}

	opinions := []types.Opinion{}
	var err error
	opinions, err = a.AuthOpinionsQuery(userID)
	if err != nil {
		log.Println(err)
		utils.RespondWithError(w, http.StatusInternalServerError, errs.InternalServerError)
		return
	}
	utils.RespondWithJSON(w, http.StatusOK, opinions)
}

// Topics of authenticated user
func (a *App) Topics(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	userID, auth := ctx.Value(middleware.AuthKeys("user_id")).(int64)
	if !auth {
		log.Println(ctx)
		utils.RespondWithError(w, http.StatusUnauthorized, errs.Unauthorized)
		return
	}

	topics := []types.TopicForList{}
	var err error
	topics, err = a.AuthTopicsQuery(userID)
	if err != nil {
		log.Println(err)
		utils.RespondWithError(w, http.StatusInternalServerError, errs.InternalServerError)
		return
	}
	utils.RespondWithJSON(w, http.StatusOK, topics)
}

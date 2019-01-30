package user

import (
	"database/sql"
	"errors"
	"net/http"

	"github.com/go-chi/chi"

	"github.com/rahulsoibam/koubru-prod-api/middleware"
	"github.com/rahulsoibam/koubru-prod-api/utils"
)

// Get details of authenticated user
func (a *App) Get(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	userID := ctx.Value(middleware.UserCtxKeys(0)).(int64)
	var err error
	user, err := a.dbGetUserByID(userID)
	if err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	utils.RespondWithJSON(w, http.StatusOK, &user)
}

// Patch details of authenticated user
func (a *App) Patch(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("patching in testing mode. Get ready to send multipart-form data"))
}

// Delete or deactivate authenticated user
func (a *App) Delete(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("need a definition of delete on our platform"))
}

// Followers to list followers to authenticated user
func (a *App) Followers(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	userID := ctx.Value(middleware.UserCtxKeys(0)).(int64)
	followers, err := a.dbGetFollowersByID(userID)
	if err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	// Construction of response
	// result := struct {
	// 	Followers *[]FollowUser `json:"followers"`
	// }{Followers: followers}

	utils.RespondWithJSON(w, http.StatusOK, &followers)
}

// Following to list users whom the authenticated user is following
func (a *App) Following(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	userID := ctx.Value(middleware.UserCtxKeys(0)).(int64)
	following, err := a.dbGetFollowingByID(userID)
	if err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	// Construction of response
	// result := struct {
	// 	Following *[]FollowUser `json:"following"`
	// }{Following: following}
	utils.RespondWithJSON(w, http.StatusOK, &following)
}

// Opinions of authenticated user
func (a *App) Opinions(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("list opinions of authenticated user"))
}

// Topics of authenticated user
func (a *App) Topics(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("list topics of authenticated user"))
}

// USERS /users/ endpoint functions

// UsersGet returns the details of a user
func (a *App) UsersGet(w http.ResponseWriter, r *http.Request) {
	username := chi.URLParam(r, "username")
	userID, err := a.validateUsernameAndGetID(username)
	if err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	user, err := a.dbGetUserByID(userID)
	if err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	utils.RespondWithJSON(w, http.StatusOK, &user)
}

func (a *App) UsersFollowers(w http.ResponseWriter, r *http.Request) {
	username := chi.URLParam(r, "username")
	userID, err := a.validateUsernameAndGetID(username)
	if err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	followers, err := a.dbGetFollowersByID(userID)
	if err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	utils.RespondWithJSON(w, http.StatusOK, &followers)

}

func (a *App) UsersFollowing(w http.ResponseWriter, r *http.Request) {
	username := chi.URLParam(r, "username")
	userID, err := a.validateUsernameAndGetID(username)
	if err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	following, err := a.dbGetFollowingByID(userID)
	if err != nil {
		utils.RespondWithJSON(w, http.StatusInternalServerError, err.Error())
		return
	}

	utils.RespondWithJSON(w, http.StatusOK, &following)
}

func (a *App) UsersTopics(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("List topics of a user"))
}

func (a *App) UsersOpinions(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("List opinions of a user"))
}

func (a *App) FollowUser(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	username := chi.URLParam(r, "username")
	followerID := ctx.Value(middleware.UserCtxKeys(0)).(int64)
	userID, err := a.validateUsernameAndGetID(username)
	if err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, err.Error())
		return
	}
	response, _ := a.DB.Exec("INSERT INTO UserMap (user_id, follower_id) VALUES ($1, $2)", userID, followerID)
	count, err := response.RowsAffected()
	if err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	if count == 0 {
		utils.RespondWithError(w, http.StatusBadRequest, "You have already followed this user")
		return
	}

	utils.RespondWithJSON(w, http.StatusOK, "You have followed "+username)
	// if err != nil {
	// 	if e, ok := err.(*pq.Error); ok {
	// 		if e.Code == "23505" {
	// 			utils.RespondWithError(w, http.StatusBadRequest, "You are already following this user")
	// 			return
	// 		} else {
	// 			utils.RespondWithError(w, http.StatusInternalServerError, e.Detail)
	// 			return
	// 		}
	// 	}
	// 	utils.RespondWithError(w, http.StatusInternalServerError, err.Error())
	// 	return
	// }
}

func (a *App) UnfollowUser(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	username := chi.URLParam(r, "username")
	followerID := ctx.Value(middleware.UserCtxKeys(0)).(int64)
	userID, err := a.validateUsernameAndGetID(username)
	if err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, err.Error())
		return
	}
	_, err = a.DB.Exec("DELETE FROM UserMap WHERE user_id = $1 AND follower_id = $2", userID, followerID)
	if err != nil {
		if err == sql.ErrNoRows {
			utils.RespondWithJSON(w, http.StatusBadRequest, "You do not follow this user")
			return
		}
		utils.RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	utils.RespondWithJSON(w, http.StatusOK, "User unfollowed")
}

func (a *App) validateUsernameAndGetID(username string) (int64, error) {
	if err := utils.ValidateUsername(username); err != nil || !utils.UsernameRegex.MatchString(username) {
		return 0, errors.New("username is invalid")
	}
	userID, err := a.dbGetUserIDUsingUsername(username)
	if err != nil {
		return 0, err
	}
	return userID, nil
}

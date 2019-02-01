package user

import (
	"errors"
	"log"
	"net/http"
	"strconv"

	"github.com/go-chi/chi"
	"github.com/lib/pq"

	"github.com/rahulsoibam/koubru-prod-api/middleware"
	"github.com/rahulsoibam/koubru-prod-api/utils"
)

// Get details of authenticated user
func (a *App) Get(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	userID := ctx.Value(middleware.UserCtxKeys(0)).(int64)
	var err error
	user, err := a.dbAuthenticatedGetUserSelf(userID)
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
	userID, ok := ctx.Value(middleware.UserCtxKeys(0)).(int64)
	if !ok {
		utils.RespondWithError(w, http.StatusUnauthorized, "You are unauthorized to perfrom this action")
		return
	}
	followers, err := a.dbAuthenticatedGetFollowersSelf(userID)
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
	following, err := a.dbAuthenticatedGetFollowing(userID)
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
	ctx := r.Context()
	userID, ok := ctx.Value(middleware.UserCtxKeys(0)).(int64)
	if !ok {
		utils.RespondWithError(w, http.StatusUnauthorized, "You are unauthorized to view this request")
	}
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
	topics, err := a.dbAuthenticatedListTopics(userID, userID, limit, offset, orderBy, order)
	if err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	utils.RespondWithJSON(w, http.StatusOK, &topics)
}

// USERS /users/ endpoint functions

// UsersGet returns the details of a user
func (a *App) UsersGet(w http.ResponseWriter, r *http.Request) {
	username := chi.URLParam(r, "username")
	ctx := r.Context()
	userID, ok := ctx.Value(middleware.UserCtxKeys(0)).(int64)
	quserID, err := a.validateUsernameAndGetID(username)
	if err != nil {
		panic(err)
		utils.RespondWithError(w, http.StatusBadRequest, err.Error())
		return
	}
	log.Println(quserID)
	var user *User
	if ok {
		user, err = a.dbAuthenticatedGetUser(userID, quserID)
	} else {
		user, err = a.dbGetUser(quserID)
	}
	if err != nil {
		panic(err)
		utils.RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	utils.RespondWithJSON(w, http.StatusOK, &user)
}

func (a *App) UsersFollowers(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	userID, ok := ctx.Value(middleware.UserCtxKeys(0)).(int64)
	username := chi.URLParam(r, "username")
	quserID, err := a.validateUsernameAndGetID(username)
	if err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	var followers *[]FollowUser

	if ok {
		followers, err = a.dbAuthenticatedGetFollowers(userID, quserID)
	} else {
		followers, err = a.dbGetFollowers(quserID)
	}
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
	ctx := r.Context()
	username := chi.URLParam(r, "username")
	quserID, err := a.validateUsernameAndGetID(username)
	userID, ok := ctx.Value(middleware.UserCtxKeys(0)).(int64)
	if err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, err.Error())
		return
	}
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
	if ok {
		topics, err = a.dbAuthenticatedListTopics(userID, quserID, limit, offset, orderBy, order)
	} else {
		topics, err = a.dbListTopics(quserID, limit, offset, orderBy, order)
	}
	if err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	utils.RespondWithJSON(w, http.StatusOK, &topics)
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
	_, err = a.DB.Exec("INSERT INTO UserMap (user_id, follower_id) VALUES ($1, $2)", userID, followerID)
	if err != nil {
		if e, ok := err.(*pq.Error); ok {
			if e.Code == "23505" {
				utils.RespondWithError(w, http.StatusBadRequest, "You are already following this user")
				return
			}

			utils.RespondWithError(w, http.StatusInternalServerError, e.Detail)
			return
		}
		utils.RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	utils.RespondWithMessage(w, http.StatusOK, "You have followed "+username)
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
	response, err := a.DB.Exec("DELETE FROM UserMap WHERE user_id = $1 AND follower_id = $2", userID, followerID)
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
		utils.RespondWithError(w, http.StatusBadRequest, "You do not follow this user")
		return
	}
	utils.RespondWithMessage(w, http.StatusOK, "User unfollowed")
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

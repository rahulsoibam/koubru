package user

import (
	"net/http"

	"github.com/rahulsoibam/koubru-prod-api/middleware"
	"github.com/rahulsoibam/koubru-prod-api/utils"
)

// Get details of authenticated user
func (a App) Get(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	userID := ctx.Value(middleware.UserCtxKeys(0)).(int64)
	var err error
	user, err := dbGetUserByID(a.DB, userID)
	if err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	utils.RespondWithJSON(w, http.StatusOK, &user)
}

// Patch details of authenticated user
func (a App) Patch(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("patching in testing mode. Get ready to send multipart-form data"))
}

// Delete or deactivate authenticated user
func (a App) Delete(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("need a definition of delete on our platform"))
}

// Followers to list followers to authenticated user
func (a App) Followers(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	userID := ctx.Value(middleware.UserCtxKeys(0)).(int64)
	followers, err := dbGetFollowersByID(a.DB, userID)
	if err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	utils.RespondWithJSON(w, http.StatusOK, &followers)
}

// Following to list users whom the authenticated user is following
func (a App) Following(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	userID := ctx.Value(middleware.UserCtxKeys(0)).(int64)
	following, err := dbGetFollowingByID(a.DB, userID)
	if err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	utils.RespondWithJSON(w, http.StatusOK, &following)
}

// Opinions of authenticated user
func (a App) Opinions(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("list opinions of authenticated user"))
}

// Topics of authenticated user
func (a App) Topics(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("list topics of authenticated user"))
}

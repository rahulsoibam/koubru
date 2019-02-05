package opinions

import (
	"net/http"

	"github.com/rahulsoibam/koubru/errs"
	"github.com/rahulsoibam/koubru/utils"

	"github.com/rahulsoibam/koubru/middleware"
	"github.com/rahulsoibam/koubru/types"
)

// List to list all opinions
func (a *App) List(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	userID, auth := ctx.Value(middleware.AuthKeys("user_id")).(int64)

	opinions := []types.Opinion{}
	var err error
	if auth {
		opinions, err = a.AuthListQuery(ctx, userID)
	} else {
		opinions, err = a.ListQuery(ctx)
	}
	if err != nil {
		a.Log.Errorln(err)
		utils.RespondWithError(w, http.StatusInternalServerError, errs.InternalServerError)
		return
	}
	utils.RespondWithJSON(w, http.StatusOK, opinions)
}

// Create to create an opinion
func (a *App) Create(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Create opinion"))
}

// Get to get details of an opinion
func (a *App) Get(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	userID, auth := ctx.Value(middleware.AuthKeys("user_id")).(int64)
	opinionID := ctx.Value(middleware.OpinionKeys("opinion_id")).(int64)

	opinion := types.Opinion{}
	var err error

	if auth {
		opinion, err = a.AuthGetQuery(userID, opinionID)
	} else {
		opinion, err = a.GetQuery(opinionID)
	}

	if err != nil {
		a.Log.Errorln(err)
		utils.RespondWithError(w, http.StatusInternalServerError, errs.InternalServerError)
		return
	}

	utils.RespondWithJSON(w, http.StatusOK, opinion)
}

// Delete to delete an opinion
func (a *App) Delete(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Delete opinion"))
}

// Followers to get followers of an opinion
func (a *App) Followers(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Get followers of opinion"))
}

// Follow to follow an opinion
func (a *App) Follow(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Follow opinion"))
}

// Unfollow to unfollow an opinion
func (a *App) Unfollow(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Unfollow opinion"))
}

// Replies to reply to an opinion
func (a *App) Replies(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Get replies of opinion"))
}

// Report to report an opinion
func (a *App) Report(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Report opinion"))
}

// Vote to vote on an opinion
func (a *App) Vote(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Vote on an opinion"))
}

func (a *App) Breadcrumbs(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Breadcrumbs of an opinion"))
}

func (a *App) Reply(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Direct replies of an opinion"))
}

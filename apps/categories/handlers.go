package categories

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/go-chi/chi"
	"github.com/lib/pq"
	"github.com/rahulsoibam/koubru-prod-api/middleware"
	"github.com/rahulsoibam/koubru-prod-api/utils"
)

// List all categories
func (a *App) List(w http.ResponseWriter, r *http.Request) {
	query := r.FormValue("q")
	perPage := r.FormValue("per_page")
	page := r.FormValue("page")

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

	categories, err := a.dbListTopics(query, limit, offset)
	if err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	utils.RespondWithJSON(w, http.StatusOK, &categories)
}

// Create a category
func (a *App) Create(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	userID, ok := ctx.Value(middleware.UserCtxKeys(0)).(int64)
	if !ok {
		utils.RespondWithError(w, http.StatusUnauthorized, "Problem with the user id associated with this token")
		return
	}

	var c Category
	err := json.NewDecoder(r.Body).Decode(&c)
	if err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, err.Error())
		return
	}
	if err = c.Validate(); err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	cres := Category{}
	err = a.DB.QueryRow("INSERT INTO Category (name, created_by) VALUES ($1, $2) RETURNING category_id, name", c.Name, userID).Scan(&cres.ID, &cres.Name)
	if err != nil {
		if e, ok := err.(*pq.Error); ok {
			utils.RespondWithError(w, http.StatusBadRequest, e.Detail)
			return
		}
		utils.RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	utils.RespondWithJSON(w, http.StatusOK, &cres)
}

// Get details of a category
func (a *App) Get(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("get details of a category"))
}

// Follow to follow a category
func (a *App) Follow(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	id := chi.URLParam(r, "id")
	followerID, ok := ctx.Value(middleware.UserCtxKeys(0)).(int64)
	if !ok {
		utils.RespondWithError(w, http.StatusUnauthorized, "Invalid user. Please try authenticating again")
		return
	}
	_, err := a.DB.Exec("INSERT INTO Category_Follower (category_id, user_id) VALUES ($1, $2)", id, followerID)
	if err != nil {
		if e, ok := err.(*pq.Error); ok {
			if e.Code == "23505" {
				utils.RespondWithError(w, http.StatusBadRequest, "You are already following this category")
				return
			} else {
				utils.RespondWithError(w, http.StatusInternalServerError, e.Detail)
				return
			}
		}
		utils.RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	utils.RespondWithMessage(w, http.StatusOK, "Successfully followed category")
}

// Unfollow to unfollow a category
func (a *App) Unfollow(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	id := chi.URLParam(r, "id")
	followerID, ok := ctx.Value(middleware.UserCtxKeys(0)).(int64)
	if !ok {
		utils.RespondWithError(w, http.StatusUnauthorized, "Invalid user. Please try authticating again")
		return
	}
	response, err := a.DB.Exec("DELETE FROM Category_Follower WHERE category_id=$1 AND user_id=$2", id, followerID)
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
		utils.RespondWithError(w, http.StatusBadRequest, "You do not follow this category")
		return
	}
	utils.RespondWithMessage(w, http.StatusOK, "Category unfollowed")
}

// BulkFollow to follow many categories at once
func (a *App) BulkFollow(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("bulk follow, first app entry"))
}

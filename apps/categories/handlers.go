package categories

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/rahulsoibam/koubru/errs"

	"github.com/lib/pq"
	"github.com/rahulsoibam/koubru/middleware"
	"github.com/rahulsoibam/koubru/types"
	"github.com/rahulsoibam/koubru/utils"
)

// List all categories
func (a *App) List(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	userID, auth := ctx.Value(middleware.AuthKeys("user_id")).(int64)

	categories := []types.CategoryForList{}
	var err error
	// Check authorization and perform query
	if auth {
		categories, err = a.AuthListQuery(ctx, userID)
	} else {
		categories, err = a.ListQuery(ctx)
	}
	if err != nil {
		log.Println(err)
		utils.RespondWithError(w, http.StatusInternalServerError, errs.InternalServerError)
		return
	}
	utils.RespondWithJSON(w, http.StatusOK, categories)
}

// Create a category
func (a *App) Create(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	userID, auth := ctx.Value(middleware.AuthKeys("user_id")).(int64)
	if !auth {
		log.Println(errs.Unauthorized)
		utils.RespondWithError(w, http.StatusUnauthorized, errs.Unauthorized)
		return
	}

	var c types.NewCategory
	err := json.NewDecoder(r.Body).Decode(&c)
	defer r.Body.Close()
	if err != nil {
		log.Println(err, r.Body)
		utils.RespondWithError(w, http.StatusBadRequest, errs.BadRequest)
		return
	}
	if err = c.Validate(); err != nil {
		log.Println(err, c)
		utils.RespondWithError(w, http.StatusBadRequest, errs.BadRequest)
		return
	}

	cres := types.Category{}

	cres, err = a.AuthCreateQuery(userID, c)
	if err != nil {
		if e, ok := err.(*pq.Error); ok {
			log.Println(e, e.Detail, e.Code)
			if e.Code == "23505" {
				utils.RespondWithError(w, http.StatusBadRequest, errs.CategoryAlreadyExists)
				return
			}
		}
		log.Println(err)
		utils.RespondWithError(w, http.StatusInternalServerError, errs.InternalServerError)
		return
	}
	utils.RespondWithJSON(w, http.StatusOK, cres)
}

// Get details of a category
func (a *App) Get(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	userID, auth := ctx.Value(middleware.AuthKeys("user_id")).(int64)
	categoryID := ctx.Value(middleware.CategoryKeys("category_id")).(int64)

	category := types.Category{}
	var err error
	if auth {
		category, err = a.AuthGetQuery(userID, categoryID)
	} else {
		category, err = a.GetQuery(categoryID)
	}
	if err != nil {
		log.Println(err)
		utils.RespondWithError(w, http.StatusInternalServerError, errs.InternalServerError)
		return
	}
	utils.RespondWithJSON(w, http.StatusOK, &category)
}

// Follow to follow a category
func (a *App) Follow(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	categoryID := ctx.Value(middleware.CategoryKeys("category_id")).(int64)
	followerID, auth := ctx.Value(middleware.AuthKeys("user_id")).(int64)

	if !auth {
		log.Println(errs.UnintendedExecution)
		utils.RespondWithError(w, http.StatusUnauthorized, errs.Unauthorized)
		return
	}
	_, err := a.DB.Exec("INSERT INTO Category_Follower (category_id, follower_id) VALUES ($1, $2)", categoryID, followerID)
	if err != nil {
		if e, ok := err.(*pq.Error); ok {
			if e.Code == "23505" {
				log.Println(err)
				utils.RespondWithError(w, http.StatusBadRequest, errs.CategoryFollowAlreadyFollowing)
				return
			}
		}
		log.Println(err)
		utils.RespondWithError(w, http.StatusInternalServerError, errs.InternalServerError)
		return
	}
	utils.RespondWithMessage(w, http.StatusOK, "Followed")
}

// Unfollow to unfollow a category
func (a *App) Unfollow(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	categoryID := ctx.Value(middleware.CategoryKeys("category_id")).(int64)
	followerID, auth := ctx.Value(middleware.AuthKeys("user_id")).(int64)
	if !auth {
		log.Println(errs.UnintendedExecution)
		utils.RespondWithError(w, http.StatusUnauthorized, errs.Unauthorized)
		return
	}
	response, err := a.DB.Exec("DELETE FROM Category_Follower WHERE category_id=$1 AND follower_id=$2", categoryID, followerID)
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

func (a *App) Followers(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	categoryID := ctx.Value(middleware.CategoryKeys("category_id")).(int64)
	userID, auth := ctx.Value(middleware.AuthKeys("user_id")).(int64)

	followers := []types.UserForFollowList{}
	var err error
	if auth {
		log.Println("authenticated category followers")
		followers, err = a.AuthFollowersQuery(ctx, userID, categoryID)
	} else {
		log.Println("unauthenticated category followers")
		followers, err = a.FollowersQuery(ctx, categoryID)
	}
	if err != nil {
		log.Println(err)
		utils.RespondWithError(w, http.StatusInternalServerError, errs.InternalServerError)
		return
	}
	utils.RespondWithJSON(w, http.StatusOK, followers)
}

func (a *App) Topics(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	categoryID := ctx.Value(middleware.CategoryKeys("category_id")).(int64)
	userID, auth := ctx.Value(middleware.AuthKeys("user_id")).(int64)

	topics := []types.TopicForList{}
	var err error
	if auth {
		topics, err = a.AuthTopicsQuery(ctx, userID, categoryID)
	} else {
		topics, err = a.TopicsQuery(ctx, categoryID)
	}
	if err != nil {
		log.Println(err)
		utils.RespondWithError(w, http.StatusInternalServerError, errs.InternalServerError)
		return
	}
	utils.RespondWithJSON(w, http.StatusOK, topics)
}

// BulkFollow to follow many categories at once
func (a *App) BulkFollow(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("bulk follow, first app entry"))
}

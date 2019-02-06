package opinions

import (
	"bytes"
	"io"
	"log"
	"net/http"
	"strconv"

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
		log.Println(err)
		utils.RespondWithError(w, http.StatusInternalServerError, errs.InternalServerError)
		return
	}
	utils.RespondWithJSON(w, http.StatusOK, opinions)
}

// Create to create an opinion
func (a *App) Create(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	userID, auth := ctx.Value(middleware.AuthKeys("user_id")).(int64)
	if !auth {
		log.Println(ctx)
		utils.RespondWithError(w, http.StatusUnauthorized, errs.Unauthorized)
		return
	}
	// Max upload size 200MB
	r.Body = http.MaxBytesReader(w, r.Body, 200<<20)
	defer r.Body.Close()
	reader, err := r.MultipartReader()
	if err != nil {
		log.Println(err)
		utils.RespondWithError(w, http.StatusInternalServerError, errs.InternalServerError)
		return
	}

	nr := types.NewReply{}
	for {
		part, err := reader.NextPart()
		if err == io.EOF {
			break
		} else if part.FormName() == "topic_id" {
			buf := new(bytes.Buffer)
			buf.ReadFrom(part)
			nr.TopicID, err = strconv.ParseInt(buf.String(), 10, 64)
			if err != nil {
				log.Println(err)
				utils.RespondWithError(w, http.StatusBadRequest, errs.BadRequest)
				return
			}
		} else if part.FormName() == "reaction" {
			buf := new(bytes.Buffer)
			buf.ReadFrom(part)
			nr.Reaction = buf.String()
		} else if part.FormName() == "is_anonymous" {
			buf := new(bytes.Buffer)
			buf.ReadFrom(part)
			nr.IsAnonymous, err = strconv.ParseBool(buf.String())
		} else if part.FormName() == "file" {
			uuid, err := a.Flake.NextID()
			if err != nil {
				log.Println(err)
				utils.RespondWithError(w, http.StatusInternalServerError, errs.InternalServerError)
			}
			filename := strconv.FormatUint(uuid, 10) + ".mp4"
			nr.Mp4, err = a.S3UploadOpinion(part, filename)
			if err != nil {
				log.Println(err)
				utils.RespondWithError(w, http.StatusBadRequest, errs.BadRequest)
				return
			}
		}
	}
	// Validate
	if err := nr.Validate(); err != nil {
		log.Println(nr, err)
		utils.RespondWithError(w, http.StatusBadRequest, err)
		return
	}

	opinion := types.Opinion{}
	opinion, err = a.AuthCreateReplyQuery(userID, nr)

	if err != nil {
		log.Println(err)
		utils.RespondWithError(w, http.StatusInternalServerError, errs.InternalServerError)
		return
	}

	utils.RespondWithJSON(w, http.StatusOK, opinion)
}

func (a *App) Reply(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	userID, auth := ctx.Value(middleware.AuthKeys("user_id")).(int64)
	ctxOpinon := ctx.Value(middleware.OpinionKeys("ctx_opinion")).(types.ContextOpinion)
	if !auth {
		log.Println(ctx)
		utils.RespondWithError(w, http.StatusUnauthorized, errs.Unauthorized)
		return
	}

	// Max upload size 200MB
	r.Body = http.MaxBytesReader(w, r.Body, 200<<20)
	defer r.Body.Close()
	reader, err := r.MultipartReader()
	if err != nil {
		log.Println(err)
		utils.RespondWithError(w, http.StatusInternalServerError, errs.InternalServerError)
		return
	}

	nr := types.NewReply{}
	nr.ParentID = ctxOpinon.ID
	nr.TopicID = ctxOpinon.TopicID
	for {
		part, err := reader.NextPart()
		if err == io.EOF {
			break
		} else if part.FormName() == "reaction" {
			buf := new(bytes.Buffer)
			buf.ReadFrom(part)
			nr.Reaction = buf.String()
		} else if part.FormName() == "is_anonymous" {
			buf := new(bytes.Buffer)
			buf.ReadFrom(part)
			nr.IsAnonymous, err = strconv.ParseBool(buf.String())
		} else if part.FormName() == "file" {
			uuid, err := a.Flake.NextID()
			if err != nil {
				log.Println(err)
				utils.RespondWithError(w, http.StatusInternalServerError, errs.InternalServerError)
			}
			filename := strconv.FormatUint(uuid, 10) + ".mp4"
			nr.Mp4, err = a.S3UploadOpinion(part, filename)
			if err != nil {
				log.Println(err)
				utils.RespondWithError(w, http.StatusBadRequest, errs.BadRequest)
				return
			}
		}
	}

	// Validate
	if err := nr.Validate(); err != nil {
		log.Println(nr, err)
		utils.RespondWithError(w, http.StatusBadRequest, err)
		return
	}

	opinion := types.Opinion{}
	opinion, err = a.AuthCreateReplyQuery(userID, nr)

	if err != nil {
		log.Println(err)
		utils.RespondWithError(w, http.StatusInternalServerError, errs.InternalServerError)
		return
	}

	utils.RespondWithJSON(w, http.StatusOK, opinion)
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
		log.Println(err)
		utils.RespondWithError(w, http.StatusInternalServerError, errs.InternalServerError)
		return
	}

	utils.RespondWithJSON(w, http.StatusOK, opinion)
}

func (a *App) Breadcrumbs(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	opinionID := ctx.Value(middleware.OpinionKeys("opinion_id")).(int64)
	userID, auth := ctx.Value(middleware.AuthKeys("user_id")).(int64)

	breadcrumbs := []types.Breadcrumb{}
	var err error
	if auth {
		breadcrumbs, err = a.AuthBreadcrumbsQuery(userID, opinionID)
	} else {
		breadcrumbs, err = a.BreadcrumbsQuery(opinionID)
	}
	if err != nil {
		log.Println(err)
		utils.RespondWithError(w, http.StatusInternalServerError, errs.InternalServerError)
		return
	}

	utils.RespondWithJSON(w, http.StatusOK, breadcrumbs)
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
	ctx := r.Context()
	opinionID := ctx.Value(middleware.OpinionKeys("opinion_id")).(int64)
	userID, auth := ctx.Value(middleware.AuthKeys("user_id")).(int64)

	replies := []types.Opinion{}
	var err error
	if auth {
		replies, err = a.AuthRepliesQuery(userID, opinionID)
	} else {
		replies, err = a.RepliesQuery(opinionID)
	}
	if err != nil {
		log.Println(err)
		utils.RespondWithError(w, http.StatusInternalServerError, errs.InternalServerError)
		return
	}

	utils.RespondWithJSON(w, http.StatusOK, replies)
}

// Report to report an opinion
func (a *App) Report(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Report opinion"))
}

// Vote to vote on an opinion
func (a *App) Vote(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Vote on an opinion"))
}

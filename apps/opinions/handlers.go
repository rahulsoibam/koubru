package opinions

import (
	"bytes"
	"database/sql"
	"errors"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/aws/aws-sdk-go/service/sqs"

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
	var filename string
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
		} else if part.FormName() == "file" {
			uuid, err := a.Flake.NextID()
			if err != nil {
				log.Println(err)
				utils.RespondWithError(w, http.StatusInternalServerError, errs.InternalServerError)
			}
			filename = strconv.FormatUint(uuid, 10)
			err = a.S3UploadOpinion(part, filename)
			if err != nil {
				log.Println(err)
				utils.RespondWithError(w, http.StatusBadRequest, errs.BadRequest)
				return
			}
		}
	}

	links, err := a.PollSQSAndGetLinks(sqs.New(a.Sess), os.Getenv("S3_BUCKET"), filename, os.Getenv("QUEUE_URL"))
	if err != nil {
		log.Println(err)
		utils.RespondWithError(w, http.StatusBadRequest, errs.OpinionBadPayload)
		return
	}
	nr.Source = links.Source
	nr.Hls = links.Hls
	nr.Thumbnails = append(nr.Thumbnails, links.Thumbnail)

	// Validate
	if err := nr.Validate(); err != nil {
		log.Println(nr, err)
		utils.RespondWithError(w, http.StatusBadRequest, err)
		return
	}

	log.Println("New Reply Object: ", nr)

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

	var filename string
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
		} else if part.FormName() == "file" {
			uuid, err := a.Flake.NextID()
			if err != nil {
				log.Println(err)
				utils.RespondWithError(w, http.StatusInternalServerError, errs.InternalServerError)
			}
			filename := strconv.FormatUint(uuid, 10)
			err = a.S3UploadOpinion(part, filename)
			if err != nil {
				log.Println(err)
				utils.RespondWithError(w, http.StatusBadRequest, errs.BadRequest)
				return
			}
		}
	}
	links, err := a.PollSQSAndGetLinks(sqs.New(a.Sess), os.Getenv("S3_BUCKET"), filename, os.Getenv("QUEUE_URL"))
	if err != nil {
		log.Println(err)
		utils.RespondWithError(w, http.StatusBadRequest, err)
	}

	nr.Source = links.Source
	nr.Hls = links.Hls
	nr.Thumbnails = append(nr.Thumbnails, links.Thumbnail)

	// Validate
	if err := nr.Validate(); err != nil {
		log.Println(nr, err)
		utils.RespondWithError(w, http.StatusBadRequest, err)
		return
	}
	log.Println("New Reply Object: ", nr)

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
	w.Write([]byte("Follow opinion is not supported yet"))
}

// Unfollow to unfollow an opinion
func (a *App) Unfollow(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Unfollow opinion is not supported by the platform yet"))
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
	ctx := r.Context()
	userID, auth := ctx.Value(middleware.AuthKeys("user_id")).(int64)
	opinionID := ctx.Value(middleware.OpinionKeys("opinion_id")).(int64)

	if !auth {
		log.Println(ctx, errs.UnintendedExecution)
		utils.RespondWithError(w, http.StatusUnauthorized, errs.Unauthorized)
		return
	}

	vote := r.FormValue("vote")
	var voteBool bool
	var safe bool

	err := a.DB.QueryRow("SELECT vote FROM Opinion_Vote WHERE voter_id=$1 AND opinion_id=$2", userID, opinionID).Scan(&voteBool)
	if err != nil {
		if err == sql.ErrNoRows {
			safe = true
		}
		log.Println(err)
		utils.RespondWithError(w, http.StatusInternalServerError, errs.InternalServerError)
		return
	}

	if vote == "upvote" {
		if safe {
			log.Println("Upvoting when user has not voted yet")
			_, err := a.DB.Exec("INSERT INTO Opinion_Vote (voter_id, opinion_id, vote) VALUES ($1, $2, $3)", userID, opinionID, true)
			if err != nil {
				log.Println(err)
				utils.RespondWithError(w, http.StatusInternalServerError, errs.InternalServerError)
				return
			}
		} else if !voteBool {
			log.Println("Upvoting when user was downvoting")
			_, err := a.DB.Exec("UPDATE Opinion_Vote SET vote=$1 WHERE voter_id=$2 AND opinion_id=$3", true, userID, opinionID)
			if err != nil {
				log.Println(err)
				utils.RespondWithError(w, http.StatusInternalServerError, errs.InternalServerError)
				return
			}
		} else if voteBool {
			log.Println("Upvoting when user has already upvoted on the post")
			_, err := a.DB.Exec("DELETE FROM Opinion_Vote WHERE voter_id=$1 AND opinion_id=$2", userID, opinionID)
			if err != nil {
				log.Println(err)
				utils.RespondWithError(w, http.StatusInternalServerError, errs.InternalServerError)
				return
			}
		}
	} else if vote == "downvote" {
		if safe {
			log.Println("Downvoting when user has not voted yet")
			_, err := a.DB.Exec("INSERT INTO Opinion_Vote (voter_id, opinion_id, vote) VALUES ($1, $2, $3)", userID, opinionID, false)
			if err != nil {
				log.Println(err)
				utils.RespondWithError(w, http.StatusInternalServerError, errs.InternalServerError)
				return
			}
		} else if voteBool {
			log.Println("Downvoting when user was upvoting")
			_, err := a.DB.Exec("UPDATE Opinion_Vote SET vote=$1 WHERE voter_id=$2 AND opinion_id=$3", false, userID, opinionID)
			if err != nil {
				log.Println(err)
				utils.RespondWithError(w, http.StatusInternalServerError, errs.InternalServerError)
				return
			}
		} else if !voteBool {
			log.Println("Downvoting when user has already downvoted on the post")
			_, err := a.DB.Exec("DELETE FROM Opinion_Vote WHERE voter_id=$1 AND opinion_id=$2", userID, opinionID)
			if err != nil {
				log.Println(err)
				utils.RespondWithError(w, http.StatusInternalServerError, errs.InternalServerError)
				return
			}
		}
	} else {
		log.Println("Invalid vote", vote)
		utils.RespondWithError(w, http.StatusBadRequest, errors.New(vote+" is not a valid vote type"))
		return
	}

	utils.RespondWithMessage(w, http.StatusOK, "Action successful")
}

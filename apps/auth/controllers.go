package auth

import (
	"crypto/hmac"
	"crypto/sha256"
	"database/sql"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"os"
	"strings"

	"github.com/lib/pq"
	"github.com/rahulsoibam/koubru-prod-api/middleware"

	"github.com/rahulsoibam/koubru-prod-api/authutils"
	"github.com/rahulsoibam/koubru-prod-api/authutils/googlejwt"
	"github.com/rahulsoibam/koubru-prod-api/utils"
)

var (
	errUserNotFound     = errors.New("User not found")
	errNoPasswordSet    = errors.New("A password is not set for this account. Login using social account or create password")
	errPasswordNotMatch = errors.New("Password does not match")
)

// Login using username/phone/email and password
func (a *App) Login(w http.ResponseWriter, r *http.Request) {
	var creds Credentials
	var encodedHash string
	var err error
	var userID int64

	err = json.NewDecoder(r.Body).Decode(&creds)
	if err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, err.Error())
		return
	}
	defer r.Body.Close()

	// Get user login type, out of email and username
	loginType, err := creds.ValidateAndLoginType()
	if err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	// Get argon2 encoded password string
	if loginType == "username" {
		userID, err = dbGetUserIDUsingUsername(a.DB, creds.User)
	} else if loginType == "email" {
		userID, err = dbGetUserIDUsingEmail(a.DB, creds.User)
	}
	if err != nil {
		switch err {
		case sql.ErrNoRows:
			utils.RespondWithError(w, http.StatusNotFound, errUserNotFound.Error())
			return
		default:
			utils.RespondWithError(w, http.StatusInternalServerError, err.Error())
			return
		}
	}

	// Get argon2 encoded password hash string
	encodedHash, err = dbGetEncodedHashUsingUserID(a.AuthDB, userID)
	if err != nil {
		switch err {
		case sql.ErrNoRows:
			utils.RespondWithError(w, http.StatusBadRequest, errNoPasswordSet.Error())
			return
		default:
			utils.RespondWithError(w, http.StatusInternalServerError, err.Error())
			return
		}
	}

	// Compare user provided password and encodedHash from the database
	var match bool
	match, err = authutils.ComparePasswordAndHash(creds.Password, encodedHash)
	if err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	if !match {
		utils.RespondWithError(w, http.StatusBadRequest, errPasswordNotMatch.Error())
		return
	}

	token, err := authutils.Authenticate(a.AuthCache, a.AuthDB, userID, r.UserAgent())
	if err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	utils.RespondWithJSON(w, http.StatusOK, &token)
}

// Register a new account
func (a *App) Register(w http.ResponseWriter, r *http.Request) {
	var nu *NewUser
	var err error
	err = json.NewDecoder(r.Body).Decode(&nu)
	defer r.Body.Close()
	if err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, err.Error())
		return
	}
	if err := nu.Validate(); err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	n := a.AuthCache.SIsMember("usernames", nu.Username)
	if n.Val() {
		utils.RespondWithError(w, http.StatusBadRequest, "Username "+nu.Username+" already exists. Please enter another username")
		return
	}

	n = a.AuthCache.SIsMember("emails", nu.Email)
	if n.Val() {
		utils.RespondWithError(w, http.StatusBadRequest, "Email "+nu.Email+" already exists. Please enter another email")
		return
	}

	userID, err := dbRegisterUser(a.DB, nu)
	if err != nil {
		if e, ok := err.(*pq.Error); ok {
			utils.RespondWithError(w, http.StatusBadRequest, e.Detail)
			return
		}
		utils.RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	// Generate password hash
	encodedHash, err := authutils.GenerateFromPassword(nu.Password, a.Argon2Params)
	if err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	// Store password in separate database
	err = dbStorePassword(a.AuthDB, userID, encodedHash)
	if err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	a.AuthCache.SAdd("usernames", nu.Username)
	a.AuthCache.SAdd("emails", nu.Email)

	// Store authentication details in data layers and return access token
	token, err := authutils.Authenticate(a.AuthCache, a.AuthDB, userID, r.UserAgent())
	if err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	utils.RespondWithJSON(w, http.StatusOK, &token)
}

// Facebook auth using facebook account
func (a *App) Facebook(w http.ResponseWriter, r *http.Request) {
	fmt.Println("inside facebook")
	facebookAccessToken := r.FormValue("facebook_access_token")
	appSecretProofGenerator := hmac.New(sha256.New, []byte(os.Getenv("FB_CLIENT_SECRET")))
	appSecretProofGenerator.Write([]byte(facebookAccessToken))

	appSecretProof := hex.EncodeToString(appSecretProofGenerator.Sum(nil))

	fmt.Println("requesting facebook")
	response, err := http.Get("https://graph.facebook.com/me?fields=id,name,picture.type(large),email&access_token=" + facebookAccessToken + "&appsecret_proof=" + appSecretProof)
	if err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, err.Error())
		return
	}
	defer response.Body.Close()
	fmt.Println("done requesting facebook")

	type FacebookUser struct {
		ID      string `json:"id"`
		Name    string `json:"name"`
		Picture struct {
			Data struct {
				Height       int    `json:"height"`
				IsSilhouette bool   `json:"is_silhouette"`
				URL          string `json:"url"`
				Width        int    `json:"width"`
			} `json:"data"`
		} `json:"picture"`
		Email string `json:"email"`
	}

	var fu FacebookUser
	fmt.Println("decoding json")
	err = json.NewDecoder(response.Body).Decode(&fu)
	if err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, err.Error())
		return
	}
	fmt.Println("done decoding json")

	userID, err := dbGetUserIDUsingFacebook(a.DB, fu.ID)
	if err != nil {
		switch err {
		case sql.ErrNoRows:
			var nu NewUser
			nu.Username = ""
			nu.Email = fu.Email
			nu.Picture = fu.Picture.Data.URL
			nu.FullName = fu.Name
			utils.RespondWithJSON(w, http.StatusUnauthorized, &nu)
			return
		default:
			utils.RespondWithError(w, http.StatusInternalServerError, err.Error())
			return
		}
	}

	token, err := authutils.Authenticate(a.AuthCache, a.AuthDB, userID, r.UserAgent())
	if err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	utils.RespondWithJSON(w, http.StatusOK, &token)
}

// Google auth using google account
func (a *App) Google(w http.ResponseWriter, r *http.Request) {
	googleIDToken := r.FormValue("google_id_token")
	v := googlejwt.GoogleIDTokenVerifier{}
	iosaud := "451796869752-sbdnk7c82edf91g3hernllknfmpngifl.apps.googleusercontent.com"
	andaud := "451796869752-muqbuv2jn8o9hce5c64gl52ibm2gbkmi.apps.googleusercontent.com"
	err := v.VerifyIDToken(googleIDToken, []string{
		iosaud, andaud,
	})
	if err != nil {
		utils.RespondWithJSON(w, http.StatusBadRequest, err.Error())
		return
	}
	claimSet, err := googlejwt.Decode(googleIDToken)

	userID, err := dbGetUserIDUsingGoogle(a.DB, claimSet.Sub)
	if err != nil {
		switch err {
		case sql.ErrNoRows:
			var nu NewUser

			// a.AuthCache.Set()
			nu.FullName = claimSet.Name
			nu.Email = claimSet.Email
			nu.Picture = claimSet.Picture
			utils.RespondWithJSON(w, http.StatusUnauthorized, &nu)
			return
		default:
			utils.RespondWithError(w, http.StatusInternalServerError, err.Error())
			return
		}
	}
	token, err := authutils.Authenticate(a.AuthCache, a.AuthDB, userID, r.UserAgent())
	if err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	utils.RespondWithJSON(w, http.StatusOK, &token)
}

// LinkGoogle is used to link a google account
func (a *App) LinkGoogle(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("link google"))
}

// LinkFacebook is used to link a facebook account
func (a *App) LinkFacebook(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("link facebook"))
}

// Logout user
func (a *App) Logout(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	userID := ctx.Value(middleware.UserCtxKeys(0)).(int64)
	token := ctx.Value(middleware.UserCtxKeys(1)).(string)

	n := a.AuthCache.Del(token)
	if n.Val() != 1 {
		return
	}
	res, err := a.AuthDB.Exec("DELETE FROM session WHERE token=$1 AND user_id=$2", token, userID)
	if err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	count, err := res.RowsAffected()
	if err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	if count == 0 {
		utils.RespondWithError(w, http.StatusInternalServerError, "Session deleted")
		return
	}

	utils.RespondWithJSON(w, http.StatusOK, "Logged out successfully")
}

// CheckEmail queries the database and checks if a username is already registered
func (a *App) CheckEmail(w http.ResponseWriter, r *http.Request) {
	email := r.FormValue("email")
	email = strings.ToLower(email)
	if err := utils.ValidateEmail(email); err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	exists := a.AuthCache.SIsMember("emails", email)
	if exists.Val() {
		utils.RespondWithError(w, http.StatusBadRequest, "Email "+email+" is already used by another account.")
		return
	}
	utils.RespondWithJSON(w, http.StatusOK, email+" is a valid email")
}

// CheckUsername queries the database and checks if an email is already registered
func (a *App) CheckUsername(w http.ResponseWriter, r *http.Request) {
	username := r.FormValue("username")
	username = strings.ToLower(username)
	if err := utils.ValidateUsername(username); err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	exists := a.AuthCache.SIsMember("usernames", username)
	if exists.Val() {
		utils.RespondWithError(w, http.StatusBadRequest, "Username "+username+" is already used by another account.")
		return
	}
	utils.RespondWithJSON(w, http.StatusOK, username+" is a valid username")
}

// VerifyEmail verifies an email by sending a one time link
func (a *App) VerifyEmail(w http.ResponseWriter, r *http.Request) {
	// TODO
	// Send confirmation mail using Sendgrid
	// Create server rendered HTML to show confirmation status
	email := r.FormValue("email")
	otp := r.FormValue("otp")
	redisKey := fmt.Sprintf("verification:email:%s:otp:%s", email, otp)
	n := a.AuthCache.Del(redisKey)
	if n.Val() == 1 {
		w.Write([]byte("Verification successful"))
	}
	w.Write([]byte("Invalid! It is possible that the code might have expired. Please try again"))
}

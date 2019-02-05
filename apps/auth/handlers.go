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
	"time"

	"github.com/rahulsoibam/koubru/errs"

	"github.com/lib/pq"
	"github.com/rahulsoibam/koubru/middleware"

	"github.com/rahulsoibam/koubru/authutils"
	"github.com/rahulsoibam/koubru/authutils/googlejwt"
	"github.com/rahulsoibam/koubru/utils"
)

func (a *App) authenticate(userID int64, bearerToken string, userAgent string) (*Token, error) {
	var err error
	expiry := 60 * 60 * 24 * 30 * time.Second
	// Store token as session in Redis
	err = a.AuthCache.Set(bearerToken, userID, expiry).Err()
	if err != nil {
		return nil, err
	}

	_, err = a.AuthDB.Exec("INSERT INTO Session (user_id, token, user_agent) VALUES ($1, $2, $3)", userID, bearerToken, userAgent)
	if err != nil {
		return nil, err
	}

	token := Token{
		TokenType:   "Bearer",
		AccessToken: bearerToken,
		Expires:     expiry.Nanoseconds() / 1e9,
	}

	return &token, nil
}

// Login using username/phone/email and password
func (a *App) Login(w http.ResponseWriter, r *http.Request) {
	var creds Credentials
	var encodedHash string
	var err error
	var userID int64

	err = json.NewDecoder(r.Body).Decode(&creds)
	if err != nil {
		a.Log.Infoln(err)
		utils.RespondWithError(w, http.StatusBadRequest, err)
		return
	}
	defer r.Body.Close()

	// Get user login type, out of email and username
	loginType, err := creds.ValidateAndLoginType()
	if err != nil {
		a.Log.Errorln(err, creds)
		utils.RespondWithError(w, http.StatusBadRequest, err)
		return
	}

	// Get argon2 encoded password string
	if loginType == "username" {
		userID, err = a.dbGetUserIDUsingUsername(creds.User)
	} else if loginType == "email" {
		userID, err = a.dbGetUserIDUsingEmail(creds.User)
	}
	if err != nil {
		switch err {
		case sql.ErrNoRows:
			a.Log.Infoln(err)
			utils.RespondWithError(w, http.StatusNotFound, errs.UserNotFound)
			return
		default:
			a.Log.Errorln(err)
			utils.RespondWithError(w, http.StatusInternalServerError, errs.InternalServerError)
			return
		}
	}

	// Get argon2 encoded password hash string
	encodedHash, err = a.dbGetEncodedHashUsingUserID(userID)
	if err != nil {
		switch err {
		case sql.ErrNoRows:
			a.Log.Infoln(err)
			utils.RespondWithError(w, http.StatusBadRequest, errs.NoPasswordSet)
			return
		default:
			a.Log.Errorln(err)
			utils.RespondWithError(w, http.StatusInternalServerError, errs.InternalServerError)
			return
		}
	}

	// Compare user provided password and encodedHash from the database
	var match bool
	match, err = authutils.ComparePasswordAndHash(creds.Password, encodedHash)
	if err != nil {
		a.Log.Errorln(err)
		utils.RespondWithError(w, http.StatusInternalServerError, errs.InternalServerError)
		return
	}
	if !match {
		utils.RespondWithError(w, http.StatusBadRequest, errs.PasswordNotMatch)
		return
	}

	// Generate bearer token
	bearerToken, err := authutils.GenerateSecureToken(256)
	if err != nil {
		a.Log.Errorln(err)
		utils.RespondWithError(w, http.StatusInternalServerError, errs.InternalServerError)
		return
	}

	// Authenticate and return access token
	token, err := a.authenticate(userID, bearerToken, r.UserAgent())
	if err != nil {
		a.Log.Errorln(err)
		utils.RespondWithError(w, http.StatusInternalServerError, err)
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
		a.Log.Infoln(err, nu)
		utils.RespondWithError(w, http.StatusBadRequest, err)
		return
	}
	if err := nu.Validate(); err != nil {
		a.Log.Infoln(err, nu)
		utils.RespondWithError(w, http.StatusBadRequest, err)
		return
	}

	n := a.AuthCache.SIsMember("usernames", nu.Username)
	if n.Val() {
		errUsernameAlreadyExist := errors.New("Username " + nu.Username + " already exists. Please enter another username")
		a.Log.Infoln()
		utils.RespondWithError(w, http.StatusBadRequest, errUsernameAlreadyExist)
		return
	}

	n = a.AuthCache.SIsMember("emails", nu.Email)
	if n.Val() {
		errEmailAlreadyExist := errors.New("Email " + nu.Email + " already exists. Please enter another email")
		utils.RespondWithError(w, http.StatusBadRequest, errEmailAlreadyExist)
		return
	}

	userID, err := a.dbRegisterUser(nu)
	if err != nil {
		if e, ok := err.(*pq.Error); ok {
			a.Log.Infoln(e, e.Detail)
			utils.RespondWithError(w, http.StatusBadRequest, errors.New(e.Detail))
			return
		}
		a.Log.Errorln(err)
		utils.RespondWithError(w, http.StatusInternalServerError, err)
		return
	}
	// Generate password hash
	encodedHash, err := authutils.GenerateFromPassword(nu.Password, a.Argon2Params)
	if err != nil {
		a.Log.Errorln(err)
		utils.RespondWithError(w, http.StatusInternalServerError, errs.InternalServerError)
		return
	}

	// Store password in separate database
	err = a.dbStorePassword(userID, encodedHash)
	if err != nil {
		a.Log.Errorln(err)
		utils.RespondWithError(w, http.StatusInternalServerError, errs.InternalServerError)
		return
	}

	a.AuthCache.SAdd("usernames", nu.Username)
	a.AuthCache.SAdd("emails", nu.Email)

	bearerToken, err := authutils.GenerateSecureToken(256)
	if err != nil {
		a.Log.Errorln(err)
		utils.RespondWithError(w, http.StatusInternalServerError, errs.InternalServerError)
		return
	}

	// Store authentication details in data layers and return access token
	token, err := a.authenticate(userID, bearerToken, r.UserAgent())
	if err != nil {
		a.Log.Errorln(err)
		utils.RespondWithError(w, http.StatusInternalServerError, errs.InternalServerError)
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
		a.Log.Infoln(err)
		utils.RespondWithError(w, http.StatusBadRequest, err)
		return
	}
	if response.StatusCode != http.StatusOK {
		a.Log.Infoln(response.StatusCode, response.Status)
		utils.RespondWithError(w, http.StatusBadRequest, errors.New("Error fetching details from facebook. Check the token and try again"))
		return
	}

	defer response.Body.Close()
	fmt.Println("done requesting facebook")

	var fu FacebookUser
	fmt.Println("decoding json")
	err = json.NewDecoder(response.Body).Decode(&fu)
	if err != nil {
		a.Log.Infoln(err, fu)
		utils.RespondWithError(w, http.StatusBadRequest, errs.BadRequest)
		return
	}
	fmt.Println("done decoding json")

	userID, err := a.dbGetUserIDUsingFacebook(fu.ID)
	if err != nil {
		switch err {
		case sql.ErrNoRows:
			username, err := fu.GenerateUsername(a.AuthCache)
			if err != nil {
				a.Log.Errorln(err)
				utils.RespondWithError(w, http.StatusBadRequest, err)
				return
			}
			userID, err = a.dbRegisterUserUsingFacebook(fu, username)
			if err != nil {
				if e, ok := err.(*pq.Error); ok {
					utils.RespondWithError(w, http.StatusInternalServerError, errors.New(e.Detail))
					return
				}

				utils.RespondWithError(w, http.StatusInternalServerError, errs.InternalServerError)
				return
			}
			a.AuthCache.SAdd("usernames", username)
		default:
			a.Log.Errorln(err)
			utils.RespondWithError(w, http.StatusInternalServerError, errs.InternalServerError)
			return
		}
	}

	bearerToken, err := authutils.GenerateSecureToken(256)
	if err != nil {
		a.Log.Errorln(err)
		utils.RespondWithError(w, http.StatusInternalServerError, err)
		return
	}

	token, err := a.authenticate(userID, bearerToken, r.UserAgent())
	if err != nil {
		a.Log.Errorln(err)
		utils.RespondWithError(w, http.StatusInternalServerError, err)
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
		utils.RespondWithError(w, http.StatusBadRequest, err)
		return
	}
	cs, err := googlejwt.Decode(googleIDToken)
	if err != nil {
		a.Log.Errorln(err)
		utils.RespondWithError(w, http.StatusInternalServerError, errs.InternalServerError)
		return
	}

	userID, err := a.dbGetUserIDUsingGoogle(cs.Sub)
	if err != nil {
		switch err {
		case sql.ErrNoRows:
			username, err := cs.GenerateUsername(a.AuthCache)
			if err != nil {
				utils.RespondWithError(w, http.StatusBadRequest, err)
				return
			}
			userID, err = a.dbRegisterUserUsingGoogle(cs, username)
			if err != nil {
				if e, ok := err.(*pq.Error); ok {
					a.Log.Errorln(e)
					utils.RespondWithError(w, http.StatusInternalServerError, errors.New(e.Detail))
					return
				}
				a.Log.Errorln(err)
				utils.RespondWithError(w, http.StatusInternalServerError, errs.InternalServerError)
				return
			}
			a.AuthCache.SAdd("usernames", username)
		default:
			a.Log.Errorln(err)
			utils.RespondWithError(w, http.StatusInternalServerError, errs.InternalServerError)
			return
		}
	}
	bearerToken, err := authutils.GenerateSecureToken(256)
	if err != nil {
		a.Log.Errorln(err)
		utils.RespondWithError(w, http.StatusInternalServerError, errs.InternalServerError)
		return
	}
	token, err := a.authenticate(userID, bearerToken, r.UserAgent())
	if err != nil {
		a.Log.Errorln(err)
		utils.RespondWithError(w, http.StatusInternalServerError, errs.InternalServerError)
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
	userID := ctx.Value(middleware.AuthKeys("user_id")).(int64)
	token := ctx.Value(middleware.AuthKeys("auth_token")).(string)

	n := a.AuthCache.Del(token)
	if n.Val() != 1 {
		utils.RespondWithError(w, http.StatusBadRequest, errors.New("The session does not exist"))
		return
	}
	res, err := a.AuthDB.Exec("DELETE FROM session WHERE token=$1 AND user_id=$2", token, userID)
	if err != nil {

		utils.RespondWithError(w, http.StatusInternalServerError, errs.InternalServerError)
		return
	}
	count, err := res.RowsAffected()
	if err != nil {
		a.Log.Errorln(err)
		utils.RespondWithError(w, http.StatusInternalServerError, errs.InternalServerError)
		return
	}
	if count == 0 {
		utils.RespondWithError(w, http.StatusInternalServerError, errors.New("Session deleted"))
		return
	}

	utils.RespondWithMessage(w, http.StatusOK, "Logged out successfully")
}

// CheckEmail queries the database and checks if a username is already registered
func (a *App) CheckEmail(w http.ResponseWriter, r *http.Request) {
	email := r.FormValue("email")
	email = strings.ToLower(email)
	if err := utils.ValidateEmail(email); err != nil {
		a.Log.Errorln(err)
		utils.RespondWithError(w, http.StatusBadRequest, errs.BadRequest)
		return
	}

	exists := a.AuthCache.SIsMember("emails", email)
	if exists.Val() {
		errEmailAlreadyExist := errors.New("Email is already used by another account.")
		utils.RespondWithError(w, http.StatusBadRequest, errEmailAlreadyExist)
		return
	}
	utils.RespondWithMessage(w, http.StatusOK, email+" is a valid email")
}

// CheckUsername queries the database and checks if an email is already registered
func (a *App) CheckUsername(w http.ResponseWriter, r *http.Request) {
	username := r.FormValue("username")
	username = strings.ToLower(username)
	if err := utils.ValidateUsername(username); err != nil {
		a.Log.Infoln(err)
		utils.RespondWithError(w, http.StatusBadRequest, err)
		return
	}

	exists := a.AuthCache.SIsMember("usernames", username)
	if exists.Val() {
		utils.RespondWithError(w, http.StatusBadRequest, errors.New("Username is already used by another account."))
		return
	}
	utils.RespondWithMessage(w, http.StatusOK, username+" is a valid username")
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

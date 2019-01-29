package auth

import (
	"errors"
	"strconv"
	"strings"

	"github.com/go-redis/redis"

	"github.com/rahulsoibam/koubru-prod-api/utils"
)

// Token struct for storing token
type Token struct {
	TokenType   string `json:"token_type"`
	AccessToken string `json:"access_token"`
	Expires     int64  `json:"expires"`
}

// Credentials struct for storing login input
type Credentials struct {
	User     string `json:"user"`
	Password string `json:"password"`
}

// ValidateAndLoginType validates the credentials and returns the login type as string
func (c *Credentials) ValidateAndLoginType() (string, error) {
	c.User = strings.ToLower(c.User)
	if err := utils.ValidatePassword(c.Password); err != nil {
		return "", err
	}
	if utils.UsernameRegex.MatchString(c.User) {
		err := utils.ValidateUsername(c.User)
		if err != nil {
			return "", err
		}
		return "username", nil
	}
	if utils.EmailRegex.MatchString(c.User) {
		err := utils.ValidateEmail(c.User)
		if err != nil {
			return "", err
		}
		return "email", nil
	}
	return "", errors.New("Invalid user field")
}

// NewUser struct for storing user registration data
type NewUser struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password,omitempty"`
	FullName string `json:"full_name"`
}

// Validate new user details
// TODO Remove dots and plus signs from emails
func (nu *NewUser) Validate() error {
	nu.Username = strings.ToLower(nu.Username)
	if err := utils.ValidateUsername(nu.Username); err != nil {
		return err
	}

	nu.Email = strings.ToLower(nu.Email)
	if err := utils.ValidateEmail(nu.Email); err != nil {
		return err
	}

	if err := utils.ValidatePassword(nu.Password); err != nil {
		return err
	}

	if err := utils.ValidateFullName(nu.FullName); err != nil {
		return err
	}

	return nil
}

// FacebookUser stores the Facebook User Data response
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

// GenerateUsername generates a username from the facebook name
func (fu *FacebookUser) GenerateUsername(authCache *redis.Client) (string, error) {
	// Trim non-username characters from Name
	username := utils.UsernameInverseRegex.ReplaceAllString(fu.Name, "")
	if utils.UsernameRegex.MatchString(username) {
		username = strings.ToLower(username)
		err := utils.ValidateUsername(username)
		if err != nil {
			return "", err
		}

		exists := authCache.SIsMember("usernames", username)
		if !exists.Val() {
			return username, nil
		}
		for i := 2; i < 100; i++ {
			exists = authCache.SIsMember("usernames", username+strconv.Itoa(i))
			if !exists.Val() {
				return username + strconv.Itoa(i), nil
			}
		}
	}

	return "", errors.New("Cannot generate username from facebook name")
}

// // GoogleUser stores the Google User Data response
// type GoogleUser struct {
// 	jws.ClaimSet
// 	Email         string `json:"email"`
// 	EmailVerified bool   `json:"email_verified"`
// 	Name          string `json:"name"`
// 	Picture       string `json:"picture"`
// 	GivenName     string `json:"given_name"`
// 	FamilyName    string `json:"family_name"`
// 	Locale        string `json:"locale"`
// }

// // GenerateUsername generates a username from the Google name
// func (gu *GoogleUser) GenerateUsername(authCache *redis.Client) (string, error) {
// 	// Trim non-username characters from Name
// 	username := utils.UsernameInverseRegex.ReplaceAllString(gu.Name, "")
// 	if utils.UsernameRegex.MatchString(username) {
// 		err := utils.ValidateEmail(username)
// 		if err != nil {
// 			return "", err
// 		}
// 	}
// 	exists := authCache.SIsMember("usernames", username)
// 	if !exists.Val() {
// 		return username, nil
// 	}
// 	for i := 2; i < 100; i++ {
// 		exists = authCache.SIsMember("usernames", username+strconv.Itoa(i))
// 		if !exists.Val() {
// 			return username + strconv.Itoa(i), nil
// 		}
// 	}

// 	return "", errors.New("Cannot generate username from facebook name")
// }

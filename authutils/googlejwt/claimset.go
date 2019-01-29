package googlejwt

import (
	"errors"
	"strconv"

	"github.com/go-redis/redis"
	"github.com/rahulsoibam/koubru-prod-api/utils"
	"golang.org/x/oauth2/jws"
)

// ClaimSet stores the decoded id token details
type ClaimSet struct {
	jws.ClaimSet
	Email         string `json:"email"`
	EmailVerified bool   `json:"email_verified"`
	Name          string `json:"name"`
	Picture       string `json:"picture"`
	GivenName     string `json:"given_name"`
	FamilyName    string `json:"family_name"`
	Locale        string `json:"locale"`
}

// GenerateUsername generates a username from the Google name
func (gu *ClaimSet) GenerateUsername(authCache *redis.Client) (string, error) {
	// Trim non-username characters from Name
	username := utils.UsernameInverseRegex.ReplaceAllString(gu.Name, "")
	if utils.UsernameRegex.MatchString(username) {
		err := utils.ValidateEmail(username)
		if err != nil {
			return "", err
		}
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

	return "", errors.New("Cannot generate username from google name")
}

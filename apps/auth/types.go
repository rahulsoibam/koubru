package auth

import (
	"errors"
	"strings"

	"github.com/rahulsoibam/koubru-prod-api/utils"
)

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
	Picture  string `json:"picture"`
	Password string `json:"password"`
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

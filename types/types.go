package types

import "errors"

import "time"

// User type stores the basic overview info of the user
type User struct {
	Username    string `json:"username"`
	FullName    string `json:"full_name"`
	PhotoURL    string `json:"picture"`
	Bio         string `json:"bio,omitempty"`
	IsFollowing bool   `json:"is_following"`
	Counts      struct {
		Followers int64 `json:"followers"`
		Following int64 `json:"following"`
		Topics    int64 `json:"topics"`
		Opinions  int64 `json:"opinions"`
	} `json:"counts"`
}

// FollowUser type stores the data pertaining to the follower and the following view
type FollowUser struct {
	Username    string    `json:"username"`
	FullName    string    `json:"full_name"`
	PhotoURL    string    `json:"picture"`
	FollowedOn  time.Time `json:"followed_on"`
	IsFollowing bool      `json:"is_following"`
}

type Topic struct {
	ID          int64      `json:"id"`
	Title       string     `json:"title"`
	Details     string     `json:"details"`
	CreatedOn   time.Time  `json:"created_on"`
	CreatedBy   TopicUser  `json:"created_by"`
	Categories  []Category `json:"categories"`
	IsFollowing bool       `json:"is_following"`
	Counts      struct {
		Followers int64 `json:"followers"`
		Opinions  int64 `json:"opinions"`
	}
}

type TopicUser struct {
	Username string `json:"username"`
	FullName string `json:"full_name"`
	Picture  string `json:"picture"`
}

type Category struct {
	ID   int64  `json:"id"`
	Name string `json:"name"`
}

func (c *Category) Validate() error {
	// TODO
	if len(c.Name) < 3 || len(c.Name) > 32 {
		return errors.New("Should be less than 32 and more than 3 characters")
	}
	return nil
}

package types

import (
	"encoding/json"
	"errors"
	"time"
)

// User type stores the basic overview info of the user
type User struct {
	Username    string `json:"username"`
	FullName    string `json:"full_name"`
	PhotoURL    string `json:"picture"`
	Bio         string `json:"bio"`
	IsSelf      bool   `json:"is_self"`
	IsFollowing bool   `json:"is_following"`
	Counts      struct {
		Followers int64 `json:"followers"`
		Following int64 `json:"following"`
		Topics    int64 `json:"topics"`
		Opinions  int64 `json:"opinions"`
	} `json:"counts"`
}

type User_ struct {
	Username    string `json:"username"`
	FullName    string `json:"full_name"`
	Picture     string `json:"picture"`
	IsFollowing bool   `json:"is_following"`
	IsSelf      bool   `json:"is_self"`
}

// FollowUser type stores the data pertaining to the follower and the following view
type FollowUser struct {
	Username   string    `json:"username"`
	FullName   string    `json:"full_name"`
	PhotoURL   string    `json:"picture"`
	FollowedOn time.Time `json:"followed_on"`
	// IsFollower bool      `json:"is_follower"`
}

type Category struct {
	ID          int64     `json:"id"`
	Name        string    `json:"name"`
	CreatedOn   time.Time `json:"created_on"`
	CreatedBy   User_     `json:"created_by"`
	IsFollowing bool      `json:"is_following"`
	Counts      struct {
		Followers int64 `json:"followers"`
		Topics    int64 `json:"topics"`
	} `json:"counts"`
}

type Topic_ struct {
	ID          int64           `json:"id"`
	Title       string          `json:"title"`
	Details     string          `json:"details"`
	Categories  json.RawMessage `json:"categories"`
	IsFollowing bool            `json:"is_following"`
	CreatedBy   struct {
		Username string `json:"username"`
		FullName string `json:"full_name"`
		Picture  string `json:"picture"`
	} `json:"created_by"`
	CreatedOn time.Time `json:"created_on"`
}

type Topic struct {
	ID          int64           `json:"id"`
	Title       string          `json:"title"`
	Details     string          `json:"details"`
	CreatedOn   time.Time       `json:"created_on"`
	CreatedBy   User_           `json:"created_by"`
	Categories  json.RawMessage `json:"categories"`
	IsFollowing bool            `json:"is_following"`
	Counts      struct {
		Followers int64 `json:"followers"`
		Opinions  int64 `json:"opinions"`
	}
}

type Opinion struct {
	ID        int64 `json:"id"`
	CreatedBy struct {
		Username string `json:"username"`
		FullName string `json:"full_name"`
		Picture  string `json:"picture"`
	} `json:"created_by"`
	Topic struct {
		ID          int64           `json:"id"`
		Title       string          `json:"title"`
		Details     string          `json:"details"`
		Categories  json.RawMessage `json:"categories"`
		IsFollowing bool            `json:"is_following"`
	} `json:"topic"`
	IsAnonymous  bool      `json:"is_anonymous"`
	Vote         string    `json:"vote"`
	IsFollowing  int64     `json:"is_following"`
	ThumbnailURL string    `json:"thumbnail_url"`
	HlsURL       string    `json:"hls_url"`  // Change to stream_src
	DashURL      string    `json:"dash_url"` // Change to alt_stream_src
	Reaction     string    `json:"reaction"`
	CreatedOn    time.Time `json:"created_on"`
	Counts       struct {
		Views     int64 `json:"views"`
		Upvotes   int64 `json:"upvotes"`
		Downvotes int64 `json:"downvotes"`
		Followers int64 `json:"followers"`
		Replies   int64 `json:"replies"`
	} `json:"counts"`
}

type Opinion_ struct {
	ID        int64 `json:"id"`
	CreatedBy struct {
		Username    string `json:"username"`
		FullName    string `json:"full_name"`
		Picture     string `json:"picture"`
		IsSelf      string `json:"is_self"`
		IsFollowing bool   `json:"is_following"`
	} `json:"created_by"`
	IsAnonymous  bool      `json:"is_anonymous"`
	ThumbnailURL string    `json:"thumbnail_url"`
	Reaction     string    `json:"reaction"`
	IsFollowing  bool      `json:"is_following"`
	CreatedOn    time.Time `json:"created_on"`
}

type Category_ struct {
	ID          int64  `json:"id"`
	Name        string `json:"name"`
	IsFollowing bool   `json:"is_following"`
}

type NewCategory struct {
	ID   int64  `json:"id"`
	Name string `json:"name"`
}

func (c *NewCategory) Validate() error {
	// TODO
	if len(c.Name) < 3 || len(c.Name) > 32 {
		return errors.New("Should be less than 32 and more than 3 characters")
	}
	return nil
}

type NewTopic struct {
	Title      string   `json:"title"`
	Details    string   `json:"details"`
	Categories [3]int64 `json:"categories"`
}

func (t *NewTopic) Validate() error {
	if len(t.Title) < 8 || len(t.Title) > 80 {
		return errors.New("Should be less than 80 and more than 8 characters")
	}
	if len(t.Details) > 1024 {
		return errors.New("Details should be less than 1024 characters")
	}
	return nil
}

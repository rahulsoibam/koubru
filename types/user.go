package types

import "time"

type User struct {
	Username    string `json:"username"`
	FullName    string `json:"full_name"`
	Picture     string `json:"picture"`
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

// Started using this implementation
type UserForCreatedBy struct {
	Username    string `json:"username"`
	FullName    string `json:"full_name"`
	Picture     string `json:"picture"`
	IsFollowing bool   `json:"is_following"`
	IsSelf      bool   `json:"is_self"`
}

type UserForFollowList struct {
	Username    string    `json:"username"`
	FullName    string    `json:"full_name"`
	Picture     string    `json:"picture"`
	FollowedOn  time.Time `json:"followed_on"`
	IsFollowing bool      `json:"is_following"`
	IsSelf      bool      `json:"is_self"`
}

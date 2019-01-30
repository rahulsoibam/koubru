package user

import "time"

// User type stores the basic overview info of the user
type User struct {
	ID           int64  `json:"id"`
	Username     string `json:"username"`
	EmailVerfied bool   `json:"email_verified"`
	FullName     string `json:"full_name"`
	Picture      string `json:"picture"`
	Bio          string `json:"bio,omitempty"`
	Counts       struct {
		Followers int64 `json:"followers"`
		Following int64 `json:"following"`
		Topics    int64 `json:"topics"`
		Opinions  int64 `json:"opinions"`
	} `json:"counts"`
}

// FollowUser type stores the data pertaining to the follower and the following view
type FollowUser struct {
	ID         int64     `json:"id"`
	Username   string    `json:"username"`
	FullName   string    `json:"full_name"`
	Picture    string    `json:"picture"`
	FollowedOn time.Time `json:"followed_on"`
}

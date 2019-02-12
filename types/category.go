package types

import "time"

type Category struct {
	ID          int64            `json:"id"`
	Name        string           `json:"name"`
	CreatedOn   time.Time        `json:"created_on"`
	CreatedBy   UserForCreatedBy `json:"created_by"`
	IsFollowing bool             `json:"is_following"`
	Counts      struct {
		Followers int64 `json:"followers"`
		Topics    int64 `json:"topics"`
	} `json:"counts"`
}

type CategoryForList struct {
	ID          int64  `json:"id"`
	Name        string `json:"name"`
	IsFollowing bool   `json:"is_following"`
}

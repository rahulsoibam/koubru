package types

import (
	"encoding/json"
	"time"
)

type Topic struct {
	ID          int64            `json:"id"`
	Title       string           `json:"title"`
	Details     string           `json:"details"`
	CreatedBy   UserForCreatedBy `json:"created_by"`
	Categories  json.RawMessage  `json:"categories"`
	IsFollowing bool             `json:"is_following"`
	CreatedOn   time.Time        `json:"created_on"`
	Counts      struct {
		Followers int64 `json:"followers"`
		Opinions  int64 `json:"opinions"`
	} `json:"counts"`
}

type TopicForList struct {
	ID          int64            `json:"id"`
	Title       string           `json:"title"`
	Details     string           `json:"details"`
	CreatedOn   time.Time        `json:"created_on"`
	Categories  json.RawMessage  `json:"categories"`
	IsFollowing bool             `json:"is_following"`
	CreatedBy   UserForCreatedBy `json:"created_by"`
}

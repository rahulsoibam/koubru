package types

import (
	"encoding/json"
	"time"
)

type Opinion struct {
	ID        int64            `json:"id"`
	ParentID  int64            `json:"parent_id"`
	CreatedBy UserForCreatedBy `json:"created_by"`
	Topic     struct {
		ID          int64           `json:"id"`
		Title       string          `json:"title"`
		Details     string          `json:"details"`
		Categories  json.RawMessage `json:"categories"`
		IsFollowing bool            `json:"is_following"`
	} `json:"topic"`
	IsAnonymous bool     `json:"is_anonymous"`
	IsFollowing bool     `json:"is_following"`
	Thumbnails  []string `json:"thumbnails"`
	Sources     struct {
		Hls  string `json:"hls"`
		Dash string `json:"dash"`
		Aac  string `json:"aac"`
	} `json:"sources"`
	Vote      string    `json:"vote"`
	Reaction  string    `json:"reaction"`
	CreatedOn time.Time `json:"created_on"`
	Counts    struct {
		Views     int64 `json:"views"`
		Upvotes   int64 `json:"upvotes"`
		Downvotes int64 `json:"downvotes"`
		Followers int64 `json:"followers"`
		Replies   int64 `json:"replies"`
	} `json:"counts"`
}

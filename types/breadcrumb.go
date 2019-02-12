package types

import "time"

type Breadcrumb struct {
	OpinionID int64            `json:"opinion_id"`
	CreatedOn time.Time        `json:"created_on"`
	CreatedBy UserForCreatedBy `json:"created_by"`
	Counts    struct {
		Replies int64 `json:"replies"`
	} `json:"counts"`
}

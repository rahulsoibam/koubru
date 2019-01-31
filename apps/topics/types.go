package topics

import (
	"errors"
	"time"
)

// // Topic to store internal topic ctx
// type Topic struct {
// 	ID        int64  `json:"id"`
// 	Title     string `json:"title"`
// 	CreatedBy struct {
// 		ID   int64  `json:"created_by"`
// 		Name string `json:"name"`
// 	}
// 	Counts struct {
// 		Followers int64 `json:"followers"`
// 	} `json:"counts"`
// }

type Topic struct {
	ID          int64      `json:"id"`
	Title       string     `json:"title"`
	Details     string     `json:"details"`
	CreatedOn   time.Time  `json:"created_on"`
	CreatedBy   User       `json:"created_by"`
	Categories  []Category `json:"categories"`
	IsFollowing bool       `json:"is_following"`
	Counts      struct {
		Followers int64 `json:"followers"`
		Opinions  int64 `json:"opinions"`
	} `json:"counts"`
}

type User struct {
	Username string `json:"username"`
	FullName string `json:"full_name"`
	Picture  string `json:"picture"`
}

type Category struct {
	ID   int64  `json:"id"`
	Name string `json:"name"`
}

type NewTopic struct {
	Title      string     `json:"title"`
	Details    string     `json:"details"`
	Categories []Category `json:"categories"`
}

func (nt *NewTopic) Validate() error {
	if len(nt.Title) < 16 {
		return errors.New("Topic title too short. Should be more than 16 characters")
	}
	if len(nt.Title) > 80 {
		return errors.New("Topic title too long. Should be less than 80 characters")
	}
	if len(nt.Details) > 1024 {
		return errors.New("Topic details too long. Should be less than 1024 characters")
	}
	return nil
}

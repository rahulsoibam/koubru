package topics

import "time"

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
	ID         int64      `json:"id"`
	Title      string     `json:"title"`
	Details    string     `json:"details"`
	CreatedOn  time.Time  `json:"created_by"`
	CreatedBy  User       `json:"created_by"`
	Categories []Category `json:"categories"`
}

type User struct {
	ID       int64  `json:"id"`
	Username string `json:"username"`
	FullName string `json:"full_name"`
}

type Category struct {
	ID   int64  `json:"id"`
	Name string `json:"name"`
}

package types

type SearchTopic struct {
	ID    int64  `json:"id"`
	Title string `json:"title"`
}

// Temporary type
type SearchUser struct {
	Username string `json:"username"`
	FullName string `json:"full_name"`
	Picture  string `json:"picture"`
}

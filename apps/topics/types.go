package topics

// Topic to store internal topic ctx
type Topic struct {
	ID int64 `json:"id"`
	Title string `json:"title"`
	CreatedBy struct {
		ID int64 `json:"created_by"`
		Name string `json:"name"`
	}
	Counts struct {
		Followers int64 `json:"followers"`
	} `json:"counts"`
}

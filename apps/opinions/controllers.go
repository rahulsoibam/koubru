package opinions

import (
	"net/http"
)

// List to list all opinions
func (a App) List(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("List opinions"))
}

// Create to create an opinion
func (a App) Create(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Create opinion"))
}

// Get to get details of an opinion
func (a App) Get(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Get opinion"))
}

// Delete to delete an opinion
func (a App) Delete(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Delete opinion"))
}

// Followers to get followers of an opinion
func (a App) Followers(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Get followers of opinion"))
}

// Follow to follow an opinion
func (a App) Follow(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Follow opinion"))
}

// Unfollow to unfollow an opinion
func (a App) Unfollow(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Unfollow opinion"))
}

// Replies to reply to an opinion
func (a App) Replies(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Get replies of opinion"))
}

// Report to report an opinion
func (a App) Report(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Report opinion"))
}

// Vote to vote on an opinion
func (a App) Vote(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Vote on an opinion"))
}

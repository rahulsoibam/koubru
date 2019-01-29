package users

import "net/http"

// List to list all opinions of a user
func (a App) List(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("List all users"))
}

// Get to get details of a user
func (a App) Get(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Get details of a user"))
}

// Followers to list followers of a user
func (a App) Followers(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("List followers of a user"))
}

// Following to list users followed by a user
func (a App) Following(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("List users a user if following"))
}

// Follow to follow a user
func (a App) Follow(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Follow a user"))
}

// Unfollow to unfollow a user
func (a App) Unfollow(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Unfollow a user"))
}

// Topics to list topics by a user
func (a App) Topics(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("List topics of a user"))
}

// Opinions to list opinions of a user
func (a App) Opinions(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("List opinions of a user"))
}

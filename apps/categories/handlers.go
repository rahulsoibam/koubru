package categories

import (
	"net/http"
)

// List all categories
func (a App) List(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("list all categories"))
}

// Create a category
func (a App) Create(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("create a category"))
}

// Get details of a category
func (a App) Get(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("get details of a category"))
}

// Follow to follow a category
func (a App) Follow(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Follow a category"))
}

// Unfollow to unfollow a category
func (a App) Unfollow(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Unfollow a category"))
}

// BulkFollow to follow many categories at once
func (a App) BulkFollow(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("bulk follow, first app entry"))
}

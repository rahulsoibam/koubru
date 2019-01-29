package search

import "net/http"

// Users - search users
func (a App) Users(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("search user"))
}

// Topics - search topics
func (a App) Topics(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("search topics"))
}

// Categories - search categpories
func (a App) Categories(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("search categories"))
}

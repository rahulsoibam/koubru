package user

import "net/http"

// Get details of authenticated user
func (a App) Get(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("get authenticated user"))
}

// Patch details of authenticated user
func (a App) Patch(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("path authenticated user"))
}

// Delete or deactivate authenticated user
func (a App) Delete(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("delete/deactivate user"))
}

// Followers to list followers to authenticated user
func (a App) Followers(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("get followers of authenticated user"))
}

// Following to list users whom the authenticated user is following
func (a App) Following(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("get users whom the authenticated user is following"))
}

// Opinions of authenticated user
func (a App) Opinions(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("list opinions of authenticated user"))
}

// Topics of authenticated user
func (a App) Topics(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("list topics of authenticated user"))
}

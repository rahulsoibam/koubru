package feed

import "net/http"

// Get feed
func (a App) Get(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Feed"))
}

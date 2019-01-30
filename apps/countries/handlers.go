package countries

import "net/http"

// List countries
func (a *App) List(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("List all countries"))
}

// BulkSelect countries
func (a *App) BulkSelect(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Bulk select countries"))
}

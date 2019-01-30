package explore

import "net/http"

// Get the explore page
func (a App) Get(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("list contents of explore page"))
}

// // Nearby to get explore page content related to nearby locations
// func (a App) Nearby(w http.ResponseWriter, r *http.Request) {
// 	w.Write([]byte("list nearby happennings"))
// }

// // Location to get explore page content related to location
// func (a App) Location(w http.ResponseWriter, r *http.Request) {
// 	w.Write([]byte("explore page tailored to set location"))
// }

package router

import "github.com/gorilla/mux"

// constant names for the named routes
const (
	NewsOverview = "News:overview"
	NewsPost     = "News:post"
)

// NewsApp constructs a mux.Router containing the routes for News html app
func NewsApp(m *mux.Router) *mux.Router {
	if m == nil {
		m = mux.NewRouter()
	}

	m.Path("/").Methods("GET").Name(NewsOverview)
	m.Path("/News/{PostID:[0-9]+}").Methods("GET").Name(NewsPost)

	return m
}

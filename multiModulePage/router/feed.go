package router

import "github.com/gorilla/mux"

// constant names for the named routes
const (
	FeedOverview = "Feed:overview"
	FeedPost     = "Feed:post"
)

// FeedApp constructs a mux.Router containing the routes for Feed html app
func FeedApp(m *mux.Router) *mux.Router {
	if m == nil {
		m = mux.NewRouter()
	}

	m.Path("/").Methods("GET").Name(FeedOverview)
	m.Path("/feed/{PostID:[0-9]+}").Methods("GET").Name(FeedPost)

	return m
}

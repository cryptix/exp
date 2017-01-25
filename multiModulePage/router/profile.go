package router

import "github.com/gorilla/mux"

// constant names for the named routes
const (
	ProfileMy   = "Profile:my"
	ProfileUser = "Profile:User"
)

// ProfileApp constructs a mux.Router containing the routes for Profile html app
func ProfileApp(m *mux.Router) *mux.Router {
	if m == nil {
		m = mux.NewRouter()
	}

	m.Path("/").Methods("GET").Name(ProfileMy)
	m.Path("/Profile/{UserID:[0-9]+}").Methods("GET").Name(ProfileUser)

	return m
}

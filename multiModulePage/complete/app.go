// Package complete uses all of the modules: feed, news, profile
//  /feed - package feed
//	/news 	- package news
//	/profile	- package profile
//	/assets		- static http files from bindata
package complete

import (
	"net/http"

	"github.com/cryptix/exp/multiModulePage"
	"github.com/cryptix/exp/multiModulePage/router"
	"github.com/cryptix/go/http/render"
	"github.com/gorilla/mux"
)

func init() {
	render.Init(multiModulePage.Assets, []string{"/tisDaemon/navbar.tmpl", "/tisDaemon/base.tmpl"})
	render.AddTemplates([]string{
		"/tisDaemon/index.tmpl",
		"/tisDaemon/info/contact.tmpl",
		"/tisDaemon/info/license.tmpl",
		"/about.tmpl",
		"/error.tmpl",
	})
}

// Handler creates a full fledged http handler for the TIS Daemon app
func Handler(m *mux.Router) (http.Handler, error) {
	if m == nil {
		m = router.CompleteApp()
	}

	// javascript, images, ...
	//	m.PathPrefix("/assets/").Handler(http.StripPrefix("/assets/", http.FileServer(assets)))

	m.PathPrefix("/feed").Handler(http.StripPrefix("/feed", feed.Handler(m)))
	m.PathPrefix("/news").Handler(http.StripPrefix("/news", news.Handler(m)))
	m.PathPrefix("/profile").Handler(http.StripPrefix("/profile", profile.Handler(m)))

	m.Get(router.CompleteIndex).Handler(render.StaticHTML("/complete/index.tmpl"))
	m.Get(router.CompleteAbout).Handler(render.StaticHTML("/complete/about.tmpl"))

	return m, nil
}

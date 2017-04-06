// Package complete uses all of the modules: feed, news, profile
//  /feed - package feed
//	/news 	- package news
//	/profile	- package profile
//	/assets		- static http files from bindata
package complete

import (
	"html/template"
	"net/http"

	"github.com/cryptix/go/http/render"
	"github.com/gorilla/mux"
	"gopkg.in/errgo.v1"

	"github.com/cryptix/exp/multiModulePage"
	"github.com/cryptix/exp/multiModulePage/feed"
	"github.com/cryptix/exp/multiModulePage/router"
)

// Handler creates a full fledged http handler for the TIS Daemon app
func Handler(m *mux.Router) (http.Handler, error) {
	if m == nil {
		m = router.CompleteApp()
	}
	r, err := render.New(multiModulePage.Assets,
		render.BaseTemplates("/complete/base.tmpl"),
		render.AddTemplates(append(feed.HTMLTemplates,
			"/complete/index.tmpl",
			"/complete/about.tmpl",
			"/error.tmpl")...),
		render.FuncMap(template.FuncMap{
			"urlTo": multiModulePage.NewURLTo(m),
		}),
	)
	if err != nil {
		return nil, errgo.Notef(err, "complete.Handler: failed to create renderer")
	}

	// javascript, images, ...
	//	m.PathPrefix("/assets/").Handler(http.StripPrefix("/assets/", http.FileServer(assets)))

	m.PathPrefix("/feed").Handler(http.StripPrefix("/feed", feed.Handler(m, r)))
	// m.PathPrefix("/news").Handler(http.StripPrefix("/news", news.Handler(m)))
	// m.PathPrefix("/profile").Handler(http.StripPrefix("/profile", profile.Handler(m)))

	m.Get(router.CompleteIndex).Handler(r.StaticHTML("/complete/index.tmpl"))
	m.Get(router.CompleteAbout).Handler(r.StaticHTML("/complete/about.tmpl"))

	return m, nil
}

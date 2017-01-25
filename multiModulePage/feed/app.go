package feed

import (
	"net/http"

	"github.com/cryptix/go/http/render"
	"github.com/cryptix/go/logging"
	"github.com/gorilla/mux"

	"github.com/cryptix/exp/multiModulePage/router"
)

var HTMLTemplates = []string{
	"/feed/overview.tmpl",
	"/feed/post.tmpl",
}

var l = logging.Logger("webApp/archive")
var r *render.Renderer

// Handler creates a http.Handler with all the archives routes attached to it
func Handler(m *mux.Router) http.Handler {
	if m == nil {
		m = router.FeedApp(nil)
	}

	m.Get(router.FeedOverview).Handler(r.HTML("/feed/overview.tmpl", showOverview))
	m.Get(router.FeedPost).Handler(r.HTML("/feed/post.tmpl", showPost))

	return m
}

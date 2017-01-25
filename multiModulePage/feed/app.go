package feed

import (
	"net/http"

	"github.com/cryptix/exp/multiModulePage/router"
	"github.com/cryptix/go/http/render"
	"github.com/cryptix/go/logging"
	"github.com/gorilla/mux"
)

func init() {
	render.AddTemplates([]string{
		"/feed/stepList.tmpl",
		"/feed/overview.tmpl",
	})
}

var l = logging.Logger("webApp/archive")

// Handler creates a http.Handler with all the archives routes attached to it
func Handler(m *mux.Router, apiurl string) http.Handler {
	if m == nil {
		m = router.FeedApp(nil)
	}

	m.Get(router.FeedOverview).Handler(render.HTML(showOverview))
	m.Get(router.FeedPost).Handler(render.HTML(showJob))

	return m
}

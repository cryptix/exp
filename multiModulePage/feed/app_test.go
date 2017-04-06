package feed

import (
	"html/template"
	"net/http"
	"testing"

	"github.com/cryptix/go/http/render"
	"github.com/cryptix/go/http/tester"
	"github.com/cryptix/go/logging/logtest"
	"github.com/pkg/errors"

	"github.com/cryptix/exp/multiModulePage"
	"github.com/cryptix/exp/multiModulePage/router"
)

var (
	testMux    *http.ServeMux
	testClient *tester.Tester
	testRouter = router.FeedApp(nil)
)

func setup(t *testing.T) {
	log := logtest.KitLogger("feed", t)
	r, err := render.New(multiModulePage.Assets,
		render.SetLogger(log),
		render.BaseTemplates("/testing/base.tmpl"),
		render.AddTemplates(append(HTMLTemplates, "/error.tmpl")...),
		render.FuncMap(template.FuncMap{
			"urlTo": multiModulePage.NewURLTo(testRouter),
		}),
	)
	if err != nil {
		t.Fatal(errors.Wrap(err, "setup: render init failed"))
	}
	testMux = http.NewServeMux()
	testMux.Handle("/", Handler(testRouter, r))
	testClient = tester.New(testMux, t)
}

func teardown() {
	testMux = nil
	testClient = nil
}

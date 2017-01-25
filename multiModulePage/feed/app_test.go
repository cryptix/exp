package feed

import (
	"html/template"
	"net/http"
	"testing"

	"github.com/cryptix/go/http/render"
	"github.com/cryptix/go/http/tester"
	"gopkg.in/errgo.v1"

	"github.com/cryptix/exp/multiModulePage"
	"github.com/cryptix/exp/multiModulePage/router"
)

var (
	testMux    *http.ServeMux
	testClient *tester.Tester
	testRouter = router.FeedApp(nil)
)

func setup(t *testing.T) {
	var err error
	r, err = render.New(multiModulePage.Assets,
		render.BaseTemplate("/testing/base.tmpl"),
		render.AddTemplates(append(HTMLTemplates, "/error.tmpl")...),
		render.FuncMap(template.FuncMap{
			"urlTo": multiModulePage.NewURLTo(testRouter),
		}),
	)
	if err != nil {
		t.Fatal(errgo.Notef(err, "setup: render init failed"))
	}
	testMux = http.NewServeMux()
	testMux.Handle("/", Handler(testRouter))
	testClient = tester.New(testMux, t)
}

func teardown() {
	r = nil
	testMux = nil
	testClient = nil
}

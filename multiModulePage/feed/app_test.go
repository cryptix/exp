package feed

import (
	"net/http"
	"testing"

	"github.com/cryptix/go/http/render"
	"github.com/cryptix/go/http/tester"

	"DST/tisdb/mockdb"
	"DST/webApp"
	"DST/webClient"
	"DST/webRouter"
)

var (
	testMux    *http.ServeMux
	testClient *tester.Tester
	testRouter = webRouter.ArchiveApp(nil)
	fakeA      *mockdb.Archive
)

func init() {
	render.Init(webApp.Assets, []string{"/testing/navbar.tmpl", "/testing/base.tmpl"})
	render.AddTemplates([]string{"/error.tmpl"})
	render.SetAppRouter(testRouter)
	render.Load()
}

func setup(t *testing.T) {
	testMux = http.NewServeMux()
	testMux.Handle("/", Handler(testRouter, ""))
	testClient = tester.New(testMux, t)
	fakeA = new(mockdb.Archive)
	apiclient = &webClient.Client{Archive: fakeA}
}

func teardown() {
	apiclient = nil
	testMux = nil
	testClient = nil
	fakeA = nil
}

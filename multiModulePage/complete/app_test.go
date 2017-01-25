package complete

import (
	"net/http"
	"testing"

	"github.com/cryptix/exp/multiModulePage/router"
	"github.com/cryptix/go/http/render"
	"github.com/cryptix/go/http/tester"
)

var (
	testMux    *http.ServeMux
	testClient *tester.Tester
	testRouter = router.CompleteApp()
)

func init() {
	render.SetAppRouter(testRouter)
	render.Load()
}

func setup(t *testing.T) {
	testMux = http.NewServeMux()
	h, err := Handler(testRouter, "")
	if err != nil {
		t.Fatalf("handler setup failed: %s", err)
	}
	testMux.Handle("/", h)
	testClient = tester.New(testMux, t)
}

func teardown() {
	testMux = nil
	testClient = nil
}

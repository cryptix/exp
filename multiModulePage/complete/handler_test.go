package complete

import (
	"net/http"
	"testing"

	"github.com/cryptix/exp/multiModulePage/router"
	"github.com/stretchr/testify/assert"
)

func TestURLTo_index(t *testing.T) {
	setup(t)
	defer teardown()
	a := assert.New(t)
	url, err := router.CompleteApp().Get(router.CompleteIndex).URL()
	a.Nil(err)
	html, resp := testClient.GetHTML(url.String(), nil)
	a.Equal(http.StatusOK, resp.Code, "wrong HTTP status code")
	a.Equal("<title>Complete - Index", html.Find("title").Text())
}

func TestURLTo_complete(t *testing.T) {
	setup(t)
	defer teardown()
	a := assert.New(t)
	url, err := router.CompleteApp().Get(router.FeedPost).URL("PostID", "1")
	a.Nil(err)
	html, resp := testClient.GetHTML(url.String(), nil)
	a.Equal(http.StatusOK, resp.Code, "wrong HTTP status code")

	lnk, ok := html.Find("#overview").Attr("href")
	a.True(ok)
	a.Equal("/feed/", lnk)
	lnk, ok = html.Find("#next").Attr("href")
	a.True(ok)
	a.Equal("/feed/post/3", lnk)
}

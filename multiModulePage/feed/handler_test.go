package feed

import (
	"net/http"
	"testing"

	"github.com/cryptix/exp/multiModulePage/router"
	"github.com/stretchr/testify/assert"
)

func TestOverview(t *testing.T) {
	setup(t)
	defer teardown()
	a := assert.New(t)
	url, err := router.FeedApp(nil).Get(router.FeedOverview).URL()
	a.Nil(err)
	html, resp := testClient.GetHTML(url.String(), nil)
	a.Equal(http.StatusOK, resp.Code, "wrong HTTP status code")
	a.Equal(html.Find("#welcome").Text(), "Welcome!")
}

func TestPost(t *testing.T) {
	setup(t)
	defer teardown()
	a := assert.New(t)
	url, err := router.FeedApp(nil).Get(router.FeedPost).URL("PostID", "1")
	a.Nil(err)
	html, resp := testClient.GetHTML(url.String(), nil)
	a.Equal(http.StatusOK, resp.Code, "wrong HTTP status code")
	a.Equal(html.Find("h1").Text(), db[1].Name)
}

func TestURLTo(t *testing.T) {
	setup(t)
	defer teardown()
	a := assert.New(t)
	url, err := router.FeedApp(nil).Get(router.FeedPost).URL("PostID", "1")
	a.Nil(err)
	html, resp := testClient.GetHTML(url.String(), nil)
	a.Equal(http.StatusOK, resp.Code, "wrong HTTP status code")
	a.Equal(html.Find("h1").Text(), db[1].Name)
	lnk, ok := html.Find("#overview").Attr("href")
	a.True(ok)
	a.Equal("/", lnk)
	lnk, ok = html.Find("#next").Attr("href")
	a.True(ok)
	a.Equal("/feed/3", lnk)
}

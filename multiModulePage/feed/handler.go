package feed

import (
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"gopkg.in/errgo.v1"
)

type Post struct {
	Name, Text string
}

var db = []Post{
	Post{"Hello", "lot's of stuff"},
	Post{"Testing", "yeeeeaaaahhhh..."},
	Post{"WAT", "i have only a partial idea of what i'm doing"},
}

func showOverview(w http.ResponseWriter, req *http.Request) (interface{}, error) {
	return db, nil
}

func showPost(w http.ResponseWriter, req *http.Request) (interface{}, error) {
	i, err := strconv.Atoi(mux.Vars(req)["PostID"])
	if err != nil {
		return nil, errgo.Notef(err, "argument parsing failed")
	}
	if i < 0 || i >= len(db) {
		return nil, errgo.Newf("db limit exceeded")
	}
	return db[i], nil
}

// +build js

package main

import (
	"net"
	"net/rpc"

	"github.com/cryptix/exp/humbleRpcTodo/frontend/views"
	"github.com/gopherjs/websocket"
	"github.com/soroushjp/humble/router"
	"github.com/soroushjp/humble/view"
	"honnef.co/go/js/console"
	"honnef.co/go/js/dom"
)

const (
	appDivSelector = "#app"
)

var (
	doc      = dom.GetWindow().Document()
	elements = struct {
		body      dom.Element
		todoList  dom.Element
		newTodo   dom.Element
		toggleBtn dom.Element
	}{}
	appHasLoaded = false
)

func init() {
	elements.body = doc.QuerySelector(appDivSelector)
}

var (
	conn   net.Conn
	client *rpc.Client
)

func main() {
	console.Log("Starting...")

	var err error
	conn, err = websocket.Dial("ws://localhost:3000/rpc-websocket")
	if err != nil {
		console.Error("Dial faild", err)
		return
	}
	client = rpc.NewClient(conn)

	console.Log("dialed...")

	//Start main app view, appView
	appView := &views.App{
		Client: client,
	}
	if err := view.ReplaceParentHTML(appView, appDivSelector); err != nil {
		panic(err)
	}

	r := router.New()
	r.HandleFunc("/", func(params map[string]string) {
		appView.InitChildren()
		if err := view.Update(appView); err != nil {
			panic(err)
		}
		if err := view.Update(appView.Footer); err != nil {
			panic(err)
		}
		appView.ApplyFilter(views.FilterAll)
	})

	r.HandleFunc("/active", func(params map[string]string) {
		appView.InitChildren()
		if err := view.Update(appView); err != nil {
			panic(err)
		}
		if err := view.Update(appView.Footer); err != nil {
			panic(err)
		}
		appView.ApplyFilter(views.FilterActive)
	})

	r.HandleFunc("/completed", func(params map[string]string) {
		appView.InitChildren()
		if err := view.Update(appView); err != nil {
			panic(err)
		}
		if err := view.Update(appView.Footer); err != nil {
			panic(err)
		}
		appView.ApplyFilter(views.FilterCompleted)
	})

	r.HandleFunc("/completed", func(params map[string]string) {
		console.Log("At Completed")
	})

	r.Start()

}

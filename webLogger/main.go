// webLogger shows an example usage of logging.GetHTTPHandler
// it exposes log events over websocket using my wslog package
package main

import (
	"net/http"
	"os"
	"time"

	"github.com/cryptix/go/logging/wslog"
	"github.com/op/go-logging"
)

var l = logging.MustGetLogger("example")

func main() {
	i := 0

	stderrbackend := logging.NewLogBackend(os.Stderr, "stderr:", 0)

	wsbackend := wslog.NewBackend()
	var format = logging.MustStringFormatter(
		"%{longfile} ▶ %{id:03x} %{message}",
	)

	wsFormatter := logging.NewBackendFormatter(wsbackend, format)

	logging.SetBackend(stderrbackend, wsFormatter)

	go func() {
		for {
			l.Warning("Logging event %d", i)
			l.Debug("Hello")
			i++
			time.Sleep(time.Second)
		}
	}()

	http.Handle("/log", wsbackend)
	check(http.ListenAndServe(":8080", nil))
}

func check(err error) {
	if err != nil {
		panic(err)
	}
}

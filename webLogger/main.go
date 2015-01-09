// webLogger shows an example usage of logging.GetHTTPHandler
// it exposes log events through http
package main

import (
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/cryptix/go/logging"
)

var l = logging.Logger("webLogger")

func main() {
	i := 0

	fileLogOutput, err := os.Create(filepath.Base(os.Args[0]) + ".log")
	logging.CheckFatal(err)
	logging.SetupLogging(fileLogOutput)

	go func() {
		time.Sleep(time.Second * 10)
		fileLogOutput.Close()
		os.Exit(0)
	}()

	go func() {
		for {
			l.Warningf("Logging event %d", i)
			l.Debug("Hello")
			i++
			time.Sleep(time.Second)
		}
	}()

	http.Handle("/log", logging.GetHTTPHandler())
	logging.CheckFatal(http.ListenAndServe(":8080", nil))
}

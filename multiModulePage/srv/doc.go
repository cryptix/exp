// server for the 'complete' page
package main

import (
	stdlog "log"
	"net"
	"net/http"
	"os"

	"github.com/cryptix/go/logging"

	"github.com/cryptix/exp/multiModulePage/complete"
)

func main() {
	logging.SetupLogging(nil)
	var mylog = logging.Logger("websrv")

	stdlog.Printf("Hello wold")
	mylog.Log("hello", "mylog")

	h, err := complete.Handler(nil)
	logging.CheckFatal(err)

	addr := os.Args[1]
	if addr == "" {
		addr = "[::]:0"
	}
	lis, err := net.Listen("tcp", addr)
	logging.CheckFatal(err)
	mylog.Log("msg", "http listening", "addr", lis.Addr().String())

	err = http.Serve(lis, h)
	logging.CheckFatal(err)
}

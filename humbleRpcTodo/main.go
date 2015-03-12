package main

import (
	"flag"
	"log"
	"net"
	"net/http"
	"net/rpc"
	"os"

	"github.com/codegangsta/negroni"
	"github.com/elazarl/go-bindata-assetfs"
	"golang.org/x/net/websocket"
)

var (
	host = flag.String("host", "localhost", "The hostname/ip to listen on.")
	port = flag.String("port", "0", "The port number to listen on.")
)

//go:generate gopherjs build -m -o public/js/app.js github.com/cryptix/exp/humbleRpcTodo/frontend
//go:generate go-bindata -prefix=public public/...

func main() {
	flag.Parse()

	mux := http.NewServeMux()

	// Rpc!
	rpc.Register(&TodoService{})
	mux.Handle("/rpc-websocket", websocket.Handler(func(conn *websocket.Conn) {
		conn.PayloadType = websocket.BinaryFrame
		rpc.ServeConn(conn)
	}))

	mux.Handle("/", http.FileServer(&assetfs.AssetFS{
		Asset:    Asset,
		AssetDir: AssetDir,
	}))

	n := negroni.New()
	n.Use(negroni.NewRecovery())
	n.Use(negroni.NewLogger())
	n.UseHandler(mux)

	// Start the server
	if *port == "0" {
		*port = os.Getenv("PORT")
	}

	l, err := net.Listen("tcp", *host+":"+*port)
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("Serving at http://%s/", l.Addr())

	if err := http.Serve(l, n); err != nil {
		log.Fatal(err)
	}
}

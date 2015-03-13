// Try net/rpc/jsonrpc between backend and frontend (via GopherJS) through a websocket connection.
//
// based on https://github.com/shurcooL/play/tree/master/42
package main

import (
	"fmt"
	"io"
	"net/http"
	"net/rpc"
	"net/rpc/jsonrpc"

	"github.com/shurcooL/go/gopherjs_http"
	"golang.org/x/net/websocket"
)

type Args struct {
	A, B int
}

type Arith struct {
	req       *http.Request
	connClose func() error
}

func (a *Arith) Multiply(args *Args, reply *int) error {
	fmt.Println("from", a.req.RemoteAddr)
	c, err := a.req.Cookie("AwesomeRPC")
	if err != nil {
		return err
	}
	fmt.Println("cookie:", c)

	*reply = args.A * args.B
	fmt.Printf("locally multiplying %v by %v -> %v\n", args.A, args.B, *reply)
	return nil
}

func main() {
	rpc.Register(&Arith{})

	http.Handle("/rpc-websocket", websocket.Handler(func(conn *websocket.Conn) {
		s := rpc.NewServer()
		a := &Arith{
			req:       conn.Request(),
			connClose: conn.Close,
		}
		s.Register(a)
		s.ServeCodec(jsonrpc.NewServerCodec(conn))
	}))

	http.HandleFunc("/index.html", func(w http.ResponseWriter, req *http.Request) {
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		io.WriteString(w, `<html>
	<head></head>
	<body>
		<pre id="output"></pre>
		<script type="text/javascript" src="/script.js"></script>
	</body>
</html>
`)
	})
	http.Handle("/script.js", gopherjs_http.GoFiles("./script.go"))

	err := http.ListenAndServe(":8880", nil)
	if err != nil {
		panic(err)
	}
}

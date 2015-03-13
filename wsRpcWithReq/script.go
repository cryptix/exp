// +build js

package main

import (
	"fmt"
	"net"
	"net/rpc/jsonrpc"
	"time"

	"github.com/gopherjs/websocket"
	"honnef.co/go/js/dom"
)

var document = dom.GetWindow().Document()

var output = document.GetElementByID("output").(*dom.HTMLPreElement)

func appendOutput(s string) {
	output.SetTextContent(output.TextContent() + s)
}

type Args struct {
	A, B int
}

var conn net.Conn

func main() {
	var err error
	conn, err = websocket.Dial("ws://localhost:8880/rpc-websocket")
	if err != nil {
		panic(err)
	}

	client := jsonrpc.NewClient(conn)

	for i := 0; i < 100; i++ {

		started := time.Now()
		args := &Args{15, 3}
		var reply int
		{
			err := client.Call("Arith.Multiply", args, &reply)
			if err != nil {
				fmt.Println("arith error:", err)
			}
		}
		appendOutput(fmt.Sprintf("Arith: %d*%d=%d taken %v\n", args.A, args.B, reply, time.Since(started).String()))
		time.Sleep(5 * time.Second)
	}
	err = client.Close()
	if err != nil {
		fmt.Println("client.Close():", err)
	}

	err = conn.Close()
	if err != nil {
		fmt.Println("conn.Close():", err)
	}
}

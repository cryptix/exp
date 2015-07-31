package main

import (
	"flag"
	"fmt"
	"log"
	"net/url"
	"os"
	"strings"

	"golang.org/x/net/context"

	client "github.com/cryptix/exp/todoKitSvc/client"
	httpclient "github.com/cryptix/exp/todoKitSvc/client/http"
	"github.com/cryptix/exp/todoKitSvc/todosvc"
)

func main() {
	fs := flag.NewFlagSet("todocli", flag.ExitOnError)
	var (
		transport = fs.String("transport", "http", "http, netrpc")
		httpAddr  = fs.String("http.addr", "localhost:8001", "HTTP (JSON) address")
		//netrpcAddr = fs.String("netrpc.addr", "localhost:8003", "net/rpc address")
	)
	flag.Usage = fs.Usage
	fs.Parse(os.Args[1:])
	log.SetFlags(0)
	log.SetOutput(os.Stdout)

	var e todosvc.Endpoints
	switch *transport {
	case "http":
		if !strings.HasPrefix(*httpAddr, "http") {
			*httpAddr = "http://" + *httpAddr
		}
		u, err := url.Parse(*httpAddr)
		if err != nil {
			log.Fatalf("url.Parse: %v", err)
		}
		e = httpclient.NewClient("GET", u.String())

	case "netrpc":
		log.Fatalf("unsupported transport %q", *transport)
		//client, err := rpc.DialHTTP("tcp", *netrpcAddr)
		//if err != nil {
		//	log.Fatalf("rpc.DialHTTP: %v", err)
		//}
		//e = netrpcclient.NewClient(client)

	default:
		log.Fatalf("unsupported transport %q", *transport)
	}

	c := client.NewClient(e)
	ctx := context.Background()

	args := fs.Args()
	if len(args) < 1 {
		fmt.Fprintln(os.Stderr, "usage: todocli <add|list> ...")
		os.Exit(1)
	}
	switch args[0] {
	case "add":
		log.Println("Adding")
		id, err := c.Add(ctx, args[1])
		check(err)
		log.Println("Added as", id)
	case "list":
		log.Println("Listing")
		list, err := c.List(ctx)
		check(err)
		fmt.Println(list)
	default:
		log.Fatalf("unsupported command: %s", args[0])
	}
}

func check(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

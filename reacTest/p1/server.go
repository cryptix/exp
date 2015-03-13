package main

import (
	"log"
	"net"
	"net/http"
	"os"

	"github.com/codegangsta/cli"
)

//go:generate -command asset go run asset.go
//go:generate asset index.html
//go:generate asset bundle.js

type Binary struct {
	asset
}

func html(a asset) Binary {
	return Binary{a}
}

func js(a asset) Binary {
	return Binary{a}
}

func main() {
	app := cli.NewApp()
	app.Name = "p1"
	app.Action = run

	app.Flags = []cli.Flag{
		cli.StringFlag{Name: "port,p", Value: "0", EnvVar: "PORT"},
		cli.StringFlag{Name: "host", Value: "localhost"},
	}
	app.Run(os.Args)

}

func run(ctx *cli.Context) {

	l, err := net.Listen("tcp", ctx.String("host")+":"+ctx.String("port"))
	check(err)
	log.Printf("Serving at http://%s/", l.Addr())
	http.Handle("/", index)
	http.Handle("/bundle.js", bundle)
	check(http.Serve(l, nil))
}

func check(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

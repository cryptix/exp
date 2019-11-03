package main

import (
	"fmt"
	"github.com/cryptix/exp/unixsockmux/client"
	"github.com/cryptix/go/logging"
	"io"
	"log"
	"os"
	"time"
)

var check = logging.CheckFatal

const path = "/tmp/unixsockmux.sock"

func main() {

	c, err := client.NewClient(path)
	check(err)

	rwc, err := c.OpenChannel("ping", struct{ Foo int }{23})
	check(err)

	go func() {
		for i:=5;i>0;i--{
			fmt.Fprintln(rwc, "ping",i)
			time.Sleep(250 * time.Millisecond)
		}
		rwc.Close()
	}()
	go func() {
		io.Copy(os.Stderr, rwc)
		log.Println("ping input closed")
	}()

	echo, err := c.OpenChannel("echo", 42)
	check(err)
	go func() {
		io.Copy(os.Stdout, echo)
	}()

	for i,s := range []string{"hello", "world", "sup?"} {
		fmt.Fprintln(echo,s,i)
		log.Println("echoed",s)
		time.Sleep(1*time.Second)
	}

	//io.Copy(echo, os.Stdin)
	log.Println("echo closed")
	echo.Close()
	c.Close()
}

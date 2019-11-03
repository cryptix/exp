package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"math/rand"
	"net"
	"os"
	"strings"
	"time"

	"github.com/cryptix/go/logging"
	"github.com/ftrvxmtrx/fd"
	"github.com/pkg/errors"
)

const path = "/tmp/unixsockmux.sock"

func main() {
	rand.Seed(time.Now().Unix())
	os.Remove(path)

	addr, err := net.ResolveUnixAddr("unix", path)
	check(err)

	l, err := net.ListenUnix("unix", addr)
	check(err)

	log.Println("Accepting unix sockets")

	done := make(chan struct{})
	errc := make(chan error)
	go func() {
		for {
			select {
			case e := <-errc:
				if e != nil {
					log.Printf("Error from errChan:\n%s\n", e)
					break
				}
			default:
				clientConn, err := l.AcceptUnix()
				if err != nil {
					log.Printf("AcceptUnix() error:%s\n", err)
					break
				}

				go handleConn(clientConn, errc)
			}

		}
		close(done)

	}()
	<-done

	log.Println("Accept loop closed.")

}

var check = logging.CheckFatal

func handleConn(c *net.UnixConn, errc chan error) {
	var session = "unset"
	f, err := c.File()
	if err == nil {
		session = fmt.Sprint(f.Fd())
	}
	log.Println("Accepted Connection", session)

	sc := bufio.NewScanner(c)
	for sc.Scan() {
		txt := sc.Text()
		parts := strings.Split(txt, ":")
		if len(parts) < 2 {
			continue
		}

		ir, iw, err := os.Pipe()
		if err != nil {
			errc <- errors.Wrap(err, "failed to make pipe 1 for echo")
			return
		}
		defer iw.Close()
		or, ow, err := os.Pipe()
		if err != nil {
			errc <- errors.Wrap(err, "failed to make pipe 2 for echo")
			return
		}
		defer ow.Close()

		switch parts[1] {
		case "ping":
			go startPing(session, ow, ir)

		case "echo":
			go startEcho(session, ow, ir)

		default:
			fmt.Fprintln(c, session,": unknown command:", parts)
			break
		}
		fd.Put(c, or, iw) // we need to put single file descriptors across this (cant squash them into a single pipe)
		log.Println(txt)
		// handle response
	}
	if err := sc.Err(); err != nil {
		log.Println("scanerr:", err)
	}
	err=c.Close()
	log.Println(session,"connection closed",err)
	return
}

func startPing(s string, w io.WriteCloser, r io.Reader) {
	go func() {
		sc := bufio.NewScanner(r)
		for sc.Scan() {
			txt := sc.Text()
			log.Println(s,"from ping:", txt)
		}
		log.Println(s,"ping read closed")
	}()
	for i := 5; i > 0; i-- {
		time.Sleep(1 * time.Second)
		fmt.Fprintln(w, "ping:", time.Now().Unix())
	}
	w.Close()
	log.Println(s,"ping write closed")
}

func startEcho(s string, w io.WriteCloser, r io.Reader) {
	sc := bufio.NewScanner(r)
	for sc.Scan() {
		txt := sc.Text()
		log.Println(s,"from echo:", txt)
		// handle response
		fmt.Fprintln(w, txt)
	}
	w.Close()
	log.Println(s,"echo closed")
}

package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"net"
	"os"
	"time"

	"github.com/cryptix/goremutake"
	"github.com/ftrvxmtrx/fd"
)

const path = "/tmp/sendMsgTest.sock"

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

				log.Println("Accepting Connection")
				go handleConn(clientConn, errc)
			}

		}
		close(done)

	}()
	<-done

	log.Println("Accept loop closed.")

}

func check(err error) {
	if err != nil {
		panic(err)
	}
}

func handleConn(c *net.UnixConn, errc chan error) {
	defer c.Close()

	fname := "/tmp/" + goremutake.Encode(uint(1024+rand.Intn(10000)))

	f, err := os.Create(fname)
	if err != nil {
		errc <- fmt.Errorf("Create rand file failed: %s", err)
		return
	}
	defer f.Close()

	err = fd.Put(c, f)
	if err != nil {
		errc <- fmt.Errorf("Create rand file failed: %s", err)
		return
	}

	log.Println("FD Send")
	time.Sleep(10 * time.Second)
	log.Println("timout.")
	go cleanupFile(fname)

	return
}

func cleanupFile(fname string) {
	time.Sleep(1 * time.Second)

	f, err := ioutil.ReadFile(fname)
	check(err)

	log.Printf("Content of File:%s\n", fname)
	log.Print(string(f))

	time.Sleep(1 * time.Second)

	log.Println("Removing file")
	os.Remove(fname)
}

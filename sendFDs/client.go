package main

import (
	"errors"
	"math/rand"
	"net"
	"time"

	"github.com/cryptix/goremutake"
	"github.com/ftrvxmtrx/fd"
)

const path = "/tmp/sendMsgTest.sock"

func main() {
	rand.Seed(time.Now().Unix())

	laddr, err := net.ResolveUnixAddr("unix", "")
	check(err)

	raddr, err := net.ResolveUnixAddr("unix", path)
	check(err)

	conn, err := net.DialUnix("unix", laddr, raddr)
	check(err)

	f, err := fd.Get(conn, 1, []string{"duh"})
	check(err)

	if len(f) < 1 {
		check(errors.New("not enough fds!"))
	}

	_, err = f[0].WriteString("Hello world!" + goremutake.Encode(uint(1024+rand.Intn(10000))))
	check(err)

	check(f[0].Close())

}

func check(err error) {
	if err != nil {
		panic(err)
	}
}

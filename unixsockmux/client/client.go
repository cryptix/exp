package client

import (
	"bufio"
	"fmt"
	"github.com/ftrvxmtrx/fd"
	"github.com/pkg/errors"
	"io"
	"log"
	"net"
	"sync/atomic"
)

type Client struct {
	reqs uint32
	conn *net.UnixConn
}

func NewClient(path string) (*Client, error) {
	laddr, err := net.ResolveUnixAddr("unix", "")
	if err != nil {
		return nil, err
	}

	raddr, err := net.ResolveUnixAddr("unix", path)
	if err != nil {
		return nil, err
	}

	conn, err := net.DialUnix("unix", laddr, raddr)
	if err != nil {
		return nil, err
	}

	c := &Client{conn: conn}
	go c.readIncomming()
	return c, nil
}

func (c *Client) Close() error {
	// close open requests
	return c.conn.Close()
}

func (c *Client) readIncomming() {

	sc := bufio.NewScanner(c.conn)
	for sc.Scan() {
		txt := sc.Text()
		log.Println(txt)
		// handle response
	}
}

func (c *Client) OpenChannel(name string, args interface{}) (io.ReadWriteCloser, error) {
	fmt.Fprintf(c.conn, "%d:%s:%+v\n", c.reqs, name, args)
	atomic.AddUint32(&c.reqs, 1)

	files, err := fd.Get(c.conn, 2, []string{name + ":r", name + ":w"})
	if err != nil {
		return nil, err
	}

	if n := len(files); n != 2 {
		return nil, errors.Errorf("not enough files?! %d", n)
	}

	return rwc{files[0], files[1]}, nil
}

type rwc struct {
	io.Reader
	io.WriteCloser
}

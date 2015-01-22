/*
Package manageSieve implements a simple client for the protocol for remotely managing sieve scripts.

	RFC5804 - https://tools.ietf.org/html/rfc5804

Done

These commands are implemented:
	STARTSSL
	AUTHENTICATE (Plain)
	LISTSCRIPTS
	GETSCRIPT

TODO

These need to be done:
	CAPABILITY
	HAVESPACE
	PUTSCRIPT
	SETACTIVE
	DELETESCRIPT
	RENAMESCRIPT
	CHECKSCRIPT
*/
package manageSieve

import (
	"bufio"
	"crypto/tls"
	"encoding/base64"
	"fmt"
	"net"
	"strings"
	"sync"
)

type Client struct {
	sl         sync.Mutex
	scanner    *bufio.Scanner
	conn       net.Conn
	serverName string
}

func Dial(addr string) (*Client, error) {
	conn, err := net.Dial("tcp", addr)
	if err != nil {
		return nil, err
	}
	// conn = debug.WrapConn(conn)
	host, _, _ := net.SplitHostPort(addr)
	return NewClient(conn, host)
}

func NewClient(conn net.Conn, host string) (*Client, error) {
	c := &Client{
		scanner:    bufio.NewScanner(conn),
		conn:       conn,
		serverName: host,
	}
	_, err := c.waitForOK() // TODO: parse capabilities
	if err != nil {
		return nil, err
	}
	return c, nil
}

func (c *Client) StartTLS(config *tls.Config) error {
	_, err := c.cmd(true, "STARTTLS")
	if err != nil {
		return err
	}
	if config == nil {
		config = &tls.Config{
			ServerName: c.serverName,
		}
	}
	c.conn = tls.Client(c.conn, config)
	// c.conn = debug.WrapConn(c.conn)
	c.scanner = bufio.NewScanner(c.conn)
	_, err = c.waitForOK() // TODO: parse capabilities
	return err
}

func (c *Client) Login(user, pass string) error {
	data := []byte("\x00" + user + "\x00" + pass)
	creds := base64.StdEncoding.EncodeToString(data)
	_, err := c.cmd(true, "AUTHENTICATE %q %q", "PLAIN", creds)
	return err
}

type Script struct {
	Name   string
	Active bool
}

func (c *Client) ListScripts() ([]Script, error) {
	b, err := c.cmd(true, "LISTSCRIPTS")
	if err != nil {
		return nil, err
	}
	var scripts []Script
	for _, l := range b[:len(b)-1] {
		var s Script
		if l[0] != '"' {
			return nil, fmt.Errorf("Illegal Line: %q", l)
		}
		s.Name = l[1:strings.LastIndex(l, "\"")]
		s.Active = strings.HasSuffix(l, "ACTIVE")
		scripts = append(scripts, s)
	}
	return scripts, nil
}

func (c *Client) GetScript(name string) (string, error) {
	b, err := c.cmd(true, "GETSCRIPT %q", name)
	if err != nil {
		return "", err
	}
	// TODO: Validate return?
	// scriptLenStr := b[0] // {int}
	// scriptLen, err := strconv.Atoi(scriptLenStr[1 : len(scriptLenStr)-1])
	// if err != nil {
	// 	return "", err
	// }
	// log.Println("scriptLen:", scriptLen)
	return strings.Join(b[1:len(b)-1], "\n"), nil
}

func (c *Client) cmd(wait bool, format string, args ...interface{}) ([]string, error) {
	c.sl.Lock()
	defer c.sl.Unlock()
	_, err := fmt.Fprintf(c.conn, format+"\r\n", args...)
	if err != nil {
		return nil, err
	}
	if wait {
		return c.waitForOK()
	}
	// TODO: not actually used..
	return []string{c.scanner.Text()}, c.scanner.Err()
	// return nil, nil
}

// read lines until OK
func (c *Client) waitForOK() ([]string, error) {
	var b []string
	for c.scanner.Scan() {
		l := c.scanner.Text()
		b = append(b, l)
		if strings.HasPrefix(l, "OK ") {
			break
		}
	}
	return b, c.scanner.Err()
}

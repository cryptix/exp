/*
Package wslog implements a (github.com/op/go-logging).Backend and a (net/http).Handler

The Handler upgrades requests to websocket connections and registeres them with the backend.

The Backend copies log records to all registerd clients.
*/
package wslog

import (
	"fmt"
	"net/http"
	"os"

	"github.com/gorilla/websocket"
	"github.com/op/go-logging"
)

// formattedRec extends a record with it's formatted Message
type formattedRec struct {
	logging.Record
	Msg string
}

// WebsocketBackend is a Hub-like Backend that copies log messages to listening connections
type WebsocketBackend struct {
	connections map[*websocketConn]bool // registered connections
	broadcast   chan *formattedRec      // sends on this
	register    chan *websocketConn     // register requests from the connections.
	unregister  chan *websocketConn     // unregister requests from connections.

	upgrader *websocket.Upgrader // elevates a normal http request to a websocket connection
}

// NewBackend creates a new WebsocketBackend.
func NewBackend() *WebsocketBackend {
	wb := &WebsocketBackend{
		broadcast:   make(chan *formattedRec),
		register:    make(chan *websocketConn),
		unregister:  make(chan *websocketConn),
		connections: make(map[*websocketConn]bool),

		// TODO(cryptix): configure sizes?
		upgrader: &websocket.Upgrader{ReadBufferSize: 1024, WriteBufferSize: 1024},
	}

	go wb.run()

	return wb
}

// run does the event handling
func (b *WebsocketBackend) run() {
	for {
		select {

		// add new connections to the pool
		case c := <-b.register:
			b.connections[c] = true

		// remove connections from the pool
		case c := <-b.unregister:
			if _, ok := b.connections[c]; ok {
				delete(b.connections, c)
				close(c.send)
			}

		// copy records to every connected client
		case rec := <-b.broadcast:
			for c := range b.connections {
				select {
				case c.send <- rec:
				default: // remove connection if send fails
					delete(b.connections, c)
					close(c.send)
				}
			}
		}
	}
}

// Log implements the logging.Backend interface. It broadcasts Records to the registerd connections
func (b *WebsocketBackend) Log(level logging.Level, calldepth int, rec *logging.Record) error {
	// format here to not loose original calling information
	b.broadcast <- &formattedRec{
		Record: *rec,
		Msg:    rec.Formatted(calldepth + 1),
	}
	return nil
}

// ServeHTTP implements net/http.Handler.
// upgrades the request to a websocket connection
// if successfull, registers it with the backend
func (b *WebsocketBackend) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ws, err := b.upgrader.Upgrade(w, r, nil)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Upgrade Error: %v\n", err)
		return
	}

	// TODO(cryptix): configure sizes?
	c := &websocketConn{
		send: make(chan *formattedRec, 10),
		ws:   ws,
	}

	// (un)register onto the hub
	b.register <- c
	defer func() { b.unregister <- c }()

	// do the i/o
	go c.writer()
	c.reader()
}

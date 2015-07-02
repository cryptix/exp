package wslog

import (
	"fmt"
	"os"

	"github.com/gorilla/websocket"
)

// TODO(cryptix): could be extended to be a leveld so that clients control which events they get
type websocketConn struct {
	ws   *websocket.Conn    // the websocket connection
	send chan *formattedRec // buffered channel of outbound messages.
}

// get messages from clients
func (c *websocketConn) reader() {
	for {
		_, message, err := c.ws.ReadMessage()
		if err != nil {
			fmt.Fprintf(os.Stderr, "ws.ReadMessage error: %v\n", err)
			break
		}
		fmt.Fprintf(os.Stderr, "From Client: %q\n", string(message))
	}
	c.ws.Close()
}

// push down record to listening clients
func (c *websocketConn) writer() {

	for rec := range c.send {
		err := c.ws.WriteJSON(rec)
		if err != nil {
			fmt.Fprintf(os.Stderr, "ws.WriteJSON error: %v\n", err)
			break
		}
	}
	c.ws.Close()
}

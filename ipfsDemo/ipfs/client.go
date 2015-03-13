package ipfs

import (
	"encoding/json"

	"github.com/jbenet/go-ipfs/commands"
	cmdsHttp "github.com/jbenet/go-ipfs/commands/http"
	coreCmds "github.com/jbenet/go-ipfs/core/commands"
	"gopkg.in/errgo.v1"
)

type Client struct {
	c cmdsHttp.Client
}

var ErrTODO = errgo.New("TODO")

func NewClient(addr string) (*Client, error) {
	// prepare connection to client
	return &Client{
		c: cmdsHttp.NewClient(addr),
	}, nil
}

func (c *Client) Ls(hashes ...string) (*coreCmds.LsOutput, error) {
	// prepareing everything we need to get a req from NewRequest
	opts := map[string]interface{}{commands.EncShort: "json"}

	// get the default options for the cmd we want to run
	optDefs, err := coreCmds.LsCmd.GetOptions([]string{})
	if err != nil {
		return nil, errgo.Notef(err, "GetOptions failed for LsCmd")
	}

	req, err := commands.NewRequest(
		[]string{"ls"}, // path - tells the daemon
		opts,           // options for the command
		hashes,         // arguments for the command
		nil,            // file, not used here
		coreCmds.LsCmd, // which command - so the client knows
		optDefs)        // default options
	if err != nil {
		return nil, errgo.Notef(err, "NewRequest failed")
	}

	// send it
	res, err := c.c.Send(req)
	if err != nil {
		return nil, errgo.Notef(err, "client.Send failed")
	}

	// read it
	r, err := res.Reader()
	if err != nil {
		return nil, errgo.Notef(err, "getting Reader failed")
	}

	o := new(coreCmds.LsOutput)

	if err = json.NewDecoder(r).Decode(o); err != nil {
		return nil, errgo.Notef(err, "json decode failed")
	}

	return o, nil
}

// small example of interfacing with the http api of the daemon
package main

import (
	"os"

	"github.com/cryptix/exp/ipfsDemo/ipfs"
	"github.com/shurcooL/go-goon"
)

func main() {

	c, err := ipfs.NewClient("127.0.0.1:5001")
	check(err)

	o, err := c.Ls(os.Args[0])
	check(err)

	goon.Dump(o)

}

func check(err error) {
	if err != nil {
		panic(err)
	}
}

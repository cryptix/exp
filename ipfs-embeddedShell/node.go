package embeddedShell

import (
	"github.com/ipfs/go-ipfs/core"
	"github.com/ipfs/go-ipfs/repo/fsrepo"
	"golang.org/x/net/context"
	"gopkg.in/errgo.v1"
)

func NewDefaultNodeWithFSRepo(ctx context.Context, repoPath string) (*core.IpfsNode, error) {
	r, err := fsrepo.Open(repoPath)
	if err != nil {
		return nil, errgo.Notef(err, "opening fsrepo failed")
	}
	node, err := core.NewNode(ctx, &core.BuildCfg{
		Online: true,
		Repo:   r,
	})
	if err != nil {
		return nil, errgo.Notef(err, "ipfs NewNode() failed.")
	}
	// TODO: can we bootsrap localy/mdns first and fall back to default?
	err = node.Bootstrap(core.DefaultBootstrapConfig)
	if err != nil {
		return nil, errgo.Notef(err, "ipfs Bootstrap() failed.")
	}
	return node, nil
}

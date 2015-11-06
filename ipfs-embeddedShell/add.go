package embeddedShell

import (
	"io"

	"gopkg.in/errgo.v1"

	"github.com/ipfs/go-ipfs/importer"
	"github.com/ipfs/go-ipfs/importer/chunk"
)

func (s *Shell) Add(r io.Reader) (string, error) {
	dag, err := importer.BuildDagFromReader(
		s.node.DAG,
		chunk.DefaultSplitter(r),
		importer.BasicPinnerCB(s.node.Pinning.GetManual()), // TODO: make pinning configurable
	)
	if err != nil {
		return "", errgo.Notef(err, "add: importing DAG failed.")
	}
	k, err := dag.Key()
	if err != nil {
		return "", errgo.Notef(err, "add: getting key from DAG failed.")
	}
	return k.B58String(), nil
}

package git

import (
	"bufio"
	"bytes"
	"compress/zlib"
	"fmt"
	"io"
	"io/ioutil"
	"strconv"

	"gopkg.in/errgo.v1"
)

type Object struct {
	Type    string // enum?
	Content []byte
}

type Tree struct {
	Mode string
	Name string

	SHA1Sum []byte
}

func (o Object) String() string {
	return fmt.Sprintf("%5s: %q", o.Type, o.Content)
}

func DecodeObject(r io.Reader) (*Object, error) {

	zr, err := zlib.NewReader(r)
	if err != nil {
		return nil, errgo.Notef(err, "zlib newReader failed")
	}

	br := bufio.NewReader(zr)

	header, err := br.ReadBytes(0)
	if err != nil {
		return nil, errgo.Notef(err, "error finding header 0byte")
	}

	o := &Object{}
	var hdrLenStr string
	switch {
	case bytes.HasPrefix(header, []byte("blob ")):
		o.Type = "blob"
		hdrLenStr = string(header[5 : len(header)-1])

	case bytes.HasPrefix(header, []byte("tree ")):
		o.Type = "tree"
		hdrLenStr = string(header[5 : len(header)-1])

	case bytes.HasPrefix(header, []byte("commit ")):
		o.Type = "tree"
		hdrLenStr = string(header[7 : len(header)-1])

	default:
		return nil, errgo.Newf("illegal git object:%q", header)
	}

	hdrLen, err := strconv.ParseInt(hdrLenStr, 10, 64)
	if err != nil {
		return nil, errgo.Notef(err, "error parsing header length")
	}
	lr := io.LimitReader(br, hdrLen)
	o.Content, err = ioutil.ReadAll(lr)
	if err != nil {
		return nil, errgo.Notef(err, "error finding header 0byte")
	}

	return o, nil
}

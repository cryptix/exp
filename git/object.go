package git

import (
	"bufio"
	"bytes"
	"compress/zlib"
	"crypto/sha1"
	"encoding/hex"
	"fmt"
	"io"
	"io/ioutil"
	"strconv"

	"gopkg.in/errgo.v1"
)

type Type int

const (
	_ Type = iota
	BlobT
	TreeT
	CommitT
)

type Object struct {
	Type Type
	Size int64
	tree []Tree
	blob *Blob
}

func newBlob(content []byte) *Blob {
	b := Blob(content)
	return &b
}

type Blob []byte

type Tree struct {
	Mode, Name string
	SHA1Sum    [sha1.Size]byte
}

func (o Object) String() string {
	switch o.Type {
	case BlobT:
		if o.blob == nil {
			return "broken blob"
		}
		return fmt.Sprintf("blob<%d> %q", o.Size, string(*o.blob))

	case TreeT:
		if o.tree == nil {
			return "broken blob"
		}
		s := fmt.Sprintf("tree<%d>\n", o.Size)
		for _, t := range o.tree {
			s += fmt.Sprintf("%q\t%q\t%s\n", t.Mode, t.Name, hex.EncodeToString(t.SHA1Sum[:]))
		}
		return s

	default:
		return "broken object"
	}
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
		o.Type = BlobT
		hdrLenStr = string(header[5 : len(header)-1])

	case bytes.HasPrefix(header, []byte("tree ")):
		o.Type = TreeT
		hdrLenStr = string(header[5 : len(header)-1])

	case bytes.HasPrefix(header, []byte("commit ")):
		o.Type = CommitT
		hdrLenStr = string(header[7 : len(header)-1])

	default:
		return nil, errgo.Newf("illegal git object:%q", header)
	}

	hdrLen, err := strconv.ParseInt(hdrLenStr, 10, 64)
	if err != nil {
		return nil, errgo.Notef(err, "error parsing header length")
	}
	o.Size = hdrLen
	lr := io.LimitReader(br, hdrLen)

	switch o.Type {
	case BlobT:
		content, err := ioutil.ReadAll(lr)
		if err != nil {
			return nil, errgo.Notef(err, "error finding header 0byte")
		}
		o.blob = newBlob(content)
		return o, nil

	case TreeT:
		o.tree, err = decodeTreeEntries(lr)
		if err != nil {
			if errgo.Cause(err) == io.EOF {
				return o, nil
			}
			return nil, errgo.Notef(err, "decodecodeTreeEntries failed")
		}
		return o, nil

	default:
		return nil, errgo.Newf("illegal object type:%T %v", o.Type, o.Type)
	}
}

func decodeTreeEntries(r io.Reader) ([]Tree, error) {
	isEOF := errgo.Is(io.EOF)
	var entries []Tree
	br := bufio.NewReader(r)
	for {
		var t Tree
		hdr, err := br.ReadSlice(0)
		if err != nil {
			return entries, errgo.NoteMask(err, "error finding modeName 0byte", isEOF)
		}
		modeName := bytes.Split(hdr[:len(hdr)-1], []byte(" "))
		if len(modeName) != 2 {
			return entries, errgo.Newf("illegal modeName block: %v", modeName)
		}
		t.Mode = string(modeName[0])
		t.Name = string(modeName[1])

		var hash [sha1.Size]byte
		n, err := br.Read(hash[:])
		if err != nil {
			return entries, errgo.NoteMask(err, "br.Read() hash failed", isEOF)
		}
		if n != 20 {
			return entries, errgo.Newf("br.Read() fell short: %d", n)
		}
		t.SHA1Sum = hash
		entries = append(entries, t)
	}
}

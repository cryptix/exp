// ipfs-waveform-kit produces waveform plots for audiofiles (using bbcrd/audiowaveform) and stores them in ipfs (ipfs.io)
package main

import (
	"encoding/json"
	"flag"
	"io"
	"log"
	"net/http"
	"os"
	"os/exec"

	"github.com/cryptix/exp/ipfs-embeddedShell"
	"golang.org/x/net/context"
	"gopkg.in/errgo.v1"

	"github.com/go-kit/kit/endpoint"
	httptransport "github.com/go-kit/kit/transport/http"
)

var (
	ErrBadRequest = errgo.New("ipfs-waveform-kit: BadRequest")
)

type WaveformService interface {
	Analyze(string) (string, error)
}

type waveformService struct {
	shell *embeddedShell.Shell
}

func (ws waveformService) Analyze(path string) (string, error) {
	rc, err := ws.shell.Cat(path)
	if err != nil {
		return "", errgo.Notef(err, "Analyze: cat(%q) failed.", path)
	}

	// throw into bbcrd audiowaveform
	// wait for stdio. https://github.com/bbcrd/audiowaveform/issues/13

	cmd := exec.Command("audiowaveform", "--input", "-")
	cmd.Stdin = rc
	outP, err := cmd.StdoutPipe()
	if err != nil {
		return "", errgo.Notef(err, "Analyze: failed to get output pipe fr audiowaveform tool")
	}

	logP, err := cmd.StderrPipe()
	if err == nil {
		go func() {
			_, err := io.Copy(os.Stderr, logP)
			log.Println("audiowaveform stderr:", err)
		}()
	}

	errc := make(chan error)

	go func() {
		errc <- cmd.Run()
	}()

	hash, err := ws.shell.Add(outP)
	if err != nil {
		return "", errgo.Notef(err, "Analyze: failed to ipfs.Add() waveform data")
	}

	if err := <-errc; err != nil {
		return "", errgo.Notef(err, "Analyze: subprocess failed")
	}

	if err := rc.Close(); err != nil {
		return "", errgo.Notef(err, "Analyze: failed to close cat.")
	}

	return hash, nil
}

var _ WaveformService = &waveformService{}

var repoPath = flag.String("ipfsrepo", "./ipfsRepo", "where to open the repo (use IPFS_PATH=... ipfs init)")

func main() {
	flag.Parse()

	_, err := exec.LookPath("audiowaveform")
	if err != nil {
		log.Fatal(err)
	}

	ctx := context.Background()

	node, err := embeddedShell.NewDefaultNodeWithFSRepo(ctx, *repoPath)
	if err != nil {
		log.Fatal(err)
	}

	shell := embeddedShell.NewShellWithContext(node, ctx)

	svc := waveformService{shell}

	analyzeHandler := httptransport.NewServer(
		ctx,
		makeAnalyzeEndpoint(svc),
		decodeAnalyzeRequest,
		encodeResponse,
	)

	http.Handle("/analyze", analyzeHandler)
	log.Fatal(http.ListenAndServe(":9080", nil))
}

func makeAnalyzeEndpoint(svc WaveformService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req, ok := request.(analyzeRequest)
		if !ok {
			return analyzeResponse{"", ErrBadRequest.Error()}, nil
		}
		v, err := svc.Analyze(req.S)
		if err != nil {
			return analyzeResponse{"", err.Error()}, nil
		}
		return analyzeResponse{v, ""}, nil
	}
}

func decodeAnalyzeRequest(r *http.Request) (interface{}, error) {
	var request analyzeRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		return nil, err
	}
	return request, nil
}

func encodeResponse(w http.ResponseWriter, resp interface{}) error {
	return json.NewEncoder(w).Encode(resp)
}

type analyzeRequest struct {
	S string `json:"s"`
}

type analyzeResponse struct {
	V   string `json:"v"`
	Err string `json:"err,omitempty"` // errors don't define JSON marshaling
}

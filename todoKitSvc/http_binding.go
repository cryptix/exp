package main

import (
	"encoding/json"
	"net/http"

	"github.com/go-kit/kit/endpoint"
	httptransport "github.com/go-kit/kit/transport/http"
	"golang.org/x/net/context"
)

func makeHTTPBinding(ctx context.Context, e endpoint.Endpoint, reqDecoder httptransport.DecodeFunc, before []httptransport.BeforeFunc, after []httptransport.AfterFunc) http.Handler {
	encode := func(w http.ResponseWriter, response interface{}) error {
		return json.NewEncoder(w).Encode(response)
	}
	return httptransport.Server{
		Context:    ctx,
		Endpoint:   e,
		DecodeFunc: reqDecoder,
		EncodeFunc: encode,
		Before:     before,
		After:      append([]httptransport.AfterFunc{httptransport.SetContentType("application/json; charset=utf-8")}, after...),
	}
}

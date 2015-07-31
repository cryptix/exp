package main

import (
	"github.com/go-kit/kit/endpoint"
	"golang.org/x/net/context"

	"github.com/cryptix/exp/todoKitSvc/reqrep"
)

type NetrpcBinding struct {
	ctx       context.Context
	endpoints todoEndpoints
}

func (b NetrpcBinding) Add(request reqrep.AddRequest, response *reqrep.AddResponse) error {
	var (
		ctx, cancel = context.WithCancel(b.ctx)
		errs        = make(chan error, 1)
		responses   = make(chan reqrep.AddResponse, 1)
	)
	defer cancel()
	go func() {
		resp, err := b.endpoints.Add(ctx, request)
		if err != nil {
			errs <- err
			return
		}
		addResp, ok := resp.(reqrep.AddResponse)
		if !ok {
			errs <- endpoint.ErrBadCast
			return
		}
		responses <- addResp
	}()
	select {
	case <-ctx.Done():
		return context.DeadlineExceeded
	case err := <-errs:
		return err
	case resp := <-responses:
		(*response) = resp
		return nil
	}

}

func (b NetrpcBinding) List(request reqrep.ListRequest, response *reqrep.ListResponse) error {
	var (
		ctx, cancel = context.WithCancel(b.ctx)
		errs        = make(chan error, 1)
		responses   = make(chan reqrep.ListResponse, 1)
	)
	defer cancel()
	go func() {
		resp, err := b.endpoints.List(ctx, request)
		if err != nil {
			errs <- err
			return
		}
		listResp, ok := resp.(reqrep.ListResponse)
		if !ok {
			errs <- endpoint.ErrBadCast
			return
		}
		responses <- listResp
	}()
	select {
	case <-ctx.Done():
		return context.DeadlineExceeded
	case err := <-errs:
		return err
	case resp := <-responses:
		(*response) = resp
		return nil
	}
}

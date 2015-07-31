package main

import (
	"github.com/go-kit/kit/endpoint"
	"golang.org/x/net/context"

	"github.com/cryptix/exp/todoKitSvc/reqrep"
	"github.com/cryptix/exp/todoKitSvc/todosvc"
)

// TODO add rest of methods
func makeTodoServerEndpoints(t todosvc.Todo) todosvc.Endpoints {
	return todosvc.Endpoints{
		Add:  makeAddEndpoint(t),
		List: makeListEndpoint(t),
	}
}

func makeAddEndpoint(t todosvc.Todo) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		select {
		default:
		case <-ctx.Done():
			return nil, endpoint.ErrContextCanceled
		}

		addReq, ok := request.(reqrep.AddRequest)
		if !ok {
			return nil, endpoint.ErrBadCast
		}

		id, err := t.Add(ctx, addReq.Title)
		return reqrep.AddResponse{ID: id, Err: err}, nil // TODO(cryptix): do we want to return the error here?..
	}
}

func makeListEndpoint(t todosvc.Todo) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		select {
		default:
		case <-ctx.Done():
			return nil, endpoint.ErrContextCanceled
		}

		_, ok := request.(reqrep.ListRequest)
		if !ok {
			return nil, endpoint.ErrBadCast
		}

		l, err := t.List(ctx)
		return reqrep.ListResponse{List: l, Err: err}, nil // TODO(cryptix): do we want to return the error here?..
	}
}

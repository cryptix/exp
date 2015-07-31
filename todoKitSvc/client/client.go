package client

import (
	"golang.org/x/net/context"

	"github.com/cryptix/exp/todoKitSvc/reqrep"
	"github.com/cryptix/exp/todoKitSvc/todosvc"
	"github.com/go-kit/kit/endpoint"
)

func NewClient(ep todosvc.Endpoints) todosvc.Todo {
	return endpointClient{ep}
}

type endpointClient struct{ ep todosvc.Endpoints }

func (c endpointClient) Add(ctx context.Context, title string) (todosvc.ID, error) {
	response, err := c.ep.Add(ctx, reqrep.AddRequest{Title: title})
	if err != nil {
		return -1, err
	}
	addResponse, ok := response.(reqrep.AddResponse)
	if !ok {
		return -1, endpoint.ErrBadCast
	}
	return addResponse.ID, addResponse.Err
}

func (c endpointClient) List(ctx context.Context) ([]todosvc.Item, error) {
	response, err := c.ep.List(ctx, reqrep.ListRequest{})
	if err != nil {
		return nil, err
	}
	listResponse, ok := response.(reqrep.ListResponse)
	if !ok {
		return nil, endpoint.ErrBadCast
	}
	return listResponse.List, listResponse.Err
}

func (c endpointClient) Toggle(ctx context.Context, id todosvc.ID) error {
	panic("not implemented")
}

func (c endpointClient) Delete(ctx context.Context, id todosvc.ID) error {
	panic("not implemented")
}

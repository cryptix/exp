package netrpc

import (
	"net/rpc"

	"golang.org/x/net/context"

	"github.com/cryptix/exp/todoKitSvc/reqrep"
	"github.com/cryptix/exp/todoKitSvc/todosvc"
)

func NewClient(c *rpc.Client) todosvc.Endpoints {
	return todosvc.Endpoints{
		Add: func(ctx context.Context, request interface{}) (interface{}, error) {
			var (
				errs      = make(chan error, 1)
				responses = make(chan interface{}, 1)
			)
			go func() {
				var response reqrep.AddResponse
				if err := c.Call("todosvc.Add", request, &response); err != nil {
					errs <- err
					return
				}
				responses <- response
			}()
			select {
			case <-ctx.Done():
				return nil, context.DeadlineExceeded
			case err := <-errs:
				return nil, err
			case response := <-responses:
				return response, nil
			}
		},
		List: func(ctx context.Context, request interface{}) (interface{}, error) {
			var (
				errs      = make(chan error, 1)
				responses = make(chan interface{}, 1)
			)
			go func() {
				var response reqrep.ListResponse
				if err := c.Call("todosvc.List", request, &response); err != nil {
					errs <- err
					return
				}
				responses <- response
			}()
			select {
			case <-ctx.Done():
				return nil, context.DeadlineExceeded
			case err := <-errs:
				return nil, err
			case response := <-responses:
				return response, nil
			}
		},
	}
}

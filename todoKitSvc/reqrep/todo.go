package reqrep

import "github.com/cryptix/exp/todoKitSvc/todosvc"

// Add

type AddRequest struct {
	Name string
}

type AddResponse struct {
	ID todosvc.ID
}

// List

type ListRequest struct {
	// TODO filtering options?
}

type ListResponse struct {
	List []todosvc.Item
	Err  error
}

package reqrep

import "github.com/cryptix/exp/todoKitSvc/todosvc"

// Add

type AddRequest struct {
	Title string
}

type AddResponse struct {
	ID  todosvc.ID
	Err error
}

// List

type ListRequest struct {
	// TODO filtering options?
}

type ListResponse struct {
	List []todosvc.Item
	Err  error
}

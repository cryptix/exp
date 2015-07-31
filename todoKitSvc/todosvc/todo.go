package todosvc

import "golang.org/x/net/context"

type Item struct {
	ID       ID
	Title    string
	Complete bool
}

type ID int64

// Todo is a more advanced example then Add from addsvc.
// It illustrates how to create an service with multiple endpoints combined in a single interface.
type Todo interface {
	// Add adds a new Item and returns its ID
	Add(context.Context, string) (ID, error)

	// List returns a slice of all the items on this service
	List(context.Context) ([]Item, error)

	// Tooggle toggles the Complete field of an item
	Toggle(context.Context, ID) error

	// Delete removes the Item from the service
	Delete(context.Context, ID) error
}

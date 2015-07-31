package todosvc

import (
	"errors"
	"sync"

	"golang.org/x/net/context"
)

func NewInmemTodo() Todo {
	idsrc := make(chan ID)

	go func() {
		var i ID = 23
		for {
			idsrc <- i
			i += 1
		}
	}()
	return inmem{
		idsrc: idsrc,
		store: make(map[ID]Item),
	}
}

var ErrNotFound = errors.New("todo: not found")

type inmem struct {
	idsrc <-chan ID

	mu    sync.RWMutex
	store map[ID]Item
}

func (t inmem) Add(_ context.Context, title string) (ID, error) {
	t.mu.Lock()
	defer t.mu.Unlock()
	id := <-t.idsrc
	t.store[id] = Item{
		ID:    id,
		Title: title,
	}
	return id, nil
}

func (t inmem) List(context.Context) ([]Item, error) {
	t.mu.RLock()
	defer t.mu.RUnlock()
	var items []Item
	for _, v := range t.store {
		items = append(items, v)
	}
	return items, nil
}

func (t inmem) Toggle(_ context.Context, id ID) error {
	t.mu.Lock()
	defer t.mu.Unlock()
	i, ok := t.store[id]
	if !ok {
		return ErrNotFound
	}
	i.Complete = !i.Complete
	t.store[id] = i
	return nil
}

func (t inmem) Delete(_ context.Context, id ID) error {
	t.mu.Lock()
	defer t.mu.Unlock()
	_, ok := t.store[id]
	if !ok {
		return ErrNotFound
	}
	delete(t.store, id)
	return nil
}

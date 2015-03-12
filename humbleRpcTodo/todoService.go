package main

import (
	"log"
	"sort"
	"sync"

	"github.com/cryptix/exp/humbleRpcTodo/types"
	"gopkg.in/errgo.v1"
)

var (
	// todosMutex protects access to the todos map
	todosMutex = sync.Mutex{}
	// todos stores all the todos as a map of id to *Todo
	todos = make(types.TodosIndex, 10)
	// todosCounter is incremented every time a new todo is created
	// it is used to set todo ids.
	todosCounter = 1
)

func init() {
	createTodo("Write a frontend framework in Go")
	createTodo("???")
	createTodo("Profit!")
}

func createTodo(title string) *types.Todo {
	todosMutex.Lock()
	defer todosMutex.Unlock()
	id := todosCounter
	todosCounter++
	todo := &types.Todo{
		Id:    id,
		Title: title,
	}
	todos[id] = todo
	return todo
}

var ErrTodoNotFound = errgo.New("todo not found")

// Todos RPC
type TodoService struct{}

func (*TodoService) List(args *types.TodoListArgs, reply *[]types.Todo) error {
	todoList := make([]types.Todo, len(todos))
	i := 0
	for _, v := range todos {
		todoList[i] = *v
		i++
	}
	sort.Sort(types.ById(todoList))
	*reply = todoList
	log.Printf("Listed %+v", reply)
	return nil
}

func (*TodoService) Save(in *types.Todo, out *types.Todo) error {
	if in == nil || in.Title == "" {
		return errgo.New("no title")
	}
	// create
	if in.Id == 0 {
		newTodo := createTodo(in.Title)
		log.Printf("Created %+v", newTodo)
		*out = *newTodo
		return nil
	}

	// save
	todosMutex.Lock()
	defer todosMutex.Unlock()

	if _, found := todos[in.Id]; !found {
		return ErrTodoNotFound
	}

	todos[in.Id] = in
	*out = *in
	log.Printf("Saved %+v", out)
	return nil
}

func (*TodoService) Delete(id *int, _ *struct{}) error {
	if id == nil {
		return errgo.New("invalid id")
	}
	todosMutex.Lock()
	defer todosMutex.Unlock()

	if _, ok := todos[*id]; !ok {
		return ErrTodoNotFound
	}

	// Delete the todo and render a response
	delete(todos, *id)
	log.Printf("Deleted %d", *id)
	return nil
}

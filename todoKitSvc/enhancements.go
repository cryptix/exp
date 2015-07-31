package main

import (
	"time"

	"github.com/cryptix/exp/todoKitSvc/todosvc"
	"github.com/go-kit/kit/log"
	"golang.org/x/net/context"
)

func NewLoggingTodo(l log.Logger, t todosvc.Todo) todosvc.Todo {
	return loggingTodo{l, t}
}

type loggingTodo struct {
	logger log.Logger
	next   todosvc.Todo
}

func (l loggingTodo) Add(ctx context.Context, title string) (id todosvc.ID, err error) {
	defer func(begin time.Time) {
		l.logger.Log("action", "add", "title", title, "ID", id, "err", err, "took", time.Since(begin))
	}(time.Now())
	id, err = l.next.Add(ctx, title)
	return id, err
}

func (l loggingTodo) List(ctx context.Context) (items []todosvc.Item, err error) {
	defer func(begin time.Time) {
		l.logger.Log("action", "list", "count", len(items), "err", err, "took", time.Since(begin))
	}(time.Now())
	items, err = l.next.List(ctx)
	return items, err
}

func (l loggingTodo) Toggle(ctx context.Context, id todosvc.ID) (err error) {
	defer func(begin time.Time) {
		l.logger.Log("action", "toggle", "ID", id, "err", err, "took", time.Since(begin))
	}(time.Now())
	err = l.next.Toggle(ctx, id)
	return err
}

func (l loggingTodo) Delete(ctx context.Context, id todosvc.ID) (err error) {
	defer func(begin time.Time) {
		l.logger.Log("action", "delete", "ID", id, "err", err, "took", time.Since(begin))
	}(time.Now())
	err = l.next.Delete(ctx, id)
	return err
}

// TODO
//func instrument(requests metrics.Counter, duration metrics.TimeHistogram) func(Todo) Todo {
//	return func(next Todo) Todo {
//		return func(ctx context.Context, a, b int64) int64 {
//			defer func(begin time.Time) {
//				requests.Add(1)
//				duration.Observe(time.Since(begin))
//			}(time.Now())
//			return next(ctx, a, b)
//		}
//	}
//}

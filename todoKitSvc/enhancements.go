package main

import (
	"time"

	"github.com/cryptix/exp/todoKitSvc/todosvc"
	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/metrics"
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

func NewInstrumentedTodo(r metrics.Counter, d metrics.TimeHistogram, t todosvc.Todo) todosvc.Todo {
	return instrumentedTodo{r, d, t}
}

type instrumentedTodo struct {
	reqcnt metrics.Counter
	reqdur metrics.TimeHistogram
	next   todosvc.Todo
}

func (l instrumentedTodo) Add(ctx context.Context, title string) (id todosvc.ID, err error) {
	defer func(begin time.Time) {
		l.reqcnt.Add(1)
		l.reqdur.Observe(time.Since(begin))
	}(time.Now())
	id, err = l.next.Add(ctx, title)
	return id, err
}

func (l instrumentedTodo) List(ctx context.Context) (items []todosvc.Item, err error) {
	defer func(begin time.Time) {
		l.reqcnt.Add(1)
		l.reqdur.Observe(time.Since(begin))
	}(time.Now())
	items, err = l.next.List(ctx)
	return items, err
}

func (l instrumentedTodo) Toggle(ctx context.Context, id todosvc.ID) (err error) {
	defer func(begin time.Time) {
		l.reqcnt.Add(1)
		l.reqdur.Observe(time.Since(begin))
	}(time.Now())
	err = l.next.Toggle(ctx, id)
	return err
}

func (l instrumentedTodo) Delete(ctx context.Context, id todosvc.ID) (err error) {
	defer func(begin time.Time) {
		l.reqcnt.Add(1)
		l.reqdur.Observe(time.Since(begin))
	}(time.Now())
	err = l.next.Delete(ctx, id)
	return err
}

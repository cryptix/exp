package main

import (
	"encoding/json"
	"flag"
	"fmt"
	stdlog "log"
	"math/rand"
	"net/http"
	"net/rpc"
	"os"
	"os/signal"
	"syscall"
	"time"

	kitlog "github.com/go-kit/kit/log"
	httptransport "github.com/go-kit/kit/transport/http"
	"golang.org/x/net/context"

	"github.com/cryptix/exp/todoKitSvc/reqrep"
	"github.com/cryptix/exp/todoKitSvc/todosvc"
)

func main() {
	// Flag domain. Note that gRPC transitively registers flags via its import
	// of glog. So, we define a new flag set, to keep those domains distinct.
	fs := flag.NewFlagSet("", flag.ExitOnError)
	var (
		debugAddr  = fs.String("debug.addr", ":8000", "Address for HTTP debug/instrumentation server")
		httpAddr   = fs.String("http.addr", ":8001", "Address for HTTP (JSON) server")
		netrpcAddr = fs.String("netrpc.addr", ":8003", "Address for net/rpc server")

		proxyHTTPURL = fs.String("proxy.http.url", "", "if set, proxy requests over HTTP to this todosvc")
	)
	flag.Usage = fs.Usage // only show our flags
	fs.Parse(os.Args[1:])

	// `package log` domain
	var logger kitlog.Logger
	logger = kitlog.NewLogfmtLogger(os.Stderr)
	logger = kitlog.NewContext(logger).With("ts", kitlog.DefaultTimestampUTC)
	stdlog.SetOutput(kitlog.NewStdlibAdapter(logger)) // redirect stdlib logging to us
	stdlog.SetFlags(0)                                // flags are handled in our logger

	// `package metrics` domain

	// TODO(cryptix): instrument
	//requests := metrics.NewMultiCounter(
	//	expvar.NewCounter("requests"),
	//	statsd.NewCounter(ioutil.Discard, "requests_total", time.Second),
	//	prometheus.NewCounter(stdprometheus.CounterOpts{
	//		Namespace: "todosvc",
	//		Subsystem: "add",
	//		Name:      "requests_total",
	//		Help:      "Total number of received requests.",
	//	}, []string{}),
	//)
	//duration := metrics.NewTimeHistogram(time.Nanosecond, metrics.NewMultiHistogram(
	//	expvar.NewHistogram("duration_nanoseconds_total", 0, 1e9, 3, 50, 95, 99),
	//	statsd.NewHistogram(ioutil.Discard, "duration_nanoseconds_total", time.Second),
	//	prometheus.NewSummary(stdprometheus.SummaryOpts{
	//		Namespace: "todosvc",
	//		Subsystem: "add",
	//		Name:      "duration_nanoseconds_total",
	//		Help:      "Total nanoseconds spend serving requests.",
	//	}, []string{}),
	//))

	// Our business and operational domain
	var t todosvc.Todo = todosvc.NewInmemTodo()
	if *proxyHTTPURL != "" {
		panic("TODO - not implemented")
		//var e endpoint.Endpoint
		//e = httpTodo.NewClient("GET", *proxyHTTPURL, nil)
		//t = proxyTodo{e, logger}
	}
	t = NewLoggingTodo(logger, t)
	//t = instrumentTodo(requests, duration, t)

	// Server domain
	todoEndpoints := makeTodoEndpoints(t)

	// Mechanical stuff
	rand.Seed(time.Now().UnixNano())
	root := context.Background()
	errc := make(chan error)

	go func() {
		errc <- interrupt()
	}()

	// Transport: HTTP (debug/instrumentation)
	go func() {
		logger.Log("addr", *debugAddr, "transport", "debug")
		errc <- http.ListenAndServe(*debugAddr, nil)
	}()

	// Transport: HTTP (JSON)
	go func() {
		ctx, cancel := context.WithCancel(root)
		defer cancel()
		before := []httptransport.BeforeFunc{}
		after := []httptransport.AfterFunc{}
		mux := http.NewServeMux()
		// Add
		addDecode := func(r *http.Request) (interface{}, error) {
			var request reqrep.AddRequest
			if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
				return nil, err
			}
			return request, r.Body.Close()
		}
		addHandler := makeHTTPBinding(ctx, todoEndpoints.Add, addDecode, before, after)
		mux.Handle("/add", addHandler)
		// List
		listDecode := func(r *http.Request) (interface{}, error) {
			var request reqrep.ListRequest
			if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
				return nil, err
			}
			return request, r.Body.Close()
		}
		listHandler := makeHTTPBinding(ctx, todoEndpoints.List, listDecode, before, after)
		mux.Handle("/list", listHandler)
		logger.Log("addr", *httpAddr, "transport", "HTTP/JSON")
		errc <- http.ListenAndServe(*httpAddr, mux)
	}()

	// Transport: net/rpc
	go func() {
		ctx, cancel := context.WithCancel(root)
		defer cancel()
		s := rpc.NewServer()
		s.RegisterName("addsvc", NetrpcBinding{ctx, todoEndpoints})
		s.HandleHTTP(rpc.DefaultRPCPath, rpc.DefaultDebugPath)
		logger.Log("addr", *netrpcAddr, "transport", "net/rpc")
		errc <- http.ListenAndServe(*netrpcAddr, s)
	}()

	logger.Log("fatal", <-errc)
}

func interrupt() error {
	c := make(chan os.Signal)
	signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)
	return fmt.Errorf("%s", <-c)
}

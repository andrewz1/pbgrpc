package main

import (
	"context"
	"log"
	"net"
	"runtime"
	"sync"

	"google.golang.org/grpc"
	"google.golang.org/grpc/peer"

	"github.com/andrewz1/pbgrpc/mygrpc"
)

type task struct {
	req *mygrpc.Request
	rsp *mygrpc.Response
	sync.WaitGroup
}

type server struct {
	ch chan *task
}

func main() {
	listener, err := net.Listen("tcp", ":5300")
	if err != nil {
		log.Fatalf("listen: %v\n", err)
	}
	srv := &server{ch: make(chan *task, 10000)}
	for i := 0; i < runtime.NumCPU(); i++ {
		go srv.worker()
	}
	grpcServer := grpc.NewServer()
	mygrpc.RegisterReverseServer(grpcServer, srv)
	log.Fatalf("serve: %v\n", grpcServer.Serve(listener))
}

func (s *server) Do(ctx context.Context, req *mygrpc.Request) (rsp *mygrpc.Response, err error) {
	if p, ok := peer.FromContext(ctx); ok {
		log.Printf("from: %v\n", p.Addr)
	}
	if req != nil {
		log.Printf("req: %v\n", req)
	}
	t := &task{req: req}
	t.Add(1)
	s.ch <- t
	t.Wait()
	rsp = t.rsp
	return
}

func (s *server) worker() {
	for t := range s.ch {
		t.rsp = &mygrpc.Response{
			Message: "200 OK",
		}
		t.Done()
	}
}

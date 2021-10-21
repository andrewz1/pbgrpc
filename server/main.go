package main

import (
	"context"
	"log"
	"net"

	"google.golang.org/grpc"
	"google.golang.org/grpc/peer"

	"github.com/andrewz1/pbgrpc/mygrpc"
)

type server struct{}

func main() {
	listener, err := net.Listen("tcp", ":5300")
	if err != nil {
		log.Fatalf("listen: %v\n", err)
	}
	grpcServer := grpc.NewServer()
	mygrpc.RegisterReverseServer(grpcServer, server{})
	log.Fatalf("serve: %v\n", grpcServer.Serve(listener))
}

func (s server) Do(ctx context.Context, req *mygrpc.Request) (rsp *mygrpc.Response, err error) {
	if p, ok := peer.FromContext(ctx); ok {
		log.Printf("from: %v\n", p.Addr)
	}
	if req != nil {
		log.Printf("req: %v\n", req)
	}
	rsp = &mygrpc.Response{
		Message: "200 OK",
	}
	return
}

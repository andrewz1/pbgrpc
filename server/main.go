package main

import (
	"context"
	"net"

	"google.golang.org/grpc"
	"google.golang.org/grpc/grpclog"
	"google.golang.org/grpc/peer"

	"github.com/andrewz1/pbgrpc/mygrpc"
)

type server struct{}

func main() {
	listener, err := net.Listen("tcp", ":5300")
	if err != nil {
		grpclog.Fatalf("listen: %v", err)
	}
	grpcServer := grpc.NewServer()

	mygrpc.RegisterReverseServer(grpcServer, server{})

	grpclog.Fatalf("serve: %v", grpcServer.Serve(listener))
}

func (s server) Do(ctx context.Context, req *mygrpc.Request) (rsp *mygrpc.Response, err error) {
	if p, ok := peer.FromContext(ctx); ok {
		grpclog.Infof("from: %v", p.Addr)
	}
	if req != nil {
		grpclog.Infof("req: %v", req)
	}
	rsp = &mygrpc.Response{
		Message: "200 OK",
	}
	return
}

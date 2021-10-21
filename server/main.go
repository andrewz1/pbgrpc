package main

import (
	"fmt"
	"net"

	pb "github.com/kogonia/protobuf_grpc/server/gen"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/grpclog"
)

func main() {
	listener, err := net.Listen("tcp", ":5300")

	if err != nil {
		grpclog.Fatalf("failed to listen: %v", err)
	}
	grpcServer := grpc.NewServer()

	pb.RegisterReverseServer(grpcServer, &server{})
	err = grpcServer.Serve(listener)
	if err != nil {
		grpclog.Fatalf("failed to serve grpcServer: %v", err)
	}
}

type server struct{}

func (s *server) Do(_ context.Context, request *pb.Request) (response *pb.Response, err error) {
	fmt.Printf("request from %s\n", request.Message)

	response = &pb.Response{
		Message: "200 OK",
	}

	return response, nil
}

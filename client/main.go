package main

import (
	"context"
	"fmt"
	"strconv"
	"sync"
	"time"

	pb "github.com/kogonia/protobuf_grpc/server/gen"
	"google.golang.org/grpc"
	"google.golang.org/grpc/grpclog"
)

func main() {
	count := 10000
	wg := &sync.WaitGroup{}
	wg.Add(count)

	for i := 0; i < count; i++ {
		id := strconv.Itoa(i)
		go func(id string, wg *sync.WaitGroup) {
			defer wg.Done()

			opts := []grpc.DialOption{
				grpc.WithInsecure(),
				grpc.WithTimeout(time.Second),
			}
			conn, err := grpc.Dial("127.0.0.1:5300", opts...)
			if err != nil {
				grpclog.Fatalf("fail to dial: %v", err)
			}
			defer conn.Close()

			client := pb.NewReverseClient(conn)
			request := &pb.Request{
				Message: id,
			}
			response, err := client.Do(context.Background(), request)

			if err != nil {
				grpclog.Fatalf("fail to dial: %v", err)
			}
			fmt.Printf("[%s]\t%s\n", id, response.Message)
		}(id, wg)
	}
	wg.Wait()
	fmt.Println("Done")
}

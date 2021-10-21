package main

import (
	"context"
	"fmt"
	"sync"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/grpclog"

	"github.com/andrewz1/pbgrpc/mygrpc"
)

const (
	target = "127.0.0.1:5300"
	count  = 10000
)

type rclient struct {
	cc *grpc.ClientConn
	wg *sync.WaitGroup
}

func (c *rclient) run(n int) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer func() {
		cancel()
		c.wg.Done()
	}()
	cl := mygrpc.NewReverseClient(c.cc)
	req := mygrpc.Request{Message: fmt.Sprintf("request %d", n)}
	if _, err := cl.Do(ctx, &req, grpc.WaitForReady(true)); err != nil {
		grpclog.Errorf("do: %v", err)
	}
}

func main() {
	cc, err := grpc.Dial(target, grpc.WithInsecure())
	if err != nil {
		grpclog.Fatalf("dial: %v", err)
	}
	var wg sync.WaitGroup
	wg.Add(count)
	rc := rclient{
		cc: cc,
		wg: &wg,
	}
	for i := 0; i < count; i++ {
		go rc.run(i)
	}
	wg.Wait()
	grpclog.Infof("done")
}

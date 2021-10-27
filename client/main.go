package main

import (
	"context"
	"fmt"
	"log"
	"sync"
	"sync/atomic"
	"time"

	"google.golang.org/grpc"

	"github.com/andrewz1/pbgrpc/mygrpc"
)

const (
	target = "10.0.10.100:5300"
	count  = 10000
)

type rclient struct {
	cc *grpc.ClientConn
	wg *sync.WaitGroup
	ee uint32
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
		log.Printf("do: %v\n", err)
		atomic.AddUint32(&c.ee, 1)
	}
}

func main() {
	cc, err := grpc.Dial(target, grpc.WithInsecure())
	if err != nil {
		log.Printf("dial: %v\n", err)
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
	log.Printf("done, errors: %d\n", rc.ee)
}

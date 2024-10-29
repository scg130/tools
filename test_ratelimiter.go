package main

import (
	"fmt"
	"sync"
	"sync/atomic"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/scg130/tools/rate_limit"
)

var rd = redis.NewClient(&redis.Options{
	Addr:        "localhost:6379",
	Password:    "smd013012",
	IdleTimeout: 3500,
	PoolSize:    50,
})

func main() {
	for i := 1; i <= 3; i++ {
		test()
	}
}

func test() {
	var pass, block int32

	rl := rate_limit.NewCRL("test2", time.Second*5, 100, rd)
	wg := sync.WaitGroup{}

	for i := 1; i <= 5000; i++ {
		wg.Add(1)
		go func() {
			// defer rl.Done()
			defer wg.Done()
			if allow, _ := rl.Allow(); allow {
				atomic.AddInt32(&pass, 1)
			} else {
				atomic.AddInt32(&block, 1)
			}
		}()
	}
	wg.Wait()
	fmt.Println(pass, block)
}

package redisync

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/go-redis/redis/v8"
)

var rc *redis.Client

func init() {
	_, exists := os.LookupEnv("REDIS_URL")
	if !exists {
		panic("Must set: $REDIS_URL")
	}
	options, err := redis.ParseURL(os.Getenv("REDIS_URL"))
	rc = redis.NewClient(options)
	if err != nil {
		panic("Unable to connect to: $REDIS_URL")
	}
}

func ExampleLock() {
	ctx := context.Background()
	ttl := time.Second
	m := NewMutex("my-lock", ttl)
	m.Lock(ctx, rc)
	defer m.Unlock(ctx, rc)

	done := make(chan bool)
	expired := make(chan bool)

	go func(e chan bool) {
		time.Sleep(ttl)
		e <- true
	}(expired)

	go func(d chan bool) {
		fmt.Printf("at=critical-section\n")
		d <- true
	}(done)

	select {
	case <-done:
		fmt.Printf("Finished.\n")
	case <-expired:
		fmt.Printf("Expired.\n")
	}
	// Output:
	// at=critical-section
	// Finished.
}

package redisync

import (
	"context"
	"os"
	"testing"
	"time"

	"github.com/go-redis/redis/v8"
)

func newConn() (*redis.Client, error) {
	options, err := redis.ParseURL(os.Getenv("REDIS_URL"))
	c := redis.NewClient(options)
	if err != nil {
		return nil, err
	}
	return c, nil
}

func TestLock(t *testing.T) {
	ctx := context.Background()
	rc, err := newConn()
	if err != nil {
		t.Error(err)
		t.FailNow()
	}
	defer rc.Close()
	ttl := time.Second
	m := NewMutex("redisync.test.1", ttl)
	m.Lock(ctx, rc)
	time.Sleep(ttl)
	ok := m.TryLock(ctx, rc)
	if !ok {
		t.Error("Expected mutex to be lockable.")
		t.FailNow()
	}
	m.Unlock(ctx, rc)
}

func TestLockLocked(t *testing.T) {
	ctx := context.Background()
	rc, err := newConn()
	if err != nil {
		t.Error(err)
		t.FailNow()
	}
	defer rc.Close()

	ttl := time.Second
	m1 := NewMutex("redisync.test.1", ttl)
	if ok := m1.TryLock(ctx, rc); !ok {
		t.Error("Expected mutex to be lockable.")
		t.FailNow()
	}

	m2 := NewMutex("redisync.test.1", ttl)
	if ok := m2.TryLock(ctx, rc); ok {
		t.Error("Expected mutex not to be lockable.")
		t.FailNow()
	}
	m1.Unlock(ctx, rc)
}

func TestUnlockOtherLocked(t *testing.T) {
	ctx := context.Background()
	rc, err := newConn()
	if err != nil {
		t.Error(err)
		t.FailNow()
	}
	defer rc.Close()

	ttl := time.Second
	m1 := NewMutex("redisync.test.1", ttl)
	if ok := m1.TryLock(ctx, rc); !ok {
		t.Error("Expected mutex to be lockable.")
		t.FailNow()
	}

	m2 := NewMutex("redisync.test.1", ttl)
	if ok, _ := m2.Unlock(ctx, rc); ok {
		t.Error("Expected mutex not to be unlockable.")
		t.FailNow()
	}
	m1.Unlock(ctx, rc)
}

func TestLockExpired(t *testing.T) {
	ctx := context.Background()
	rc, err := newConn()
	if err != nil {
		t.Error(err)
		t.FailNow()
	}
	defer rc.Close()

	ttl := time.Second
	m1 := NewMutex("redisync.test.1", ttl)
	if ok := m1.TryLock(ctx, rc); !ok {
		t.Error("Expected mutex to be lockable.")
		t.FailNow()
	}
	time.Sleep(ttl)

	m2 := NewMutex("redisync.test.1", ttl)
	if ok := m2.TryLock(ctx, rc); !ok {
		t.Error("Expected mutex to be lockable.")
		t.FailNow()
	}
	m2.Unlock(ctx, rc)
}

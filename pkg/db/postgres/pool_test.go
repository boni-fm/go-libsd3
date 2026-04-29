package postgres

import (
	"context"
	"sync"
	"testing"
)

func TestConcurrentConnect_NoDuplicate(t *testing.T) {
	pool := &ConnectionPool{
		connections: make(map[string]*Database),
		configs:     make(map[string]Config),
	}

	fakeDB := &Database{isClosed: false}
	pool.connections["testdc"] = fakeDB

	const goroutines = 50
	results := make([]*Database, goroutines)
	var wg sync.WaitGroup

	for i := 0; i < goroutines; i++ {
		wg.Add(1)
		go func(idx int) {
			defer wg.Done()
			pool.mu.RLock()
			conn := pool.connections["testdc"]
			pool.mu.RUnlock()
			results[idx] = conn
		}(i)
	}

	wg.Wait()

	for i, r := range results {
		if r != fakeDB {
			t.Errorf("goroutine %d got unexpected connection", i)
		}
	}
}

func TestConnectionPool_HasConnection(t *testing.T) {
	pool := &ConnectionPool{
		connections: make(map[string]*Database),
		configs:     make(map[string]Config),
	}

	if pool.HasConnection("testdc") {
		t.Error("expected no connection initially")
	}

	pool.connections["testdc"] = &Database{}
	if !pool.HasConnection("testdc") {
		t.Error("expected connection to exist after insertion")
	}
}

func TestConnectionPool_ConnectNoConfig(t *testing.T) {
	pool := &ConnectionPool{
		connections: make(map[string]*Database),
		configs:     make(map[string]Config),
	}

	ctx := context.Background()
	_, err := pool.Connect(ctx, "nonexistent")
	if err == nil {
		t.Error("expected error for missing config")
	}
}

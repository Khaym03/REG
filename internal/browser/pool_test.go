package browser

import (
	"context"
	"sync"
	"testing"
	"time"
)

func newMockFactory(t *testing.T) BrowserRuntimeFactory {
	t.Helper()

	return func(context.Context) (BrowserRuntime, error) {
		return NewMockBrowserRuntime(t), nil
	}
}

func waitOrFail(t *testing.T, done <-chan struct{}, msg string) {
	t.Helper()

	select {
	case <-done:
	case <-time.After(100 * time.Millisecond):
		t.Fatal(msg)
	}
}

func TestAcquireCreatesRuntime(t *testing.T) {
	t.Parallel()

	pool := NewBrowserPool(BrowserPoolConfig{
		factory: newMockFactory(t),
	})

	rt, err := pool.Acquire(t.Context())
	if err != nil {
		t.Fatalf("Acquire() returned error: %v", err)
	}

	if rt == nil {
		t.Fatal("Acquire() returned nil runtime")
	}
}

func TestAcquireReusesReleasedRuntime(t *testing.T) {
	t.Parallel()

	pool := NewBrowserPool(BrowserPoolConfig{
		factory: newMockFactory(t),
	})

	rt1, err := pool.Acquire(t.Context())
	if err != nil {
		t.Fatal(err)
	}

	pool.Release(rt1)

	rt2, err := pool.Acquire(t.Context())
	if err != nil {
		t.Fatal(err)
	}

	if rt1 != rt2 {
		t.Fatal("expected released runtime to be reused")
	}
}

func TestReleaseWhenPoolIsFullDoesNotBlock(t *testing.T) {
	t.Parallel()

	pool := NewBrowserPool(BrowserPoolConfig{
		factory: newMockFactory(t),
	})

	rt1 := NewMockBrowserRuntime(t)
	rt2 := NewMockBrowserRuntime(t)

	rt2.On("Close").Return(nil).Once()

	pool.Release(rt1)

	done := make(chan struct{})

	go func() {
		defer close(done)
		pool.Release(rt2)
	}()

	waitOrFail(t, done, "Release blocked")
}

func TestAcquireNeverBlocks(t *testing.T) {
	t.Parallel()

	pool := NewBrowserPool(BrowserPoolConfig{
		factory: newMockFactory(t),
	})

	done := make(chan struct{})

	go func() {
		defer close(done)

		for range 5 {
			if _, err := pool.Acquire(t.Context()); err != nil {
				t.Errorf("Acquire() returned error: %v", err)
				return
			}
		}
	}()

	waitOrFail(t, done, "Acquire blocked")
}

func TestStressConcurrentAcquireRelease(t *testing.T) {
	t.Parallel()

	const (
		poolSize   = 100
		workers    = 100
		iterations = 1000
	)

	pool := NewBrowserPool(BrowserPoolConfig{
		size:    poolSize,
		factory: newMockFactory(t),
	})

	var wg sync.WaitGroup

	for range workers {
		wg.Go(func() {
			for range iterations {
				rt, err := pool.Acquire(t.Context())
				if err != nil {
					t.Errorf("Acquire() returned error: %v", err)
					return
				}

				pool.Release(rt)
			}
		})
	}

	wg.Wait()
}

package utilworker_test

import (
	"context"
	"sync"
	"testing"
	"time"

	"github.com/kiln-mid/pkg/utilworker"
	"github.com/stretchr/testify/assert"
)

func TestNewIntervalWorker(t *testing.T) {
	fn := func(wg *sync.WaitGroup) {
		wg.Done()
	}

	worker := utilworker.NewIntervalWorker("testWorker", fn, 5*time.Second)
	assert.Equal(t, "testWorker", worker.Name)
	assert.Equal(t, 5*time.Second, worker.Interval)
}

func TestWorker_StartAndStop(t *testing.T) {
	fnCalled := false
	fn := func(wg *sync.WaitGroup) {
		defer wg.Done()
		fnCalled = true
	}

	worker := utilworker.NewIntervalWorker("testWorker", fn, 1*time.Second)
	ctx, cancel := context.WithCancel(context.Background())
	go worker.Start(ctx)

	time.Sleep(1500 * time.Millisecond)
	assert.True(t, fnCalled)

	// Cancel the context to stop the worker
	cancel()

	fnCalled = false

	// Ensure that the worker no longer runs the function
	time.Sleep(2 * time.Second)
	assert.False(t, fnCalled)
}

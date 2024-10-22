package utilworker_test

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/kiln-mid/pkg/utilworker"
	"github.com/stretchr/testify/assert"
)

func TestWorker_StartAndStop(t *testing.T) {
	fnCalled := false
	fn := func(ctx context.Context) error {
		fnCalled = true
		return nil
	}

	ctx, cancel := context.WithCancel(context.Background())
	go utilworker.StartNewIntervalWorker("testWorker", fn, 5*time.Millisecond, ctx)

	time.Sleep(50 * time.Millisecond)
	assert.True(t, fnCalled)

	cancel()

	fnCalled = false

	time.Sleep(50 * time.Millisecond)

	assert.False(t, fnCalled)
}

func TestWorker_StartWithError(t *testing.T) {
	fnCalled := false
	fn := func(ctx context.Context) error {
		fnCalled = true
		return fmt.Errorf("AN ERROR OCCURED")
	}

	go utilworker.StartNewIntervalWorker("testWorker", fn, 5*time.Millisecond, context.Background())

	time.Sleep(50 * time.Millisecond)
	assert.True(t, fnCalled)

	fnCalled = false

	time.Sleep(50 * time.Millisecond)

	assert.False(t, fnCalled)
}

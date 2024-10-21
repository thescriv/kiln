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
	fn := func() error {
		fnCalled = true
		return nil
	}

	ctx, cancel := context.WithCancel(context.Background())
	go utilworker.StartNewIntervalWorker("testWorker", fn, 1*time.Second, ctx)

	time.Sleep(1500 * time.Millisecond)
	assert.True(t, fnCalled)

	cancel()

	fnCalled = false

	time.Sleep(2 * time.Second)

	assert.False(t, fnCalled)
}

func TestWorker_StartWithError(t *testing.T) {
	fnCalled := false
	fn := func() error {
		fnCalled = true
		return fmt.Errorf("AN ERROR OCCURED")
	}

	go utilworker.StartNewIntervalWorker("testWorker", fn, 1*time.Second, context.Background())

	time.Sleep(1500 * time.Millisecond)
	assert.True(t, fnCalled)

	fnCalled = false

	time.Sleep(2 * time.Second)

	assert.False(t, fnCalled)
}

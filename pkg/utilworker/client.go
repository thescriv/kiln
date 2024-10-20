package utilworker

import (
	"context"
	"fmt"
	"sync"
	"time"
)

var wg sync.WaitGroup

// DefaultWorkerInterval is the default interval duration used by TimeoutClient when no duration is provided.
const DefaultWorkerInterval = 10 * time.Second

type Worker struct {
	Name     string
	Function func(*sync.WaitGroup)
	Interval time.Duration
}

// NewIntervalWorker return a new intervalWorker take in params a name to handle logging, the called function and an interval value.
// the name will help the logging system.
// The worker will be called each X seconds, based on the interval value provided by defautl the timer is set to every 10seconds.
// the function provided should received a *sync.WaitGroup in param as the worker once started will use the system of waitGroup and will Add and Wait.
func NewIntervalWorker(name string, function func(*sync.WaitGroup), interval time.Duration) *Worker {
	if interval == 0 {
		interval = DefaultWorkerInterval
	}

	return &Worker{
		Name:     name,
		Function: function,
		Interval: interval,
	}
}

// Start launch the worker and call
func (w *Worker) Start(ctx context.Context) {
	fmt.Printf("[WORKER] %s started.\n", w.Name)
	go func() {
		ticker := time.NewTicker(w.Interval)
		defer ticker.Stop()

		for {
			select {
			case <-ticker.C:
				fmt.Printf("[WORKER] %s called.\n", w.Name)

				wg.Add(1)

				go w.Function(&wg)

				wg.Wait()
			case <-ctx.Done(): // Arrêter si le contexte est annulé
				fmt.Printf("[WORKER] %s stopped.\n", w.Name)
				return
			}
		}
	}()
}

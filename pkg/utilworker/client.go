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

// StartNewIntervalWorker start a new worker called at an interval which is provided in params
// worker name is used only for logging purpose.
// fct represent the function called inside the worker.
// a context should also be passed to cancel the worker.
func StartNewIntervalWorker(name string, fct func(context.Context) error, interval time.Duration, ctx context.Context) {
	if interval == 0 {
		interval = DefaultWorkerInterval
	}

	fmt.Printf("[WORKER] %s started.\n", name)
	go func() {
		cancel := make(chan bool)
		defer close(cancel)

		ticker := time.NewTicker(interval)
		defer ticker.Stop()

		for {
			select {
			case <-ticker.C:
				fmt.Printf("[WORKER] %s called.\n", name)

				wg.Add(1)

				go func() {
					defer wg.Done()
					err := fct(ctx)
					if err != nil {
						fmt.Println(err.Error())

						cancel <- true
					}
				}()

				wg.Wait()
			case <-cancel:
			case <-ctx.Done():
				fmt.Printf("[WORKER] %s stopped.\n", name)
				return
			}
		}
	}()
}

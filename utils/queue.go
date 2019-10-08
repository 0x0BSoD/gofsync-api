package utils

import (
	"runtime"
)

// =====================================================================================================================
// Goroutine workers
// =====================================================================================================================
// WorkQueue is a channel type that you can send Work on.
type WorkQueue chan Work

// New creates a WorkQueue with runtime.NumCPU() workers.
func New() WorkQueue {
	return NewN(runtime.NumCPU() + 1)
}

// NewN creates and returns a new WorkQueue that has the specified number
// of workers.
func NewN(numWorkers int) WorkQueue {
	queue := make(WorkQueue)
	d := make(dispatcher, numWorkers)
	go d.dispatch(queue)
	return queue
}

// Work is a task to perform that can be sent over a WorkQueue.
type Work func()

type dispatcher chan chan Work

func (d dispatcher) dispatch(queue WorkQueue) {
	// Create and start all of our workers.
	for i := 0; i < cap(d); i++ {
		w := make(worker)
		go w.work(d)
	}

	// Start the main loop in a goroutine.
	go func() {
		for work := range queue {
			go func(work Work) {
				worker := <-d
				worker <- work
			}(work)
		}

		// If we get here, the work queue has been closed, and we should
		// stop all the workers.
		for i := 0; i < cap(d); i++ {
			w := <-d
			close(w)
		}
	}()
}

type worker chan Work

func (w worker) work(d dispatcher) {
	// Add ourselves to the dispatcher.
	d <- w

	// Start the main loop.
	go w.wait(d)
}

func (w worker) wait(d dispatcher) {
	for work := range w {
		// Do the work.
		if work == nil {
			panic("nil work received")
		}

		work()

		// Re-add ourselves to the dispatcher.
		d <- w
	}
}

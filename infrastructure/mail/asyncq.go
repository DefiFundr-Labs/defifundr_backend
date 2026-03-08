package mail

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/demola234/defifundr/infrastructure/common/logging"
)

// AsyncQueue is a simple in-memory queue for async processing
type AsyncQueue struct {
	queue     chan any
	logger    logging.Logger
	workers   int
	wg        sync.WaitGroup
	processor func(any) error
	stopCh    chan struct{}
	mu        sync.Mutex
	running   bool
}

// NewAsyncQueue creates a new async queue
func NewAsyncQueue(capacity int, workers int, logger logging.Logger, processor func(any) error) *AsyncQueue {
	return &AsyncQueue{
		queue:     make(chan any, capacity),
		logger:    logger,
		workers:   workers,
		processor: processor,
		stopCh:    make(chan struct{}),
	}
}

// Start starts the queue workers
func (q *AsyncQueue) Start() {
	q.mu.Lock()
	defer q.mu.Unlock()

	if q.running {
		return
	}

	q.running = true
	q.stopCh = make(chan struct{})

	for i := 0; i < q.workers; i++ {
		q.wg.Add(1)
		go q.worker(i)
	}

	q.logger.Info("AsyncQueue started", map[string]any{
		"workers":  q.workers,
		"capacity": cap(q.queue),
	})
}

// Stop stops the queue workers
func (q *AsyncQueue) Stop() {
	q.mu.Lock()
	if !q.running {
		q.mu.Unlock()
		return
	}
	q.running = false
	close(q.stopCh)
	q.mu.Unlock()

	q.wg.Wait()
	q.logger.Info("AsyncQueue stopped")
}

// Enqueue adds an item to the queue
func (q *AsyncQueue) Enqueue(item any) error {
	select {
	case q.queue <- item:
		return nil
	case <-time.After(2 * time.Second):
		return ErrQueueFull
	}
}

// EnqueueWithContext adds an item to the queue with context
func (q *AsyncQueue) EnqueueWithContext(ctx context.Context, item any) error {
	select {
	case q.queue <- item:
		return nil
	case <-ctx.Done():
		return ctx.Err()
	case <-time.After(2 * time.Second):
		return ErrQueueFull
	}
}

// worker processes items from the queue
func (q *AsyncQueue) worker(id int) {
	defer q.wg.Done()

	q.logger.Info("AsyncQueue worker started", map[string]any{
		"worker_id": id,
	})

	for {
		select {
		case <-q.stopCh:
			q.logger.Info("AsyncQueue worker stopping", map[string]any{
				"worker_id": id,
			})
			return
		case item, ok := <-q.queue:
			if !ok {
				q.logger.Info("AsyncQueue channel closed", map[string]any{
					"worker_id": id,
				})
				return
			}

			err := q.processor(item)
			if err != nil {
				q.logger.Error("Failed to process item", err, map[string]any{
					"worker_id": id,
				})
			}
		}
	}
}

// ErrQueueFull is returned when the queue is full
var ErrQueueFull = fmt.Errorf("queue is full")
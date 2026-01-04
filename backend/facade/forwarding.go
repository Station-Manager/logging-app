package facade

import (
	"context"
	"fmt"
	"sync"
	"sync/atomic"
	"time"

	"github.com/Station-Manager/errors"
	"github.com/Station-Manager/logging"
	"github.com/Station-Manager/types"
)

type forwarding struct {
	pollInterval    time.Duration
	maxWorkers      int
	forwardingQueue chan types.QsoUpload
	dbWriteQueue    chan func() error
	fetchPending    func() ([]types.QsoUpload, error)     // See: s.DatabaseService.FetchPendingUploads()
	sendAndMarkDone func(qsoUpload types.QsoUpload) error // See: s.forwardQso(qsoUpload)
	logger          *logging.Service

	// Worker lifecycle management
	wg             sync.WaitGroup
	workerRegistry sync.Map // map[string]bool - tracks running workers
	started        atomic.Bool
	stopping       atomic.Bool
}

// start initializes and starts the worker and polling goroutines for the forwarding process. Returns an error if context is nil.
func (f *forwarding) start(ctx context.Context, shutdown <-chan struct{}) error {
	const op errors.Op = "forwarding.start"
	if ctx == nil {
		return errors.New(op).Msg("Context is nil")
	}

	if !f.started.CompareAndSwap(false, true) {
		return errors.New(op).Msg("Forwarding already started")
	}

	// Start worker goroutines with tracking
	for i := 0; i < f.maxWorkers; i++ {
		workerID := i
		workerName := fmt.Sprintf("worker-%d", workerID)

		f.launchWorker(ctx, shutdown, workerName, func() {
			f.workerLoop(ctx, shutdown, workerID)
		})
	}

	// Start the database write worker
	f.launchWorker(ctx, shutdown, "db-writer", func() {
		f.dbWriteWorkerLoop(ctx, shutdown)
	})

	// Start the polling goroutine
	f.launchWorker(ctx, shutdown, "poller", func() {
		f.pollerLoop(ctx, shutdown)
	})

	f.logger.InfoWith().
		Int("workers", f.maxWorkers+2).
		Msg("All forwarding workers launched")

	return nil
}

// launchWorker starts a worker goroutine with proper lifecycle tracking
func (f *forwarding) launchWorker(ctx context.Context, shutdown <-chan struct{}, name string, work func()) {
	// Register worker before launching
	f.workerRegistry.Store(name, true)

	f.wg.Add(1)
	go func() {
		defer func() {
			// Cleanup on exit
			f.workerRegistry.Delete(name)
			f.wg.Done()
			f.logger.DebugWith().Str("worker", name).Msg("Worker exited")

			// Recover from panics
			if r := recover(); r != nil {
				f.logger.ErrorWith().
					Str("worker", name).
					Interface("panic", r).
					Msg("Worker panicked")
			}
		}()

		f.logger.DebugWith().Str("worker", name).Msg("Worker started")

		// Run the actual work
		work()
	}()
}

// ActiveWorkerCount returns the number of currently running workers
func (f *forwarding) ActiveWorkerCount() int {
	count := 0
	f.workerRegistry.Range(func(key, value interface{}) bool {
		count++
		return true
	})
	return count
}

// stop gracefully shuts down all workers with timeout
func (f *forwarding) stop(timeout time.Duration) error {
	if !f.stopping.CompareAndSwap(false, true) {
		return nil // Already stopping
	}

	f.logger.InfoWith().Msg("Stopping forwarding workers")

	// Close queues to signal workers to exit
	close(f.forwardingQueue)
	close(f.dbWriteQueue)

	// Wait with timeout
	done := make(chan struct{})
	go func() {
		f.wg.Wait()
		close(done)
	}()

	select {
	case <-done:
		f.logger.InfoWith().Msg("All forwarding workers stopped gracefully")
		return nil
	case <-time.After(timeout):
		activeWorkers := f.ActiveWorkerCount()
		f.logger.ErrorWith().
			Int("stuck_workers", activeWorkers).
			Msg("Timeout waiting for workers to stop")

		// Log which workers are still running
		f.workerRegistry.Range(func(key, value interface{}) bool {
			f.logger.WarnWith().
				Str("worker", key.(string)).
				Msg("Worker still running after timeout")
			return true
		})

		return errors.New("forwarding.stop").
			Msgf("Timeout after %v with %d workers still running", timeout, activeWorkers)
	}
}

// pollerLoop starts a loop that periodically fetches pending QSO uploads and attempts to enqueue them for processing.
func (f *forwarding) pollerLoop(ctx context.Context, shutdown <-chan struct{}) {
	f.logger.InfoWith().Msg("Starting forwarding poller")

	ticker := time.NewTicker(f.pollInterval)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-shutdown:
			return
		case <-ticker.C:
			qsoUploads, err := f.fetchPending()
			if err != nil {
				f.logger.ErrorWith().Err(err).Msg("Failed to fetch pending uploads")
				continue
			}
			for _, qsoUpload := range qsoUploads {
				select {
				case f.forwardingQueue <- qsoUpload:
					// forwarded to the forwarding queue
				case <-ctx.Done():
					return
				case <-shutdown:
					return
				default:
					f.logger.WarnWith().
						Int64("upload_id", qsoUpload.ID).
						Msg("Forwarding queue full, dropping upload")
				}
			}
		}
	}
}

// workerLoop runs a worker goroutine to process QSO uploads from the forwarding queue until shutdown or context cancellation.
func (f *forwarding) workerLoop(ctx context.Context, shutdown <-chan struct{}, workerID int) {
	f.logger.InfoWith().Int("workerID", workerID).Msg("Starting forwarding worker")

	for {
		select {
		case <-ctx.Done():
			f.logger.InfoWith().Msg("Context done, shutting down forwarding worker")
			return
		case <-shutdown:
			f.logger.InfoWith().Msg("Forwarding worker shutting down")
			return
		case qsoUpload, ok := <-f.forwardingQueue:
			if !ok {
				return
			}

			// Do network call (can be concurrent)
			err := f.sendAndMarkDone(qsoUpload)

			// Note: Database writes are now handled within sendAndMarkDone via the dbWriteQueue
			// This maintains backward compatibility while ensuring serialized DB access
			if err != nil {
				f.logger.ErrorWith().Err(err).Msg("Error in processing or forwarding QSO")
			}
		}
	}
}

func (f *forwarding) dbWriteWorkerLoop(ctx context.Context, shutdown <-chan struct{}) {
	f.logger.InfoWith().Msg("Starting database write worker")

	for {
		select {
		case <-ctx.Done():
			return
		case <-shutdown:
			return
		case writeOp, ok := <-f.dbWriteQueue:
			if !ok {
				return
			}
			if err := writeOp(); err != nil {
				f.logger.ErrorWith().Err(err).Msg("Database write operation failed")
			}
		}
	}
}

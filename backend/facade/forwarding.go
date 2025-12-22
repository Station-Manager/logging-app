package facade

import (
	"context"
	"sync"
	"time"

	"github.com/Station-Manager/errors"
	"github.com/Station-Manager/logging"
	"github.com/Station-Manager/types"
)

type forwarding struct {
	pollInterval    time.Duration
	maxWorkers      int
	forwardingQueue chan types.QsoUpload
	fetchPending    func() ([]types.QsoUpload, error)     // See: s.DatabaseService.FetchPendingUploads()
	sendAndMarkDone func(qsoUpload types.QsoUpload) error // See: s.forwardQso(qsoUpload)
	logger          *logging.Service
	wg              sync.WaitGroup
}

// start initializes and starts the worker and polling goroutines for the forwarding process. Returns an error if context is nil.
func (f *forwarding) start(ctx context.Context, shutdown <-chan struct{}) error {
	const op errors.Op = "forwarding.start"
	if ctx == nil {
		return errors.New(op).Msg("Context is nil")
	}

	// Add to WaitGroup BEFORE starting goroutines to prevent race condition where the parent calls wg.Wait()
	// before a child calls wd.Add(...)
	f.wg.Add(f.maxWorkers + 1) // +1 for poller

	// Start worker goroutines
	for i := 0; i < f.maxWorkers; i++ {
		go f.workerLoop(ctx, shutdown, i)
	}

	// Start the polling goroutine
	go f.pollerLoop(ctx, shutdown)

	return nil
}

// pollerLoop starts a loop that periodically fetches pending QSO uploads and attempts to enqueue them for processing.
func (f *forwarding) pollerLoop(ctx context.Context, shutdown <-chan struct{}) {
	defer f.wg.Done()

	f.logger.DebugWith().Msg("Starting forwarding poller")

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
	defer f.wg.Done()

	f.logger.InfoWith().Int("workerID", workerID).Msg("Starting forwarding worker")

	for {
		select {
		case <-ctx.Done():
			f.logger.InfoWith().Msg("Context done, shutting down forwarding poller")
			return
		case <-shutdown:
			f.logger.InfoWith().Msg("Forwarding poller shutting down")
			return
		case qsoUpload, ok := <-f.forwardingQueue:
			if !ok {
				return
			}

			if err := f.sendAndMarkDone(qsoUpload); err != nil {
				f.logger.ErrorWith().Err(err).Msg("Failed to forward QSO")
			}
		}
	}
}

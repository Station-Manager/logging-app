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
	fetchPending    func() ([]types.QsoUpload, error)
	sendAndMarkDone func(qsoUpload types.QsoUpload) error
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
			f.logger.DebugWith().Msg("Context done, shutting down forwarding poller")
			return
		case <-shutdown:
			f.logger.DebugWith().Msg("Forwarding poller shutting down")
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
					// sent successfully
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

	f.logger.DebugWith().Int("workerID", workerID).Msg("Starting forwarding worker")

	for {
		select {
		case <-ctx.Done():
			f.logger.DebugWith().Msg("Context done, shutting down forwarding poller")
			return
		case <-shutdown:
			f.logger.DebugWith().Msg("Forwarding poller shutting down")
			return
		case qsoUpload, ok := <-f.forwardingQueue:
			if !ok {
				return
			}

			f.logger.DebugWith().
				Int64("upload_id", qsoUpload.ID).
				//				Int64("qso_id", qsoUpload.Qso.ID).
				//				Str("callsign", qsoUpload.Qso.Call).
				Str("service", qsoUpload.Service).
				Int("workerID", workerID).
				Msg("Processing QSO upload for forwarding")

			if err := f.sendAndMarkDone(qsoUpload); err != nil {
				f.logger.ErrorWith().Err(err).Msg("Failed to forward QSO")
			}
		}
	}
}

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

func (f *forwarding) start(ctx context.Context, shutdown <-chan struct{}) error {
	const op errors.Op = "forwarding.start"
	if ctx == nil {
		return errors.New(op).Msg("Context is nil")
	}

	// Start worker goroutines
	for i := 0; i < f.maxWorkers; i++ {
		go f.workerLoop(ctx, shutdown, i)
	}

	// Start the polling goroutine
	go f.pollerLoop(ctx, shutdown)

	return nil
}

func (f *forwarding) pollerLoop(ctx context.Context, shutdown <-chan struct{}) {
	f.logger.DebugWith().Msg("Starting forwarding poller")

	f.wg.Add(1)
	defer f.wg.Done()

	for {
		select {
		case <-ctx.Done():
			f.logger.DebugWith().Msg("Context done, shutting down forwarding poller")
			return
		case <-shutdown:
			f.logger.DebugWith().Msg("Forwarding poller shutting down")
			return
		default:
			qsoUploads, err := f.fetchPending()
			if err != nil {
				f.logger.ErrorWith().Err(err).Msg("Failed to fetch pending uploads")
				continue
			}
			for _, qsoUpload := range qsoUploads {
				f.forwardingQueue <- qsoUpload
			}
		}
	}
}

func (f *forwarding) workerLoop(ctx context.Context, shutdown <-chan struct{}, workerID int) {
	f.logger.DebugWith().Int("workerID", workerID).Msg("Starting forwarding worker")

	f.wg.Add(1)
	defer f.wg.Done()

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
				Int64("qso_id", qsoUpload.Qso.ID).
				Str("callsign", qsoUpload.Qso.Call).
				Str("service", qsoUpload.Service).
				Int("workerID", workerID).
				Msg("Processing QSO upload for forwarding")

			// sendAndMarkDone() here
		}
	}
}

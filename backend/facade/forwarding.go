package facade

import (
	"context"
	"fmt"
	"time"

	"github.com/Station-Manager/errors"
	"github.com/Station-Manager/types"
)

type forwarding struct {
	pollInterval    time.Duration
	maxWorkers      int
	queue           chan types.Qso
	fetchPending    func(ctx context.Context) ([]types.QsoUpload, error)
	sendAndMarkDone func(ctx context.Context, qso types.QsoUpload) error
}

func (f *forwarding) start(ctx context.Context) error {
	const op errors.Op = "forwarding.start"
	if ctx == nil {
		return errors.New(op).Msg("Context is nil")
	}

	for i := 0; i < f.maxWorkers; i++ {
		go f.workerLoop(ctx)
	}

	return nil
}

func (f *forwarding) workerLoop(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			return
		case qso, ok := <-f.queue:
			if !ok {
				return
			}
			fmt.Println("Forwarding QSO:", qso)
		}
	}
}

// qsoForwarder forwards queued QSOs until a shutdown signal is received.
// When a QSO is logged locally, it is placed into a queue to be forwarded to
// any configured external logging services. This function processes that queue and runs as
// a goroutine.
func (s *Service) startForwarding(shutdown <-chan struct{}) {
	const op errors.Op = "facade.Service.startForwarding"

	if s.ctx == nil {
		err := errors.New(op).Msg("Service is not initialized.")
		s.LoggerService.ErrorWith().Err(err).Msg("Failed to start QSO forwarding.")
		return
	}

	readTicker := time.NewTicker(s.requiredCfgs.QsoForwardingIntervalSeconds * time.Second)
	defer readTicker.Stop()

	for {
		select {
		case <-shutdown:
			return
		case <-readTicker.C:
			s.LoggerService.DebugWith().Msg("Check for QSOs to be forwareded...")
			slice, err := s.DatabaseService.FetchPendingUploads()
			if err != nil {
				s.LoggerService.ErrorWith().Err(err).Msg("Failed to fetch QSOs to be forwarded.")
				continue
			}
			if len(slice) > 0 {
				s.LoggerService.DebugWith().Int("count", len(slice)).Msg("We have something to do...")
			}
		}
	}
}

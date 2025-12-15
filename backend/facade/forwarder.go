package facade

import (
	"time"
)

// qsoForwarder forwards queued QSOs until a shutdown signal is received.
// When a QSO is logged locally, it is placed into a queue to be forwarded to
// any configured external logging services. This function processes that queue and runs as
// a goroutine.
func (s *Service) qsoForwarder(shutdown <-chan struct{}) {
	readTicker := time.NewTicker(s.requiredCfgs.QsoForwardingIntervalSeconds * time.Second)
	defer readTicker.Stop()

	for {
		select {
		case <-shutdown:
			return
		case <-readTicker.C:
			s.LoggerService.DebugWith().Msg("Check for QSOs to be forwareded...")
			slice, err := s.DatabaseService.FetchQsoSliceNotForwarded()
			if err != nil {
				s.LoggerService.ErrorWith().Err(err).Msg("Failed to fetch QSOs to be forwarded.")
				continue
			}
			if slice != nil && len(slice) > 0 {
				s.LoggerService.DebugWith().Int("count", len(slice)).Msg("We have something to do...")
			}
		}
	}
}

package facade

import (
	"github.com/Station-Manager/types"
	"strings"
)

// launchWorkerThread starts a new goroutine for the given worker function and manages its lifecycle using a wait group.
func (s *Service) launchWorkerThread(run *runState, workerFunc func(<-chan struct{}), workerName string) {
	run.wg.Add(1)
	go func() {
		defer run.wg.Done()
		s.LoggerService.InfoWith().Str("worker", workerName).Msg("Facade worker starting")
		workerFunc(run.shutdownChannel)
		s.LoggerService.InfoWith().Str("worker", workerName).Msg("Facade worker stopped")
	}()
}

// parseCallsign processes and cleans up a callsign string, removing trailing known modifiers and excess whitespace.
func (s *Service) parseCallsign(callsign string) string {
	str := strings.TrimSpace(callsign)
	if str == "" {
		return str
	}
	parts := strings.Split(str, "/")
	if len(parts) == 1 {
		return str
	}
	// Known trailing modifiers to strip
	known := map[string]struct{}{
		"P":        {},
		"PORTABLE": {},
		"M":        {},
		"MM":       {},
		"MOBILE":   {},
		"QRP":      {},
		"QRO":      {},
		"AM":       {},
		"PM":       {},
	}
	for len(parts) > 1 { // keep at least one segment
		last := strings.ToUpper(strings.TrimSpace(parts[len(parts)-1]))
		if _, ok := known[last]; ok {
			parts = parts[:len(parts)-1]
			continue
		}
		break
	}
	result := strings.Join(parts, "/")
	return result
}

// lookupCallsignOnline fetches details about the specified callsign from an online service (QRZ.com, HamQTH, etc.).
// Returns a ContactedStation object or an error in case of failure.
func (s *Service) lookupCallsignOnline(callsign string) (types.ContactedStation, error) {

	return types.ContactedStation{}, nil
}

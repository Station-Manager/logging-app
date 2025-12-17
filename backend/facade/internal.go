package facade

import (
	"strconv"
	"strings"
	"time"

	"github.com/Station-Manager/errors"
	"github.com/Station-Manager/maidenhead"
	"github.com/Station-Manager/types"
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
	const op errors.Op = "facade.Service.lookupCallsignOnline"
	emptyRetVal := types.ContactedStation{}

	if !s.initialized.Load() {
		return emptyRetVal, errors.New(op).Msg("service is not initialized")
	}

	station, err := s.QrzLookupService.Lookup(callsign)
	if err != nil {
		return emptyRetVal, errors.New(op).Err(err).Msg("Failed to lookup callsign")
	}

	return station, nil
}

func (s *Service) calulatedBearingAndDistance(country *types.Country, ls types.LoggingStation, cs types.ContactedStation) error {
	const op errors.Op = "facade.Service.calulatedBearingAndDistance"
	if country == nil {
		return errors.New(op).Msg("country parameter is nil")
	}
	if ls.MyGridsquare == "" {
		return errors.New(op).Msg("logging station's gridsquare is empty")
	}
	if cs.Gridsquare == "" {
		return errors.New(op).Msg("contacted station's gridsquare is empty")
	}

	if location, err := maidenhead.GetLocation(ls.MyGridsquare, cs.Gridsquare); err != nil {
		s.LoggerService.WarnWith().Err(err).Msg("Failed to get location between logging station and contacted station")
	} else {
		country.ShortPathBearing = strconv.FormatFloat(location.ShortPathBearing, 'f', -1, 64)
		country.ShortPathDistance = strconv.Itoa(int(location.ShortPathDistanceKm))
		country.LongPathBearing = strconv.FormatFloat(location.LongPathBearing, 'f', -1, 64)
		country.LongPathDistance = strconv.Itoa(int(location.LongPathDistanceKm))
	}
	return nil
}

func (s *Service) initializeForwarding() error {
	s.forwarder = &forwarding{
		pollInterval:    30 * time.Second, // Will be updated in Start
		maxWorkers:      5,
		forwardingQueue: make(chan types.QsoUpload, 10),
		fetchPending: func() ([]types.QsoUpload, error) {
			return s.DatabaseService.FetchPendingUploads()
		},
		sendAndMarkDone: func(qsoUpload types.QsoUpload) error {
			return s.forwardQso(qsoUpload)
		},
		logger: s.LoggerService,
	}
	return nil
}

// forwardQso forwards a single QSO upload to the appropriate service and updates the database
func (s *Service) forwardQso(qsoUpload types.QsoUpload) error {
	const op errors.Op = "facade.Service.forwardQso"

	if s.container == nil {
		return errors.New(op).Msg("container is nil")
	}

	s.LoggerService.InfoWith().Str("service", qsoUpload.Service).Str("qso", qsoUpload.Qso.Call).Msg("Forwarding QSO to:")
	// Get the appropriate forwarder from the container based on the service name
	//forwarderBeanID := qsoUpload.Service + "forwarder"
	//forwarderInterface, err := s.container.ResolveSafe(forwarderBeanID)
	//if err != nil {
	//	// Increment attempt count and record error
	//	if updateErr := s.DatabaseService.UpdateQsoUploadStatus(
	//		qsoUpload.ID,
	//		"failed",
	//		qsoUpload.Attempts+1,
	//		"forwarder not found: "+err.Error(),
	//	); updateErr != nil {
	//		s.LoggerService.ErrorWith().Err(updateErr).Msg("Failed to update upload status after forwarder resolution failure")
	//	}
	//	return errors.New(op).Err(err).Msgf("failed to resolve forwarder for service: %s", qsoUpload.Service)
	//}
	//
	//forwarder, ok := forwarderInterface.(interface {
	//	Forward(qso types.Qso, param ...string) error
	//	IsEnabled() bool
	//})
	//if !ok {
	//	err := errors.New(op).Msgf("resolved bean is not a valid forwarder: %s", forwarderBeanID)
	//	if updateErr := s.DatabaseService.UpdateQsoUploadStatus(
	//		qsoUpload.ID,
	//		"failed",
	//		qsoUpload.Attempts+1,
	//		err.Error(),
	//	); updateErr != nil {
	//		s.LoggerService.ErrorWith().Err(updateErr).Msg("Failed to update upload status after type assertion failure")
	//	}
	//	return err
	//}
	//
	//// Check if forwarder is enabled
	//if !forwarder.IsEnabled() {
	//	err := errors.New(op).Msgf("forwarder is disabled: %s", qsoUpload.Service)
	//	if updateErr := s.DatabaseService.UpdateQsoUploadStatus(
	//		qsoUpload.ID,
	//		"disabled",
	//		qsoUpload.Attempts,
	//		"forwarder is disabled",
	//	); updateErr != nil {
	//		s.LoggerService.ErrorWith().Err(updateErr).Msg("Failed to update upload status to disabled")
	//	}
	//	return err
	//}
	//
	//// Forward the QSO
	//if err := forwarder.Forward(qsoUpload.Qso); err != nil {
	//	// Increment attempt count and record error
	//	if updateErr := s.DatabaseService.UpdateQsoUploadStatus(
	//		qsoUpload.ID,
	//		"pending",
	//		qsoUpload.Attempts+1,
	//		err.Error(),
	//	); updateErr != nil {
	//		s.LoggerService.ErrorWith().Err(updateErr).Msg("Failed to update upload status after forward failure")
	//	}
	//	return errors.New(op).Err(err).Msgf("failed to forward QSO to %s", qsoUpload.Service)
	//}
	//
	//// Mark as completed
	//if err := s.DatabaseService.UpdateQsoUploadStatus(
	//	qsoUpload.ID,
	//	"completed",
	//	qsoUpload.Attempts+1,
	//	"",
	//); err != nil {
	//	s.LoggerService.ErrorWith().Err(err).Msg("Failed to mark upload as completed")
	//	return errors.New(op).Err(err).Msg("failed to mark upload as completed")
	//}

	return nil
}

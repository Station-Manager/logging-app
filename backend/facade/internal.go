package facade

import (
	"strconv"
	"strings"
	"time"

	"github.com/Station-Manager/enums/upload/action"
	"github.com/Station-Manager/enums/upload/status"
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
	// Initialize the internal forwarding service workers
	s.forwarding = &forwarding{
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
		return errors.New(op).Msg("container is nil. Call SetContainer before using the facade.")
	}

	provider, ok := s.forwarders[qsoUpload.Service]
	if !ok {
		return errors.New(op).Msgf("no forwarder found for service: %s", qsoUpload.Service)
	}

	errState := ""
	uploadStatus := status.Failed
	err := provider.Forward(qsoUpload.Qso, qsoUpload.Action)
	if err != nil {
		s.LoggerService.ErrorWith().Err(err).Msgf("failed to forward QSO to %s", qsoUpload.Service)
		qsoUpload.Attempts++
		errState = errors.Root(err).Error()
	} else {
		qsoUpload.Attempts = 0
		qsoUpload.LastError = ""
		uploadStatus = status.Uploaded
	}

	act, _ := action.Parse(qsoUpload.Action)
	uerr := s.DatabaseService.UpdateQsoUploadStatus(qsoUpload.ID, uploadStatus, act, qsoUpload.Attempts, errState)
	if uerr != nil {
		s.LoggerService.ErrorWith().Int64("qso_id", qsoUpload.QsoID).Str("service", qsoUpload.Service).Err(uerr).Msg("Database error: Failed to update upload status after forward failure")
	}

	if err == nil {
		s.LoggerService.InfoWith().Int64("qso_id", qsoUpload.QsoID).Str("service", qsoUpload.Service).Msg("QSO uploaded successfully")
	} else {
		s.LoggerService.WarnWith().Int64("qso_id", qsoUpload.QsoID).Str("service", qsoUpload.Service).Err(err).Msg("Failed to upload QSO")
	}

	return nil
}

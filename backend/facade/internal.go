package facade

import (
	"database/sql"
	stderr "errors"
	"github.com/Station-Manager/errors"
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

// openAndLoadFromDatabase initializes the database connection, applies migrations, loads the default
// logbook, and generates a session ID.
func (s *Service) openAndLoadFromDatabase() error {
	const op errors.Op = "facade.Service.loadFromDatabase"

	// Open and migrate the database. Don't need to ping as opening the database will do that.
	if err := s.DatabaseService.Open(); err != nil {
		err = errors.New(op).Err(err)
		s.LoggerService.ErrorWith().Err(err).Msg("Failed to open database.")
		return err
	}
	if err := s.DatabaseService.Migrate(); err != nil {
		err = errors.New(op).Err(err)
		s.LoggerService.ErrorWith().Err(err).Msg("Failed to migrate database.")
		return err
	}

	// Load the default logbook
	logbook, err := s.DatabaseService.FetchLogbookByID(s.requiredCfgs.DefaultRigID)
	if err != nil {
		err = errors.New(op).Err(err)
		s.LoggerService.ErrorWith().Err(err).Msg("Failed to fetch logbook.")
		return err
	}
	s.CurrentLogbook = logbook

	// Generate a new session id
	s.sessionID, err = s.DatabaseService.GenerateNewSessionID()
	if err != nil {
		err = errors.New(op).Err(err)
		s.LoggerService.ErrorWith().Err(err).Msg("Failed to generate new session ID.")
		return err
	}

	return nil
}

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

func (s *Service) lookupCallsignOnline(callsign string) (types.ContactedStation, error) {

	return types.ContactedStation{}, nil
}

// contactedStationExistsByCallsign checks if a contacted station exists in the database using the given callsign.
// The callsign is trimmed, uppercased, and validated before querying.
// Returns a boolean indicating existence and an error if any operation fails.
func (s *Service) contactedStationExistsByCallsign(callsign string) (bool, error) {
	const op errors.Op = "facade.Service.contactedStationExistsByCallsign"
	if !s.initialized.Load() {
		err := errors.New(op).Msg(errMsgServiceNotInit)
		s.LoggerService.ErrorWith().Err(err).Msg(errMsgServiceNotInit)
		return false, errors.Root(err)
	}

	if !s.started.Load() {
		err := errors.New(op).Msg(errMsgServiceNotStarted)
		s.LoggerService.ErrorWith().Err(err).Msg(errMsgServiceNotStarted)
		return false, errors.Root(err)
	}

	callsign = strings.ToUpper(strings.TrimSpace(callsign))

	if len(callsign) < 3 {
		return false, errors.New(op).Msg(errMsgInvalidCallsign)
	}

	parsedCallsign := s.parseCallsign(callsign)

	exists, err := s.DatabaseService.ContactedStationExistsByCallsign(parsedCallsign)
	if err != nil {
		return false, errors.New(op).Err(err)
	}

	return exists, nil
}

func (s *Service) insertOrUpdateContactedStation(station types.ContactedStation) error {
	const op errors.Op = "facade.Service.insertOrUpdateContactedStation"

	model, err := s.DatabaseService.FetchContactedStationByCallsign(station.Call)
	if err != nil && !stderr.Is(err, sql.ErrNoRows) {
		return errors.New(op).Err(err)
	}

	// Error == sql.ErrNoRows
	if err != nil {
		s.LoggerService.DebugWith().Str("callsign", station.Call).Msg("Contacted station does not exist in database, inserting.")
		if _, err = s.DatabaseService.InsertContactedStation(station); err != nil {
			s.LoggerService.ErrorWith().Err(err).Msg("Failed to insert contacted station into database.")
			return errors.Root(err)
		}
		return nil
	}

	// Do this before the comparison to ensure the ID is set and the comparison doesn't fail because of it.
	station.ID = model.ID

	if model != station {
		s.LoggerService.DebugWith().Str("callsign", station.Call).Msg("Contacted station exists in database, but needs updating.")
		if err = s.DatabaseService.UpdateContactedStation(station); err != nil {
			s.LoggerService.ErrorWith().Err(err).Msg("Failed to update contacted station in database.")
			return errors.Root(err)
		}
	}

	return nil
}

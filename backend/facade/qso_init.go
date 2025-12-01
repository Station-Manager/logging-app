package facade

import (
	"database/sql"
	stderr "errors"
	"fmt"
	"github.com/Station-Manager/errors"
	"github.com/Station-Manager/types"
)

func (s *Service) initializeQso(callsign string) (*types.Qso, error) {
	const op errors.Op = "facade.Service.initializeQso"

	loggingStation, err := s.initLoggingStationSection()
	if err != nil {
		s.LoggerService.ErrorWith().Err(err).Msg("Failed to initialize the QSO's logging station section")
		return nil, errors.New(op).Err(err)
	}

	contactedStation, err := s.initContactedStationSection(callsign)
	if err != nil {
		s.LoggerService.ErrorWith().Err(err).Msg("Failed to initialize the QSO's contacted station section")
		return nil, errors.New(op).Err(err)
	}

	country, err := s.initCountrySection(callsign)
	if err != nil {
		s.LoggerService.ErrorWith().Err(err).Msg("Failed to initialize the QSO's country section")
		return nil, errors.New(op).Err(err)
	}

	qso := &types.Qso{
		LoggingStation:   loggingStation,
		ContactedStation: *contactedStation,
		CountryDetails:   country,
	}

	return qso, nil
}

// initLoggingStationSection initializes the logging station using the current logbook's callsign and configuration data.
// Returns a LoggingStation instance and an error if the station configuration retrieval fails.
func (s *Service) initLoggingStationSection() (types.LoggingStation, error) {
	const op errors.Op = "facade.Service.initLoggingStationSection"

	loggingStation, err := s.ConfigService.LoggingStationConfigs()
	if err != nil {
		return types.LoggingStation{}, errors.New(op).Err(err)
	}

	// All QSO must be logged to the current logbook, this is done by using the logbook's callsign
	// as the station callsign.
	loggingStation.StationCallsign = s.CurrentLogbook.Callsign

	return loggingStation, nil
}

func (s *Service) initContactedStationSection(callsign string) (*types.ContactedStation, error) {
	const op errors.Op = "facade.Service.initContactedStation"
	parsedCallsign := s.parseCallsign(callsign)

	contactedStation, err := s.DatabaseService.FetchContactedStationByCallsign(parsedCallsign)
	if err != nil && !stderr.Is(err, sql.ErrNoRows) {
		// There is something seriously wrong with the database, so we can't continue.
		s.LoggerService.ErrorWith().Err(err).Msgf("Failed to fetch contacted station with callsign %s", parsedCallsign)
		return nil, errors.New(op).Err(err)
	}

	// New contact, so we must look up the details from somewhere else: online services.
	if stderr.Is(err, sql.ErrNoRows) {
		contactedStation, err = s.lookupCallsignOnline(parsedCallsign)
		if err != nil {
			// This is not a show-stopper, so we can continue without the contacted station details.
			s.LoggerService.ErrorWith().Err(err).Msgf("Failed to look up contacted station with callsign %s", parsedCallsign)
		}
	}

	// We use the callsign provided by the caller as this might contain more information than the callsign returned by
	// the lookup (for example, portable/mobile/etc. suffixes).
	contactedStation.Call = callsign

	return &contactedStation, nil
}

func (s *Service) initCountrySection(callsign string) (types.Country, error) {
	const op errors.Op = "facade.Service.initCountrySection"

	parsedCallsign := s.parseCallsign(callsign)

	obj, err := s.HamnutService.Lookup(parsedCallsign)
	if err != nil {
		return types.Country{}, errors.New(op).Err(err)
	}

	fmt.Println(obj)

	return types.Country{}, nil
}

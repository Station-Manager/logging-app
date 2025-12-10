package facade

import (
	"github.com/Station-Manager/enums/cmds"
	"github.com/Station-Manager/errors"
	"github.com/Station-Manager/types"
	"github.com/wailsapp/wails/v2/pkg/runtime"
	"net/url"
	"strings"
)

// FetchUiConfig retrieves the UI configuration object. It returns an error if the service is not initialized, or the underlying
// ConfigService returns an error.
func (s *Service) FetchUiConfig() (*types.UiConfig, error) {
	const op errors.Op = "facade.Service.UiConfig"

	if !s.initialized.Load() {
		err := errors.New(op).Msg(errMsgServiceNotInit)
		s.LoggerService.ErrorWith().Err(err).Msg(errMsgServiceNotInit)
		return nil, err
	}

	requiredCfg, err := s.ConfigService.RequiredConfigs()
	if err != nil {
		err = errors.New(op).Err(err)
		s.LoggerService.ErrorWith().Err(err).Msg("Failed to fetch required configs.")
		return nil, err
	}

	return &types.UiConfig{
		DefaultRigID:       requiredCfg.DefaultRigID,
		Logbook:            s.CurrentLogbook,
		RigName:            s.CatService.RigConfig().Name,
		DefaultIsRandomQso: requiredCfg.DefaultIsRandomQso,
		DefaultTxPower:     requiredCfg.DefaultTxPower,
		UsePowerMultiplier: requiredCfg.UsePowerMultiplier,
		PowerMultiplier:    requiredCfg.PowerMultiplier,
		DefaultFreq:        requiredCfg.DefaultFreq,
		DefaultMode:        requiredCfg.DefaultMode,
	}, nil
}

// FetchCatStateValues retrieves the cat state values from the ConfigService.
func (s *Service) FetchCatStateValues() (map[string]map[string]string, error) {
	const op errors.Op = "facade.Service.FetchCatStateValues"

	if !s.initialized.Load() {
		err := errors.New(op).Msg(errMsgServiceNotInit)
		s.LoggerService.ErrorWith().Err(err).Msg(errMsgServiceNotInit)
		return nil, err
	}

	values, err := s.ConfigService.CatStateValues()
	if err != nil {
		err = errors.New(op).Err(err)
		s.LoggerService.ErrorWith().Err(err).Msg("Failed to fetch CAT state values.")
		return nil, err
	}

	return values, nil
}

// Ready checks if the service is initialized and started, then enqueues initialization and read commands to the CatService.
func (s *Service) Ready() error {
	const op errors.Op = "facade.Service.Ready"
	if !s.initialized.Load() {
		err := errors.New(op).Msg(errMsgServiceNotInit)
		s.LoggerService.ErrorWith().Err(err).Msg(errMsgServiceNotInit)
		return errors.Root(err)
	}

	if !s.started.Load() {
		err := errors.New(op).Msg(errMsgServiceNotStarted)
		s.LoggerService.ErrorWith().Err(err).Msg(errMsgServiceNotStarted)
		return errors.Root(err)
	}

	if err := s.CatService.EnqueueCommand(cmds.Init); err != nil {
		s.LoggerService.ErrorWith().Err(err).Msgf("Failed to enqueue command: %s", cmds.Init)
		return errors.New(op).Err(err)
	}

	if err := s.CatService.EnqueueCommand(cmds.Read); err != nil {
		s.LoggerService.ErrorWith().Err(err).Msgf("Failed to enqueue command: %s", cmds.Read)
		return errors.New(op).Err(err)
	}

	return nil
}

// NewQso initializes a new QSO object with the given callsign.
func (s *Service) NewQso(callsign string) (*types.Qso, error) {
	const op errors.Op = "facade.Service.NewQso"
	if !s.initialized.Load() {
		err := errors.New(op).Msg(errMsgServiceNotInit)
		s.LoggerService.ErrorWith().Err(err).Msg(errMsgServiceNotInit)
		return nil, errors.Root(err)
	}

	if !s.started.Load() {
		err := errors.New(op).Msg(errMsgServiceNotStarted)
		s.LoggerService.ErrorWith().Err(err).Msg(errMsgServiceNotStarted)
		return nil, errors.Root(err)
	}

	callsign = strings.ToUpper(strings.TrimSpace(callsign))

	if len(callsign) < 3 {
		return nil, errors.New(op).Msg(errMsgInvalidCallsign)
	}

	qso, err := s.initializeQso(callsign)
	if err != nil {
		return nil, errors.Root(err)
	}

	return qso, nil
}

// LogQso inserts a new QSO into the database.
func (s *Service) LogQso(qso types.Qso) error {
	const op errors.Op = "facade.Service.LogQso"
	if !s.initialized.Load() {
		err := errors.New(op).Msg(errMsgServiceNotInit)
		s.LoggerService.ErrorWith().Err(err).Msg(errMsgServiceNotInit)
		return errors.Root(err)
	}

	if !s.started.Load() {
		err := errors.New(op).Msg(errMsgServiceNotStarted)
		s.LoggerService.ErrorWith().Err(err).Msg(errMsgServiceNotStarted)
		return errors.Root(err)
	}

	// Set the current session ID
	qso.SessionID = s.sessionID

	// Insert the QSO into the database
	_, err := s.DatabaseService.InsertQso(qso)
	if err != nil {
		err = errors.New(op).Err(err)
		s.LoggerService.ErrorWith().Err(err).Msg("Failed to insert QSO into database.")
		return errors.Root(err)
	}
	s.LoggerService.InfoWith().Str("callsign", qso.Call).Msg("QSO logged successfully")

	// Check if the contacted station exists in the database and insert or update it if it does not
	// match the current QSO's contacted station. The ContactedStation object is loaded when
	// the QSO is initialized.
	if err = s.insertOrUpdateContactedStation(qso.ContactedStation); err != nil {
		// This is a serious error, but not fatal, so log and carry on.
		s.LoggerService.ErrorWith().Err(err).Msg("Failed to insert or update contacted station.")
	}

	if err = s.insertOrUpdateCountry(qso.CountryDetails); err != nil {
		// This is a serious error, but not fatal, so log and carry on.
		s.LoggerService.ErrorWith().Err(err).Msg("Failed to insert or update country.")
	}

	return nil
}

func (s *Service) CurrentSessionQsoSlice() ([]types.Qso, error) {
	const op errors.Op = "facade.Service.CurrentSessionQsoSlice"
	if !s.initialized.Load() {
		err := errors.New(op).Msg(errMsgServiceNotInit)
		s.LoggerService.ErrorWith().Err(err).Msg(errMsgServiceNotInit)
		return nil, err
	}

	if !s.started.Load() {
		err := errors.New(op).Msg(errMsgServiceNotStarted)
		s.LoggerService.ErrorWith().Err(err).Msg(errMsgServiceNotStarted)
		return nil, errors.Root(err)
	}

	list, err := s.DatabaseService.FetchQsoSliceBySessionID(s.sessionID)
	if err != nil {
		err = errors.New(op).Err(err)
		s.LoggerService.ErrorWith().Err(err).Msg("Failed to fetch QSOs by session ID.")
		return nil, errors.Root(err)
	}

	return list, nil
}

func (s *Service) OpenInBrowser(urlStr string) error {
	const op errors.Op = "facade.Service.OpenInBrowser"
	if !s.initialized.Load() {
		err := errors.New(op).Msg(errMsgServiceNotInit)
		s.LoggerService.ErrorWith().Err(err).Msg(errMsgServiceNotInit)
		return err
	}

	if !s.started.Load() {
		err := errors.New(op).Msg(errMsgServiceNotStarted)
		s.LoggerService.ErrorWith().Err(err).Msg(errMsgServiceNotStarted)
		return errors.Root(err)
	}

	if s.ctx == nil || (s.ctx.Err() != nil) {
		s.LoggerService.ErrorWith().Msg("Context is not set")
		return errors.New(op).Msg("Context is not set")
	}
	u, err := url.ParseRequestURI(urlStr)
	if err != nil || u.Scheme != "https" {
		err = errors.New(op).Err(err).Msg("Invalid or unsafe URL")
		s.LoggerService.ErrorWith().Err(err).Str("url_scheme", u.Scheme).Msg("Invalid or unsafe URL")
		return err
	}

	runtime.BrowserOpenURL(s.ctx, urlStr)

	return nil
}

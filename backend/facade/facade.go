package facade

import (
	"github.com/Station-Manager/errors"
	"github.com/Station-Manager/types"
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
		DefaultRigID: requiredCfg.DefaultRigID,
		Logbook:      s.CurrentLogbook,
		RigName:      s.CatService.RigConfig().Name,
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

func (s *Service) NewQso(callsign string) (*types.Qso, error) {
	const op errors.Op = "facade.Service.NewQso"
	if !s.initialized.Load() {
		err := errors.New(op).Msg(errMsgServiceNotInit)
		s.LoggerService.ErrorWith().Err(err).Msg(errMsgServiceNotInit)
		return nil, err
	}

	return nil, nil
}

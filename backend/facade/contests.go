package facade

import "github.com/Station-Manager/errors"

func (s *Service) IsContestDuplicate(callsign, band string) (bool, error) {
	const op errors.Op = "facade.Service.IsContestDuplicate"
	if !s.initialized.Load() {
		err := errors.New(op).Msg(errMsgServiceNotInit)
		s.LoggerService.ErrorWith().Err(err).Msg(errMsgServiceNotInit)
		return false, err
	}

	if !s.started.Load() {
		err := errors.New(op).Msg(errMsgServiceNotStarted)
		s.LoggerService.ErrorWith().Err(err).Msg(errMsgServiceNotStarted)
		return false, err
	}

	return false, nil
}

func (s *Service) TotalQsosByLogbookId(logbookId int64) (int64, error) {
	const op errors.Op = "facade.Service.TotalQsosByLogbookId"
	if !s.initialized.Load() {
		err := errors.New(op).Msg(errMsgServiceNotInit)
		s.LoggerService.ErrorWith().Err(err).Msg(errMsgServiceNotInit)
		return 0, err
	}

	if !s.started.Load() {
		err := errors.New(op).Msg(errMsgServiceNotStarted)
		s.LoggerService.ErrorWith().Err(err).Msg(errMsgServiceNotStarted)
		return 0, errors.Root(err)
	}

	if logbookId < 1 {
		err := errors.New(op).Msg("Invalid logbook id")
		s.LoggerService.ErrorWith().Err(err).Msg("Invalid logbook id")
		return 0, errors.Root(err)
	}

	count, err := s.DatabaseService.QsoCountByLogbookId(logbookId)
	if err != nil {
		err = errors.New(op).Err(err)
		s.LoggerService.ErrorWith().Err(err).Msg("Failed to get QSO count by logbook ID")
		return 0, errors.Root(err)
	}

	return count, nil
}

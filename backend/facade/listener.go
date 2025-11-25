package facade

import "github.com/Station-Manager/errors"

func (s *Service) catStatusChannelListener(shutdown <-chan struct{}) {
	const op errors.Op = "facade.Service.catStatusChannelListener"

	statusChannel, err := s.CatService.StatusChannel()
	if err != nil {
		err = errors.New(op).Err(err).Msg("Failed to get cat status channel.")
		s.LoggerService.ErrorWith().Err(err).Msg("Failed to start CAT status channel listener.")
		return
	}

	for {
		select {
		case <-shutdown:
			return
		case <-s.ctx.Done():
			return
		case status := <-statusChannel:
			s.LoggerService.DebugWith().Int("status", len(status)).Msg("Received CAT status update.")
		}
	}
}

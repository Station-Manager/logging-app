package facade

import (
	"github.com/Station-Manager/enums/events"
	"github.com/Station-Manager/errors"
	"github.com/wailsapp/wails/v2/pkg/runtime"
)

// catStatusChannelListener listens to the CAT status updates channel and logs received updates or handles shutdown signals.
func (s *Service) catStatusChannelListener(shutdown <-chan struct{}) {
	const op errors.Op = "facade.Service.catStatusChannelListener"

	statusChannel, err := s.CatService.StatusChannel()
	if err != nil {
		err = errors.New(op).Err(err).Msg("Failed to get cat status channel.")
		s.LoggerService.ErrorWith().Err(err).Msg("Failed to start CAT status channel listener.")
		return
	}

	if statusChannel == nil {
		err = errors.New(op).Msg("CAT status channel is nil")
		s.LoggerService.ErrorWith().Err(err).Msg("Cannot start listener with nil channel")
		return
	}

	for {
		select {
		case <-shutdown:
			return
		case <-s.ctx.Done():
			return
		case status, ok := <-statusChannel:
			if !ok {
				s.LoggerService.InfoWith().Msg("CAT status channel closed, listener exiting")
				return
			}
			// Emit the status update to the frontend
			runtime.EventsEmit(s.ctx, events.Status.String(), status)
		}
	}
}

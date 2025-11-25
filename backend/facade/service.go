package facade

import (
	"github.com/Station-Manager/config"
	"github.com/Station-Manager/errors"
	"github.com/Station-Manager/logging"
	"sync"
	"sync/atomic"
)

type Service struct {
	ConfigService *config.Service  `di.inject:"configservice"`
	LoggerService *logging.Service `di.inject:"loggingservice"`

	initialized atomic.Bool
	started     atomic.Bool // guarded via atomic operations; Start/Stop also hold mu for broader state

	initOnce sync.Once
	mu       sync.Mutex
}

func (s *Service) Initialize() error {
	const op errors.Op = "facade.Service.Initialize"

	var initErr error
	s.initOnce.Do(func() {
		if s.ConfigService == nil {
			initErr = errors.New(op).Msg(errMsgNilConfigService)
			return
		}

		if s.LoggerService == nil {
			initErr = errors.New(op).Msg(errMsgNilLoggerService)
			return
		}

		s.initialized.Store(true)
	})

	return initErr
}

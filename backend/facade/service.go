package facade

import (
	"context"
	"github.com/Station-Manager/config"
	"github.com/Station-Manager/database"
	"github.com/Station-Manager/errors"
	"github.com/Station-Manager/iocdi"
	"github.com/Station-Manager/logging"
	"sync"
	"sync/atomic"
)

const (
	ServiceName = "logging-app-facade"
)

type runState struct {
	shutdownChannel chan struct{}
	wg              sync.WaitGroup
}

type Service struct {
	ConfigService   *config.Service   `di.inject:"configservice"`
	LoggerService   *logging.Service  `di.inject:"loggingservice"`
	DatabaseService *database.Service `di.inject:"databaseservice"`

	container *iocdi.Container
	ctx       context.Context

	currentRun *runState

	initialized atomic.Bool
	started     atomic.Bool // guarded via atomic operations; Start/Stop also hold mu for a broader state

	initOnce sync.Once
	mu       sync.Mutex
}

// Initialize sets up the Service instance by verifying required dependencies and initializing its state.
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

		if s.DatabaseService == nil {
			initErr = errors.New(op).Msg(errMsgNilDatabaseService)
			return
		}

		s.initialized.Store(true)
	})

	return initErr
}

func (s *Service) SetContainer(container *iocdi.Container) error {
	const op errors.Op = "facade.Service.SetContainer"
	if !s.initialized.Load() {
		return errors.New(op).Msg(errMsgServiceNotInit)
	}

	if s.started.Load() {
		return nil
	}

	if container == nil {
		return errors.New(op).Msg("Container cannot be nil")
	}

	s.container = container

	return nil
}

func (s *Service) Start() error {
	const op errors.Op = "facade.Service.Start"
	if !s.initialized.Load() {
		return errors.New(op).Msg(errMsgServiceNotInit)
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	if s.started.Load() {
		return nil
	}

	if err := s.DatabaseService.Open(); err != nil {
		return errors.New(op).Err(err)
	}

	run := &runState{
		shutdownChannel: make(chan struct{}),
	}
	s.currentRun = run

	s.started.Store(true)

	return nil
}

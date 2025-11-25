package facade

import (
	"context"
	"github.com/Station-Manager/cat"
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
	CatService      *cat.Service      `di.inject:"catservice"`

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

		if s.CatService == nil {
			initErr = errors.New(op).Msg(errMsgNilCatService)
			return
		}

		s.initialized.Store(true)
	})

	return initErr
}

// SetContainer sets the IOC container for the Service. Returns an error if the Service is uninitialized or the container is nil.
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

// Start begins the Service lifecycle by initializing dependencies, opening the database, and marking it as started.
func (s *Service) Start(ctx context.Context) error {
	const op errors.Op = "facade.Service.Start"
	if !s.initialized.Load() {
		return errors.New(op).Msg(errMsgServiceNotInit)
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	if s.started.Load() {
		return nil
	}

	if ctx == nil || ctx.Err() != nil {
		return errors.New(op).Msg("Context cannot be nil or cancelled")
	}
	s.ctx = ctx

	// Open and migrate the database. Don't need to ping as opening the database will do that.
	if err := s.DatabaseService.Open(); err != nil {
		return errors.New(op).Err(err)
	}
	if err := s.DatabaseService.Migrate(); err != nil {
		return errors.New(op).Err(err)
	}

	run := &runState{
		shutdownChannel: make(chan struct{}),
	}
	s.currentRun = run

	s.launchWorkerThread(run, s.catStatusChannelListener, "catStatusChannelListener")

	s.started.Store(true)

	return nil
}

// Stop gracefully shuts down the service, closes resources, and resets the service state. Returns an error if any failure occurs.
func (s *Service) Stop() error {
	const op errors.Op = "facade.Service.Stop"
	if !s.initialized.Load() {
		return errors.New(op).Msg(errMsgServiceNotInit)
	}

	run := s.currentRun
	if run != nil && run.shutdownChannel != nil {
		select {
		case <-run.shutdownChannel:
			// already closed; nothing to do
		default:
			close(run.shutdownChannel)
		}
	}

	if run != nil {
		run.wg.Wait()
	}

	if err := s.DatabaseService.Close(); err != nil {
		s.LoggerService.ErrorWith().Err(err).Msg("Failed to close database")
	}

	s.currentRun = nil
	s.started.Store(false)

	return nil
}

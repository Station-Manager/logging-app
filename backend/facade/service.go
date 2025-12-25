package facade

import (
	"context"
	"sync"
	"sync/atomic"
	"time"

	"github.com/Station-Manager/cat"
	"github.com/Station-Manager/config"
	"github.com/Station-Manager/database/sqlite"
	"github.com/Station-Manager/email"
	"github.com/Station-Manager/errors"
	fwdrs "github.com/Station-Manager/forwarding"
	"github.com/Station-Manager/iocdi"
	"github.com/Station-Manager/logging"
	"github.com/Station-Manager/lookup/hamnut"
	"github.com/Station-Manager/lookup/qrz"
	"github.com/Station-Manager/types"
	"github.com/go-playground/validator/v10"
)

const (
	ServiceName = "logging-app-facade"
)

type runState struct {
	shutdownChannel chan struct{}
	wg              sync.WaitGroup
}

type Service struct {
	ConfigService       *config.Service  `di.inject:"configservice"`
	LoggerService       *logging.Service `di.inject:"loggingservice"`
	DatabaseService     *sqlite.Service  `di.inject:"sqliteservice"`
	CatService          *cat.Service     `di.inject:"catservice"`
	HamnutLookupService *hamnut.Service  `di.inject:"hamnutlookupservice"`
	QrzLookupService    *qrz.Service     `di.inject:"qrzlookupservice"`
	EmailService        *email.Service   `di.inject:"emailservice"`

	forwarders map[string]fwdrs.Forwarder

	requiredCfgs   *types.RequiredConfigs
	CurrentLogbook types.Logbook
	sessionID      int64

	container *iocdi.Container
	ctx       context.Context

	currentRun *runState

	forwarding *forwarding

	initialized atomic.Bool
	started     atomic.Bool // guarded via atomic operations; Start/Stop also hold mu for a broader state

	initOnce sync.Once
	mu       sync.Mutex

	validate *validator.Validate
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

		if s.HamnutLookupService == nil {
			initErr = errors.New(op).Msg(errMsgNilHamnutService)
			return
		}

		if s.QrzLookupService == nil {
			initErr = errors.New(op).Msg(errMsgNilQrzService)
			return
		}

		if err := s.initializeValidation(); err != nil {
			initErr = errors.New(op).Err(err)
			return
		}

		if err := s.initializeForwarding(); err != nil {
			initErr = errors.New(op).Err(err)
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
		err := errors.New(op).Msg(errMsgServiceNotInit)
		s.LoggerService.ErrorWith().Err(err).Msg(errMsgServiceNotInit)
		return err
	}

	if s.started.Load() {
		return nil
	}

	if container == nil {
		err := errors.New(op).Msg("Container cannot be nil")
		s.LoggerService.ErrorWith().Err(err).Msg("Container cannot be nil")
		return err
	}

	s.container = container

	return nil
}

// Start begins the Service lifecycle by initializing dependencies, opening the database, and marking it as started.
func (s *Service) Start(ctx context.Context) error {
	const op errors.Op = "facade.Service.Start"
	if !s.initialized.Load() {
		err := errors.New(op).Msg(errMsgServiceNotInit)
		s.LoggerService.ErrorWith().Err(err).Msg(errMsgServiceNotInit)
		return errors.Root(err)
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	if s.started.Load() {
		return nil
	}

	if s.container == nil {
		return errors.New(op).Msg("Container is nil. Please call SetContainer() before calling Start()")
	}

	if ctx == nil || ctx.Err() != nil {
		err := errors.New(op).Msg("Context cannot be nil or cancelled")
		s.LoggerService.ErrorWith().Msg("Context cannot be nil or cancelled")
		return errors.Root(err)
	}
	s.ctx = ctx

	reqCfg, err := s.ConfigService.RequiredConfigs()
	if err != nil {
		err = errors.New(op).Err(err)
		s.LoggerService.ErrorWith().Err(err).Msg("Failed to fetch required configs.")
		return errors.Root(err)
	}
	s.requiredCfgs = &reqCfg

	// Start the database service
	if err = s.openAndLoadFromDatabase(); err != nil {
		err = errors.New(op).Err(err)
		s.LoggerService.ErrorWith().Err(err).Msg("Failed to open and load from database.")
		return errors.Root(err)
	}

	// Start the CAT service
	if err = s.CatService.Start(); err != nil {
		err = errors.New(op).Err(err)
		s.LoggerService.ErrorWith().Err(err).Msg("Failed to start CAT service.")
		return errors.Root(err)
	}

	run := &runState{
		shutdownChannel: make(chan struct{}),
	}
	s.currentRun = run

	s.launchWorkerThread(run, s.catStatusChannelListener, "catStatusChannelListener")

	// Create a map of all the configured forwarders
	cfgs, _ := s.ConfigService.ForwarderConfigs() // Error discarded as ForwarderConfigs err is always nil
	s.forwarders = make(map[string]fwdrs.Forwarder, len(cfgs))
	for _, cfg := range cfgs {
		name := cfg.Name
		obj, serr := s.container.ResolveSafe(name)
		if serr != nil {
			s.LoggerService.WarnWith().Err(serr).Str("name", name).Msg("Failed to resolve forwarder service")
			continue
		}
		fwd, ok := obj.(fwdrs.Forwarder)
		if !ok {
			return errors.New(op).Msg("Failed to cast Forwarder service")
		}

		s.forwarders[name] = fwd
	}

	// Update forwarder poll interval from config
	s.forwarding.pollInterval = s.requiredCfgs.QsoForwardingIntervalSeconds * time.Second

	// Start the forwarder
	if s.forwarding != nil {
		if err = s.forwarding.start(s.ctx, run.shutdownChannel); err != nil {
			err = errors.New(op).Err(err)
			s.LoggerService.ErrorWith().Err(err).Msg("Failed to start QSO forwarder.")
			return errors.Root(err)
		}
	}

	s.started.Store(true)

	return nil
}

// Stop gracefully shuts down the service, closes resources, and resets the service state. Returns an error if any failure occurs.
func (s *Service) Stop() error {
	const op errors.Op = "facade.Service.Stop"
	if !s.initialized.Load() {
		err := errors.New(op).Msg(errMsgServiceNotInit)
		s.LoggerService.ErrorWith().Err(err).Msg(errMsgServiceNotInit)
		return err
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

	if s.forwarding != nil {
		close(s.forwarding.forwardingQueue)
		s.forwarding.wg.Wait()
	}

	if run != nil {
		run.wg.Wait()
	}

	// Stop the CAT service
	if err := s.CatService.Stop(); err != nil {
		s.LoggerService.ErrorWith().Err(err).Msg("Failed to stop CAT service")
	}

	// Soft-delete the session ID
	if err := s.DatabaseService.SoftDeleteSessionByID(s.sessionID); err != nil {
		// Not a show-stopper, just log the error
		s.LoggerService.ErrorWith().Err(err).Msg("Failed to soft-delete session ID")
	}

	// Stop the database service
	if err := s.DatabaseService.Close(); err != nil {
		// Log the error, but it's not fatal
		s.LoggerService.ErrorWith().Err(err).Msg("Failed to close database")
	}

	s.currentRun = nil
	s.started.Store(false)

	return nil
}

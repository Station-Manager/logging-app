package main

import (
	"github.com/Station-Manager/config"
	"github.com/Station-Manager/database"
	"github.com/Station-Manager/errors"
	"github.com/Station-Manager/iocdi"
	"github.com/Station-Manager/logging"
	"reflect"
	"strings"
)

func initializeContainer(workingDir string) error {
	const op errors.Op = "logging-app.main.initializeContainer"

	container = iocdi.New()

	if err := container.RegisterInstance("workingdir", workingDir); err != nil {
		return errors.New(op).Err(err)
	}
	if err := container.Register(config.ServiceName, reflect.TypeOf((*config.Service)(nil))); err != nil {
		return errors.New(op).Err(err)
	}
	if err := container.Register(logging.ServiceName, reflect.TypeOf((*logging.Service)(nil))); err != nil {
		return errors.New(op).Err(err)
	}
	if err := container.Register(database.ServiceName, reflect.TypeOf((*database.Service)(nil))); err != nil {
		return errors.New(op).Err(err)
	}

	if err := container.Build(); err != nil {
		return errors.New(op).Err(err)
	}

	return nil
}

// isDevelopment determines if the current application version is a development version by checking if "dev" is in its name.
func isDevelopment() bool {
	return strings.Contains(strings.ToLower(version), "dev")
}

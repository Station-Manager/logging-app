package main

import (
	"github.com/Station-Manager/errors"
	"github.com/Station-Manager/iocdi"
	"strings"
)

func initializeContainer(workingDir string) error {
	const op errors.Op = "logging-app.main.initializeContainer"

	container = iocdi.New()

	if err := container.RegisterInstance("workingdir", workingDir); err != nil {
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

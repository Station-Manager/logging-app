package main

import (
	"embed"
	"fmt"
	"github.com/Station-Manager/iocdi"
	"github.com/Station-Manager/utils"
	"os"
)

const (
	minWidth  int = 1024
	minHeight int = 751
)

var (
	version string
)

var container *iocdi.Container

//go:embed all:frontend/build
var assets embed.FS

func main() {
	workingDir, err := utils.WorkingDir()
	if err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "failed to determine working directory: %v\n", err)
		os.Exit(ExitWorkingDir)
	}

	if err = initializeContainer(workingDir); err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "failed to initialize container: %v\n", err)
		os.Exit(ExitContainerInit)
	}
}

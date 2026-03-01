// logging-app
package main

import (
	"embed"
	"fmt"
	"os"
	"runtime/debug"

	"github.com/Station-Manager/errors"
	"github.com/Station-Manager/iocdi"
	"github.com/Station-Manager/utils"
	"github.com/wailsapp/wails/v2"
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
	defer func() {
		if r := recover(); r != nil {
			_, _ = fmt.Fprintf(os.Stderr, "PANIC in main: %v\n", r)
			_, _ = fmt.Fprintf(os.Stderr, "Stack trace:\n%s\n", debug.Stack())
			os.Exit(ExitPanic)
		}
	}()

	workingDir, err := utils.WorkingDir()
	if err != nil {
		errors.PrintChain(err)
		_, _ = fmt.Fprintf(os.Stderr, "failed to determine working directory: %v\n", errors.Root(err))
		os.Exit(ExitWorkingDir)
	}

	if err = initializeContainer(workingDir); err != nil {
		errors.PrintChain(err)
		_, _ = fmt.Fprintf(os.Stderr, "failed to initialize container: %v\n", errors.Root(err))
		os.Exit(ExitContainerInit)
	}

	facade, err := getFacadeService()
	if err != nil {
		errors.PrintChain(err)
		_, _ = fmt.Fprintf(os.Stderr, "failed to get facade service: %v\n", errors.Root(err))
		os.Exit(ExitFacadeService)
	}

	if err = facade.SetContainer(container); err != nil {
		errors.PrintChain(err)
		_, _ = fmt.Fprintf(os.Stderr, "failed to set container: %v\n", errors.Root(err))
		os.Exit(ExitFacadeService)
	}

	if err = wails.Run(mainOpts(facade)); err != nil {
		errors.PrintChain(err)
		_, _ = fmt.Fprintf(os.Stderr, "failed to run wails: %v\n", errors.Root(err))
		os.Exit(ExitWailsRun)
	}
}

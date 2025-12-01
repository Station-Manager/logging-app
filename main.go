package main

import (
	"context"
	"embed"
	"fmt"
	"github.com/Station-Manager/enums/events"
	"github.com/Station-Manager/enums/tags"
	"github.com/Station-Manager/errors"
	"github.com/Station-Manager/iocdi"
	"github.com/Station-Manager/utils"
	"github.com/wailsapp/wails/v2"
	"github.com/wailsapp/wails/v2/pkg/logger"
	"github.com/wailsapp/wails/v2/pkg/options"
	"github.com/wailsapp/wails/v2/pkg/options/assetserver"
	"github.com/wailsapp/wails/v2/pkg/options/windows"
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
		_, _ = fmt.Fprintf(os.Stderr, "container initialization failed: %s\n", errors.Root(err))
		os.Exit(ExitContainerInit)
	}

	facade, err := getFacadeService()
	if err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "failed to get facade service: %v\n", err)
		os.Exit(ExitFacadeService)
	}

	startup := func(ctx context.Context) {
		if err = facade.Start(ctx); err != nil {
			panic(err)
		}
	}

	shutdown := func(ctx context.Context) {
		if err = facade.Stop(); err != nil {
			panic(err)
		}
	}

	opts := &options.App{
		Title:             fmt.Sprintf("%s: %s", AppTitle, version),
		Width:             minWidth,
		Height:            minHeight,
		DisableResize:     true,
		Fullscreen:        false,
		Frameless:         false,
		MinWidth:          minWidth,
		MinHeight:         minHeight,
		MaxWidth:          minWidth,
		MaxHeight:         minHeight,
		StartHidden:       false,
		HideWindowOnClose: false,
		AlwaysOnTop:       false,
		BackgroundColour:  &options.RGBA{R: 255, G: 255, B: 255, A: 255},
		AssetServer: &assetserver.Options{
			Assets: assets,
		},
		Menu:               nil,
		Logger:             nil,
		LogLevel:           logger.WARNING,
		LogLevelProduction: logger.ERROR,
		OnStartup:          startup,
		OnDomReady:         nil,
		OnShutdown:         shutdown,
		OnBeforeClose:      nil,
		Bind: []interface{}{
			facade,
		},
		EnumBind: []interface{}{
			tags.AllCatStateTags,
			events.AllEvents,
		},
		WindowStartState:                 options.Normal,
		ErrorFormatter:                   nil,
		CSSDragProperty:                  "",
		CSSDragValue:                     "",
		EnableDefaultContextMenu:         false,
		EnableFraudulentWebsiteDetection: false,
		SingleInstanceLock:               nil,
		Windows: &windows.Options{
			WebviewIsTransparent:                true,
			WindowIsTranslucent:                 false,
			DisableWindowIcon:                   false,
			IsZoomControlEnabled:                false,
			ZoomFactor:                          0,
			DisablePinchZoom:                    false,
			DisableFramelessWindowDecorations:   false,
			WebviewUserDataPath:                 "",
			WebviewBrowserPath:                  "",
			Theme:                               windows.SystemDefault,
			CustomTheme:                         nil,
			BackdropType:                        0,
			Messages:                            nil,
			ResizeDebounceMS:                    0,
			OnSuspend:                           nil,
			OnResume:                            nil,
			WebviewGpuIsDisabled:                false,
			WebviewDisableRendererCodeIntegrity: false,
			EnableSwipeGestures:                 false,
		},
		Mac:          nil,
		Linux:        nil,
		Experimental: nil,
		Debug: options.Debug{
			OpenInspectorOnStartup: isDevelopment(),
		},
		DragAndDrop: nil,
	}

	if err = wails.Run(opts); err != nil {
		panic(err)
	}
}

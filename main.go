package main

import (
	"embed"
	"fmt"
	"os"
	"runtime"
	"strings"
	"sync/atomic"

	"github.com/SUDOKU-ASCII/sudoku-desktop/internal/core"
	"github.com/SUDOKU-ASCII/sudoku/pkg/logx"
	"github.com/wailsapp/wails/v3/pkg/application"
	"github.com/wailsapp/wails/v3/pkg/events"
)

//go:embed all:frontend/dist
var assets embed.FS

//go:embed all:runtime/bin
var bundledRuntime embed.FS

//go:embed build/appicon.png
var trayIcon []byte

//go:embed build/tray/logo-monochrome.svg
var trayTemplateIcon []byte

func init() {
	application.RegisterEvent[core.RuntimeState](core.EventStateUpdated)
	application.RegisterEvent[core.LogEntry](core.EventLogAdded)
}

func main() {
	if runNetworkHelper() {
		return
	}
	startHidden := false
	for _, arg := range os.Args[1:] {
		switch strings.ToLower(strings.TrimSpace(arg)) {
		case "--autostart", "--background", "--hidden":
			startHidden = true
		}
	}

	appService := NewApp(bundledRuntime, "runtime/bin")
	var quitting atomic.Bool
	var app *application.App
	var mainWindow application.Window

	showExistingInstance := func() {
		application.InvokeAsync(func() {
			if app != nil {
				app.Show()
			}
			if mainWindow != nil {
				if mainWindow.IsMinimised() {
					mainWindow.UnMinimise()
				}
				mainWindow.Restore()
				mainWindow.Show().Focus()
			}
			if app != nil {
				dialog := app.Dialog.Info().
					SetTitle("Sudoku уже запущен").
					SetMessage("Переключились на открытое окно Sudoku.")
				if mainWindow != nil {
					dialog.AttachToWindow(mainWindow)
				}
				dialog.Show()
			}
		})
	}

	app = application.New(application.Options{
		Name:        "sudoku-app",
		Description: "Минималистичный VPN-клиент Sudoku",
		Icon:        trayIcon,
		Services: []application.Service{
			application.NewService(appService),
		},
		Assets: application.AssetOptions{
			Handler: application.BundledAssetFileServer(assets),
		},
		SingleInstance: &application.SingleInstanceOptions{
			UniqueID: "com.sanderstripa.sudoku.vpn",
			ExitCode: 0,
			OnSecondInstanceLaunch: func(_ application.SecondInstanceData) {
				showExistingInstance()
			},
		},
		Windows: application.WindowsOptions{
			DisableQuitOnLastWindowClosed: true,
		},
		Linux: application.LinuxOptions{
			DisableQuitOnLastWindowClosed: true,
			ProgramName:                   "sudoku-app",
		},
		Mac: application.MacOptions{
			ApplicationShouldTerminateAfterLastWindowClosed: false,
		},
		ShouldQuit: func() bool {
			quitting.Store(true)
			appService.ShutdownNow()
			return true
		},
	})

	requestQuit := func() {
		if !quitting.CompareAndSwap(false, true) {
			return
		}
		appService.ShutdownNow()
		app.Quit()
	}

	backgroundType, backgroundColour := mainWindowBackground()
	mainWindow = app.Window.NewWithOptions(application.WebviewWindowOptions{
		Name:                   "main",
		Title:                  "Sudoku",
		Width:                  410,
		Height:                 660,
		MinWidth:               360,
		MinHeight:              660,
		URL:                    "/",
		Frameless:              true,
		BackgroundType:         backgroundType,
		BackgroundColour:       backgroundColour,
		Hidden:                 startHidden,
		OpenInspectorOnStartup: false,
		Mac: application.MacWindow{
			Backdrop: application.MacBackdropTransparent,
		},
		Windows: application.WindowsWindow{
			DisableFramelessWindowDecorations: false,
		},
		Linux: application.LinuxWindow{
			Icon:                trayIcon,
			WindowIsTranslucent: true,
		},
		KeyBindings: map[string]func(window application.Window){
			"cmd+q": func(window application.Window) {
				requestQuit()
			},
			"ctrl+q": func(window application.Window) {
				requestQuit()
			},
			"cmd+w": func(window application.Window) {
				if runtime.GOOS == "darwin" {
					window.Hide()
				}
			},
			"ctrl+w": func(window application.Window) {
				if runtime.GOOS == "darwin" {
					window.Hide()
				}
			},
		},
	})

	mainWindow.RegisterHook(events.Common.WindowClosing, func(event *application.WindowEvent) {
		if quitting.Load() {
			return
		}
		event.Cancel()
		mainWindow.Hide()
	})
	configureApplicationMenu(app, mainWindow)

	trayWindow := app.Window.NewWithOptions(application.WebviewWindowOptions{
		Name: "tray-menu", Title: "Sudoku — меню", Width: 250, Height: 255,
		URL: "/?tray=1", Frameless: true, Hidden: true, AlwaysOnTop: true,
		ContentProtectionEnabled: false, BackgroundColour: backgroundColour,
		Windows: application.WindowsWindow{DisableFramelessWindowDecorations: false},
	})
	trayWindow.RegisterHook(events.Common.WindowClosing, func(event *application.WindowEvent) {
		event.Cancel()
		trayWindow.Hide()
	})
	// Keep the popup visible when screen clipping takes focus. Escape,
	// selecting an action, or clicking the tray icon again closes it.
	tray := app.SystemTray.New()
	tray.SetTooltip("Sudoku")
	if runtime.GOOS == "darwin" {
		tray.SetTemplateIcon(trayTemplateIcon)
	} else {
		tray.SetIcon(trayIcon)
	}
	showTray := func() {
		if trayWindow.IsVisible() {
			trayWindow.Hide()
			return
		}
		_ = tray.PositionWindow(trayWindow, 8)
		trayWindow.Show().Focus()
	}
	tray.OnClick(showTray)
	tray.OnRightClick(showTray)
	app.Event.On("tray:open", func(event *application.CustomEvent) {
		page, _ := event.Data.(string)
		if page != "settings" && page != "servers" {
			page = "home"
		}
		application.InvokeAsync(func() {
			trayWindow.Hide()
			mainWindow.Show().Focus()
			app.Event.Emit("main:navigate", page)
		})
	})
	app.Event.On("tray:quit", func(_ *application.CustomEvent) { go requestQuit() })
	err := app.Run()
	if err != nil {
		logx.Errorf("Desktop", "wails run failed: %v", err)
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

package tui

import (
	"capture-board-selector/internal/captureboard/app"
	"capture-board-selector/internal/captureboard/domain"

	"github.com/rivo/tview"
)

type App struct {
	tv             *tview.Application
	pages          *tview.Pages
	currentVersion string
	checkUpdate    app.CheckUpdateUseCase
	installUpdate  app.InstallUpdateUseCase
	checkFFmpeg    app.CheckFFmpegUseCase
	installFFmpeg  app.InstallFFmpegUseCase
	listDevices    app.ListDevicesUseCase
	runPreview     app.RunPreviewUseCase
	reporter       domain.ErrorReporter

	selectedVideo *domain.Device
	selectedAudio *domain.Device
}

func NewApp(
	version string,
	checkUpdate app.CheckUpdateUseCase,
	installUpdate app.InstallUpdateUseCase,
	check app.CheckFFmpegUseCase,
	install app.InstallFFmpegUseCase,
	list app.ListDevicesUseCase,
	run app.RunPreviewUseCase,
	reporter domain.ErrorReporter,
) *App {
	return &App{
		tv:             tview.NewApplication(),
		pages:          tview.NewPages(),
		currentVersion: version,
		checkUpdate:    checkUpdate,
		installUpdate:  installUpdate,
		checkFFmpeg:    check,
		installFFmpeg:  install,
		listDevices:    list,
		runPreview:     run,
		reporter:       reporter,
	}
}

func (a *App) Run() error {
	proceed := func() {
		if a.checkFFmpeg.Execute() {
			a.showDeviceSelection()
		} else {
			a.showInstallPrompt()
		}
	}

	a.checkForUpdate(proceed)

	if err := a.tv.SetRoot(a.pages, true).Run(); err != nil {
		return err
	}

	if a.selectedVideo != nil && a.selectedAudio != nil {
		return a.runPreview.Execute(*a.selectedVideo, *a.selectedAudio)
	}
	return nil
}

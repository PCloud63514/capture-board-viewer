package tui

import (
	"capture-board-selector/internal/captureboard/app"
	"capture-board-selector/internal/captureboard/domain"

	"github.com/rivo/tview"
)

type App struct {
	tv            *tview.Application
	pages         *tview.Pages
	checkFFmpeg   app.CheckFFmpegUseCase
	installFFmpeg app.InstallFFmpegUseCase
	listDevices   app.ListDevicesUseCase
	runPreview    app.RunPreviewUseCase

	selectedVideo *domain.Device
	selectedAudio *domain.Device
}

func NewApp(
	check app.CheckFFmpegUseCase,
	install app.InstallFFmpegUseCase,
	list app.ListDevicesUseCase,
	run app.RunPreviewUseCase,
) *App {
	return &App{
		tv:            tview.NewApplication(),
		pages:         tview.NewPages(),
		checkFFmpeg:   check,
		installFFmpeg: install,
		listDevices:   list,
		runPreview:    run,
	}
}

func (a *App) Run() error {
	if a.checkFFmpeg.Execute() {
		a.showDeviceSelection()
	} else {
		a.showInstallPrompt()
	}

	if err := a.tv.SetRoot(a.pages, true).Run(); err != nil {
		return err
	}

	if a.selectedVideo != nil && a.selectedAudio != nil {
		return a.runPreview.Execute(*a.selectedVideo, *a.selectedAudio)
	}
	return nil
}

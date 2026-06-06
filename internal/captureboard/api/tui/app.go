package tui

import (
	"fmt"

	"capture-board-selector/internal/captureboard/app"
	"capture-board-selector/internal/captureboard/domain"
	"capture-board-selector/internal/captureboard/infra/logger"

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
	logger.Start(a.currentVersion)
	defer logger.Close()

	logger.Log("시작", fmt.Sprintf("버전: %s", a.currentVersion))

	proceed := func() {
		logger.Log("FFmpeg", "설치 여부 확인 중")
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
		logger.Log("실행", fmt.Sprintf("ffplay 시작 — 비디오: %s, 오디오: %s", a.selectedVideo.Name, a.selectedAudio.Name))
		return a.runPreview.Execute(*a.selectedVideo, *a.selectedAudio)
	}
	return nil
}

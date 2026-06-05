package captureboard

import (
	"capture-board-selector/internal/captureboard/api/tui"
	"capture-board-selector/internal/captureboard/app"
	"capture-board-selector/internal/captureboard/domain"
	"capture-board-selector/internal/captureboard/infra/ffmpeg"
	infrasentry "capture-board-selector/internal/captureboard/infra/sentry"
)

// Infra
func NewFFmpegChecker() domain.FFmpegChecker      { return ffmpeg.NewChecker() }
func NewFFmpegInstaller() domain.FFmpegInstaller   { return ffmpeg.NewInstaller() }
func NewDeviceDiscoverer() domain.DeviceDiscoverer { return ffmpeg.NewDShowDiscoverer() }
func NewPreviewRunner() domain.PreviewRunner       { return ffmpeg.NewFFplayRunner() }
func NewErrorReporter() domain.ErrorReporter       { return infrasentry.NewReporter(SentryDSN) }

// App
func NewCheckFFmpegService(c domain.FFmpegChecker) app.CheckFFmpegUseCase {
	return app.NewCheckFFmpegService(c)
}
func NewInstallFFmpegService(i domain.FFmpegInstaller) app.InstallFFmpegUseCase {
	return app.NewInstallFFmpegService(i)
}
func NewListDevicesService(d domain.DeviceDiscoverer) app.ListDevicesUseCase {
	return app.NewListDevicesService(d)
}
func NewRunPreviewService(r domain.PreviewRunner) app.RunPreviewUseCase {
	return app.NewRunPreviewService(r)
}

// API

type TUIApp interface {
	Run() error
}

func NewTUIApp(
	check app.CheckFFmpegUseCase,
	install app.InstallFFmpegUseCase,
	list app.ListDevicesUseCase,
	run app.RunPreviewUseCase,
	reporter domain.ErrorReporter,
) TUIApp {
	return tui.NewApp(check, install, list, run, reporter)
}

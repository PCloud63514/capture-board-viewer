package captureboard

import (
	"capture-board-selector/internal/captureboard/api/tui"
	"capture-board-selector/internal/captureboard/app"
	"capture-board-selector/internal/captureboard/domain"
	"capture-board-selector/internal/captureboard/infra/ffmpeg"
	"capture-board-selector/internal/captureboard/infra/ghrelease"

	infrasentry "capture-board-selector/internal/captureboard/infra/sentry"
	"capture-board-selector/internal/captureboard/infra/updater"
)

// Infra
func NewFFmpegChecker() domain.FFmpegChecker      { return ffmpeg.NewChecker() }
func NewFFmpegInstaller() domain.FFmpegInstaller   { return ffmpeg.NewInstaller() }
func NewDeviceDiscoverer() domain.DeviceDiscoverer { return ffmpeg.NewDShowDiscoverer() }
func NewPreviewRunner() domain.PreviewRunner       { return ffmpeg.NewFFplayRunner() }
func NewErrorReporter() domain.ErrorReporter       { return infrasentry.NewReporter(SentryDSN) }
func NewUpdateChecker() domain.UpdateChecker       { return ghrelease.NewGitHubChecker() }
func NewUpdateInstaller() domain.UpdateInstaller   { return updater.NewInstaller() }

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
func NewCheckUpdateService(c domain.UpdateChecker) app.CheckUpdateUseCase {
	return app.NewCheckUpdateService(c)
}
func NewInstallUpdateService(i domain.UpdateInstaller) app.InstallUpdateUseCase {
	return app.NewInstallUpdateService(i)
}

// API

type TUIApp interface {
	Run() error
}

func NewTUIApp(
	checkUpdate app.CheckUpdateUseCase,
	installUpdate app.InstallUpdateUseCase,
	check app.CheckFFmpegUseCase,
	install app.InstallFFmpegUseCase,
	list app.ListDevicesUseCase,
	run app.RunPreviewUseCase,
	reporter domain.ErrorReporter,
) TUIApp {
	return tui.NewApp(Version, checkUpdate, installUpdate, check, install, list, run, reporter)
}

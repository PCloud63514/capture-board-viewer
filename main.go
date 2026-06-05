package main

import (
	captureboard "capture-board-selector/internal/captureboard"

	"go.uber.org/fx"
	"go.uber.org/fx/fxevent"
)

func main() {
	fx.New(
		fx.WithLogger(func() fxevent.Logger { return fxevent.NopLogger }),
		fx.Provide(
			captureboard.NewFFmpegChecker,
			captureboard.NewFFmpegInstaller,
			captureboard.NewDeviceDiscoverer,
			captureboard.NewPreviewRunner,
			captureboard.NewErrorReporter,
			captureboard.NewCheckFFmpegService,
			captureboard.NewInstallFFmpegService,
			captureboard.NewListDevicesService,
			captureboard.NewRunPreviewService,
			captureboard.NewTUIApp,
		),
		fx.Invoke(func(app captureboard.TUIApp, shut fx.Shutdowner) {
			go func() {
				app.Run()
				shut.Shutdown()
			}()
		}),
	).Run()
}

package domain

type FFmpegChecker interface {
	IsInstalled() bool
}

type FFmpegInstaller interface {
	Install(onProgress func(downloaded, total int64)) error
}

type DeviceDiscoverer interface {
	Discover() ([]Device, error)
}

type PreviewRunner interface {
	Run(video Device, audio Device) error
}

type ErrorReporter interface {
	Report(err error)
	Flush()
}

type UpdateInfo struct {
	Version     string
	DownloadURL string
}

type UpdateChecker interface {
	Check(currentVersion string) (*UpdateInfo, error)
}

type UpdateInstaller interface {
	Install(info UpdateInfo, onProgress func(downloaded, total int64)) error
}

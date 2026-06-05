package app

import "capture-board-selector/internal/captureboard/domain"

type InstallFFmpegUseCase interface {
	Execute(onProgress func(downloaded, total int64)) error
}

type InstallFFmpegService struct {
	installer domain.FFmpegInstaller
}

func NewInstallFFmpegService(installer domain.FFmpegInstaller) InstallFFmpegUseCase {
	return &InstallFFmpegService{installer: installer}
}

func (s *InstallFFmpegService) Execute(onProgress func(downloaded, total int64)) error {
	return s.installer.Install(onProgress)
}

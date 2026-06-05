package app

import "capture-board-selector/internal/captureboard/domain"

type InstallUpdateUseCase interface {
	Execute(info domain.UpdateInfo, onProgress func(downloaded, total int64)) error
}

type InstallUpdateService struct {
	installer domain.UpdateInstaller
}

func NewInstallUpdateService(installer domain.UpdateInstaller) InstallUpdateUseCase {
	return &InstallUpdateService{installer: installer}
}

func (s *InstallUpdateService) Execute(info domain.UpdateInfo, onProgress func(downloaded, total int64)) error {
	return s.installer.Install(info, onProgress)
}

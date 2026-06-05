package app

import "capture-board-selector/internal/captureboard/domain"

type RunPreviewUseCase interface {
	Execute(video domain.Device, audio domain.Device) error
}

type RunPreviewService struct {
	runner domain.PreviewRunner
}

func NewRunPreviewService(runner domain.PreviewRunner) RunPreviewUseCase {
	return &RunPreviewService{runner: runner}
}

func (s *RunPreviewService) Execute(video domain.Device, audio domain.Device) error {
	return s.runner.Run(video, audio)
}

package app

import "capture-board-selector/internal/captureboard/domain"

type CheckFFmpegUseCase interface {
	Execute() bool
}

type CheckFFmpegService struct {
	checker domain.FFmpegChecker
}

func NewCheckFFmpegService(checker domain.FFmpegChecker) CheckFFmpegUseCase {
	return &CheckFFmpegService{checker: checker}
}

func (s *CheckFFmpegService) Execute() bool {
	return s.checker.IsInstalled()
}

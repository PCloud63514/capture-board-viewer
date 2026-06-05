package app

import "capture-board-selector/internal/captureboard/domain"

type CheckUpdateUseCase interface {
	Execute(currentVersion string) (*domain.UpdateInfo, error)
}

type CheckUpdateService struct {
	checker domain.UpdateChecker
}

func NewCheckUpdateService(checker domain.UpdateChecker) CheckUpdateUseCase {
	return &CheckUpdateService{checker: checker}
}

func (s *CheckUpdateService) Execute(currentVersion string) (*domain.UpdateInfo, error) {
	return s.checker.Check(currentVersion)
}

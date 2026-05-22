package app

import (
	"context"
	"fmt"

	"capture-board-selector/internal/captureboard/domain"
)

type SelectorUseCase interface {
	LoadDevices(ctx context.Context) (domain.DeviceCatalog, error)
	RunPreview(ctx context.Context, selection domain.Selection) error
}

type SelectorService struct {
	discoverer domain.DeviceDiscoverer
	runner     domain.PreviewRunner
	checker    domain.DependencyChecker
}

func NewSelectorService(
	discoverer domain.DeviceDiscoverer,
	runner domain.PreviewRunner,
	checker domain.DependencyChecker,
) SelectorUseCase {
	return &SelectorService{
		discoverer: discoverer,
		runner:     runner,
		checker:    checker,
	}
}

func (s *SelectorService) LoadDevices(ctx context.Context) (domain.DeviceCatalog, error) {
	if err := s.checker.Ensure("ffmpeg", "ffplay"); err != nil {
		return domain.DeviceCatalog{}, err
	}

	devices, logFile, err := s.discoverer.Discover(ctx)
	if err != nil {
		return domain.DeviceCatalog{}, err
	}

	if !devices.HasSelectableDevices() {
		return domain.DeviceCatalog{}, fmt.Errorf("캡처보드 장치가 연결되어 있지 않습니다. 로그 확인: %s", logFile)
	}

	return devices, nil
}

func (s *SelectorService) RunPreview(ctx context.Context, selection domain.Selection) error {
	return s.runner.Run(ctx, selection)
}

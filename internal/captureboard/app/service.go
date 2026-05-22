package app

import (
	"context"

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

	devices, _, err := s.discoverer.Discover(ctx)
	if err != nil {
		return domain.DeviceCatalog{}, err
	}

	return devices, nil
}

func (s *SelectorService) RunPreview(ctx context.Context, selection domain.Selection) error {
	return s.runner.Run(ctx, selection)
}

package app

import "capture-board-selector/internal/captureboard/domain"

type ListDevicesUseCase interface {
	Execute() ([]domain.Device, error)
}

type ListDevicesService struct {
	discoverer domain.DeviceDiscoverer
}

func NewListDevicesService(discoverer domain.DeviceDiscoverer) ListDevicesUseCase {
	return &ListDevicesService{discoverer: discoverer}
}

func (s *ListDevicesService) Execute() ([]domain.Device, error) {
	return s.discoverer.Discover()
}

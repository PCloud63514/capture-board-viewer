package domain

import "context"

type DeviceDiscoverer interface {
	Discover(ctx context.Context) (DeviceCatalog, string, error)
}

type PreviewRunner interface {
	Run(ctx context.Context, selection Selection) error
}

type DependencyChecker interface {
	Ensure(names ...string) error
}

package captureboard

import (
	"context"

	"capture-board-selector/internal/captureboard/api/tui"
	"capture-board-selector/internal/captureboard/app"
	"capture-board-selector/internal/captureboard/domain"
	"capture-board-selector/internal/captureboard/infra/command"
	"capture-board-selector/internal/captureboard/infra/ffmpeg"
)

type Selector struct {
	useCase app.SelectorUseCase
}

func NewSelector() *Selector {
	return &Selector{
		useCase: app.NewSelectorService(
			ffmpeg.NewDiscoverer(),
			ffmpeg.NewPreviewRunner(),
			command.NewChecker(),
		),
	}
}

func (s *Selector) Select(ctx context.Context) (domain.Selection, error) {
	return tui.Run(ctx, s.useCase)
}

func (s *Selector) RunPreview(ctx context.Context, selection domain.Selection) error {
	return s.useCase.RunPreview(ctx, selection)
}

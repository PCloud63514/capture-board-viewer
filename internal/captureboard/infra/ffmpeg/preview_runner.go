package ffmpeg

import (
	"context"
	"fmt"
	"os"
	"os/exec"

	"capture-board-selector/internal/captureboard/domain"
)

type PreviewRunner struct{}

func NewPreviewRunner() *PreviewRunner {
	return &PreviewRunner{}
}

func (r *PreviewRunner) Run(ctx context.Context, selection domain.Selection) error {
	cmd := exec.CommandContext(ctx, "ffplay",
		"-hide_banner",
		"-loglevel", "error",
		"-f", "dshow",
		"-fflags", "nobuffer",
		"-flags", "low_delay",
		"-avioflags", "direct",
		"-rtbufsize", "512M",
		"-af", "aresample=resampler=soxr:osf=s32:async=1:min_comp=0.001:first_pts=0,adelay=0|0",
		"-use_wallclock_as_timestamps", "1",
		"-audio_buffer_size", "50",
		"-async", "1",
		"-i", fmt.Sprintf("video=%q:audio=%q", selection.Video, selection.Audio),
		"-video_size", "1920x1080",
		"-vf", "scale=1920:1080,setdar=16/9",
		"-window_title", "Preview",
		"-x", "1920",
		"-y", "1080",
	)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	return cmd.Run()
}

package ffmpeg

import (
	"fmt"
	"os"
	"os/exec"

	"capture-board-selector/internal/captureboard/domain"
)

type FFplayRunner struct{}

func NewFFplayRunner() domain.PreviewRunner {
	return &FFplayRunner{}
}

func (r *FFplayRunner) Run(video domain.Device, audio domain.Device) error {
	ffplay := resolveExe("ffplay", FFplayExe())
	cmd := exec.Command(ffplay,
		"-f", "dshow",
		"-rtbufsize", "100M",
		"-use_wallclock_as_timestamps", "1",
		"-audio_buffer_size", "30",
		"-i", fmt.Sprintf("video=%s:audio=%s", video.Name, audio.Name),
		"-video_size", "1920x1080",
		"-vf", "scale=1920:1080,setdar=16/9",
		"-window_title", "Preview",
		"-x", "1920", "-y", "1080",
		"-sync", "video",
		"-af", "adelay=0|0",
	)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

// resolveExe returns the system command if available, otherwise falls back to the APPDATA path.
func resolveExe(name, fallback string) string {
	if _, err := exec.LookPath(name); err == nil {
		return name
	}
	return fallback
}

package ffmpeg

import (
	"os"
	"os/exec"

	"capture-board-selector/internal/captureboard/domain"
)

type Checker struct{}

func NewChecker() domain.FFmpegChecker {
	return &Checker{}
}

func (c *Checker) IsInstalled() bool {
	if _, err := exec.LookPath("ffmpeg"); err == nil {
		return true
	}
	_, err := os.Stat(FFmpegExe())
	return err == nil
}

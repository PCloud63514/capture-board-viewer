package ffmpeg

import (
	"os"
	"os/exec"

	"capture-board-selector/internal/captureboard/domain"
	"capture-board-selector/internal/captureboard/infra/logger"
)

type Checker struct{}

func NewChecker() domain.FFmpegChecker {
	return &Checker{}
}

func (c *Checker) IsInstalled() bool {
	if _, err := exec.LookPath("ffmpeg"); err == nil {
		logger.Log("FFmpeg", "PATH에서 발견")
		return true
	}
	if _, err := os.Stat(FFmpegExe()); err == nil {
		logger.Log("FFmpeg", "AppData에서 발견: "+FFmpegExe())
		return true
	}
	logger.Log("FFmpeg", "설치되지 않음")
	return false
}

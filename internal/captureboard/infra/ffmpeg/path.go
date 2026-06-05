package ffmpeg

import (
	"os"
	"path/filepath"
)

func installDir() string {
	return filepath.Join(os.Getenv("APPDATA"), "capture-board-viewer", "ffmpeg")
}

func FFmpegExe() string {
	return filepath.Join(installDir(), "bin", "ffmpeg.exe")
}

func FFplayExe() string {
	return filepath.Join(installDir(), "bin", "ffplay.exe")
}

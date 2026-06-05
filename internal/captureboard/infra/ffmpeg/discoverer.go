package ffmpeg

import (
	"bytes"
	"os/exec"
	"strings"

	"capture-board-selector/internal/captureboard/domain"
)

type DShowDiscoverer struct{}

func NewDShowDiscoverer() domain.DeviceDiscoverer {
	return &DShowDiscoverer{}
}

func (d *DShowDiscoverer) Discover() ([]domain.Device, error) {
	ffmpeg := resolveExe("ffmpeg", FFmpegExe())
	cmd := exec.Command(ffmpeg, "-hide_banner", "-list_devices", "true", "-f", "dshow", "-i", "dummy")
	var stderr bytes.Buffer
	cmd.Stderr = &stderr
	_ = cmd.Run()

	var devices []domain.Device
	for _, line := range strings.Split(stderr.String(), "\n") {
		line = strings.TrimSpace(strings.ReplaceAll(line, "\uFEFF", ""))
		name := extractDeviceName(line)
		if name == "" {
			continue
		}
		if strings.Contains(line, "(video)") {
			devices = append(devices, domain.Device{Name: name, Type: domain.DeviceTypeVideo})
		} else if strings.Contains(line, "(audio)") {
			devices = append(devices, domain.Device{Name: name, Type: domain.DeviceTypeAudio})
		}
	}
	return devices, nil
}

func extractDeviceName(line string) string {
	start := strings.Index(line, "\"")
	end := strings.LastIndex(line, "\"")
	if start >= 0 && end > start {
		return line[start+1 : end]
	}
	return ""
}

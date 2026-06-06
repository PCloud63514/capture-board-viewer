package ffmpeg

import (
	"bytes"
	"fmt"
	"os/exec"
	"strings"
	"time"

	"capture-board-selector/internal/captureboard/domain"
	"capture-board-selector/internal/captureboard/infra/logger"
)

type DShowDiscoverer struct {
	version string
}

func NewDShowDiscoverer(version string) domain.DeviceDiscoverer {
	return &DShowDiscoverer{version: version}
}

func (d *DShowDiscoverer) Discover() ([]domain.Device, error) {
	ffmpeg := resolveExe("ffmpeg", FFmpegExe())
	cmd := exec.Command(ffmpeg, "-hide_banner", "-list_devices", "true", "-f", "dshow", "-i", "dummy")
	var stderr bytes.Buffer
	cmd.Stderr = &stderr
	_ = cmd.Run()

	output := stderr.String()
	d.saveLog(output)

	var devices []domain.Device
	for _, line := range strings.Split(output, "\n") {
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

func (d *DShowDiscoverer) saveLog(output string) {
	content := fmt.Sprintf("시각: %s\n버전: %s\n\n%s",
		time.Now().Format("2006-01-02 15:04:05"),
		d.version,
		output,
	)
	logger.Write("device", content)
}

func extractDeviceName(line string) string {
	start := strings.Index(line, "\"")
	end := strings.LastIndex(line, "\"")
	if start >= 0 && end > start {
		return line[start+1 : end]
	}
	return ""
}

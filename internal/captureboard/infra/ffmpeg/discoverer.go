package ffmpeg

import (
	"bytes"
	"context"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"capture-board-selector/internal/captureboard/domain"
)

type Discoverer struct{}

func NewDiscoverer() *Discoverer {
	return &Discoverer{}
}

func (d *Discoverer) Discover(ctx context.Context) (domain.DeviceCatalog, string, error) {
	timestamp := time.Now().Format("20060102_150405")
	if err := os.MkdirAll("logs", 0o755); err != nil {
		return domain.DeviceCatalog{}, "", err
	}

	logFile := filepath.Join("logs", fmt.Sprintf("device_output_%s.txt", timestamp))
	cmd := exec.CommandContext(ctx, "ffmpeg", "-hide_banner", "-list_devices", "true", "-f", "dshow", "-i", "dummy")

	var stderr bytes.Buffer
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil && stderr.Len() == 0 {
		return domain.DeviceCatalog{}, logFile, err
	}

	output := stderr.String()
	if err := os.WriteFile(logFile, []byte(output), 0o644); err != nil {
		return domain.DeviceCatalog{}, logFile, err
	}

	return parseDevices(output), logFile, nil
}

func parseDevices(output string) domain.DeviceCatalog {
	lines := strings.Split(output, "\n")
	catalog := domain.DeviceCatalog{}

	for _, line := range lines {
		line = strings.TrimSpace(strings.ReplaceAll(line, "\uFEFF", ""))
		switch {
		case strings.Contains(line, "(video)"):
			catalog.Videos = append(catalog.Videos, domain.Device{
				Name: extractDeviceName(line),
				Kind: domain.DeviceKindVideo,
			})
		case strings.Contains(line, "(audio)"):
			catalog.Audios = append(catalog.Audios, domain.Device{
				Name: extractDeviceName(line),
				Kind: domain.DeviceKindAudio,
			})
		}
	}

	return catalog
}

func extractDeviceName(line string) string {
	start := strings.Index(line, "\"")
	end := strings.LastIndex(line, "\"")
	if start >= 0 && end > start {
		return line[start+1 : end]
	}
	return line
}

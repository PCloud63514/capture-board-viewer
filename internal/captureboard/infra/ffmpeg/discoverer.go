package ffmpeg

import (
	"bytes"
	"fmt"
	"os/exec"
	"strings"

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
	logger.Log("장치", "ffmpeg dshow 장치 검색 시작")

	ffmpeg := resolveExe("ffmpeg", FFmpegExe())
	cmd := exec.Command(ffmpeg, "-hide_banner", "-list_devices", "true", "-f", "dshow", "-i", "dummy")
	var stderr bytes.Buffer
	cmd.Stderr = &stderr
	_ = cmd.Run()

	output := stderr.String()
	logger.Log("장치", "ffmpeg 원시 출력:\n"+output)

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

	d.logParsedDevices(devices)
	return devices, nil
}

func (d *DShowDiscoverer) logParsedDevices(devices []domain.Device) {
	var videos, audios []string
	for _, dev := range devices {
		if dev.Type == domain.DeviceTypeVideo {
			videos = append(videos, dev.Name)
		} else {
			audios = append(audios, dev.Name)
		}
	}

	logger.Log("장치", fmt.Sprintf("발견: 비디오 %d개, 오디오 %d개", len(videos), len(audios)))
	if len(videos) > 0 {
		logger.Log("장치", "비디오 목록: "+strings.Join(videos, ", "))
	}
	if len(audios) > 0 {
		logger.Log("장치", "오디오 목록: "+strings.Join(audios, ", "))
	}
}

func extractDeviceName(line string) string {
	start := strings.Index(line, "\"")
	end := strings.LastIndex(line, "\"")
	if start >= 0 && end > start {
		return line[start+1 : end]
	}
	return ""
}

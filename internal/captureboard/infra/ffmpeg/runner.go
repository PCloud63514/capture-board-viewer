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
		"-f", "dshow", // Windows DirectShow 입력
		"-rtbufsize", "100M", // 실시간 입력 버퍼 크기 (끊김 방지)
		"-use_wallclock_as_timestamps", "1", // 시스템 시계 기준으로 타임스탬프 설정 (립싱크 안정화)
		"-audio_buffer_size", "30", // 오디오 버퍼 크기 (ms)
		"-i", fmt.Sprintf("video=%s:audio=%s", video.Name, audio.Name),
		"-video_size", "1920x1080", // 입력 해상도
		"-vf", "scale=1920:1080,setdar=16/9", // 출력 해상도 고정 + 화면비 16:9 강제
		"-window_title", "Preview", // 미리보기 창 제목
		"-x", "1920", "-y", "1080", // 창 크기
		"-sync", "video", // 비디오 기준으로 동기화 (오디오가 비디오에 맞춤)
		"-af", "adelay=0|0", // 오디오 지연 없음 (좌/우 채널 모두 0ms)
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

package ffmpeg

import (
	"bufio"
	"fmt"
	"os/exec"
	"strings"

	"capture-board-selector/internal/captureboard/domain"
	"capture-board-selector/internal/captureboard/infra/logger"
)

type FFplayRunner struct{}

func NewFFplayRunner() domain.PreviewRunner {
	return &FFplayRunner{}
}

func (r *FFplayRunner) Run(video domain.Device, audio domain.Device) error {
	ffplay := resolveExe("ffplay", FFplayExe())
	cmd := exec.Command(ffplay,
		"-hide_banner",                                                                   // 버전/설정 출력 숨김
		"-loglevel", "error",                                                             // 에러 외 로그 숨김
		"-f", "dshow",                                                                    // Windows DirectShow 입력
		"-fflags", "nobuffer",                                                            // 입력 버퍼링 비활성화 (지연 감소)
		"-flags", "low_delay",                                                            // 저지연 모드
		"-avioflags", "direct",                                                           // I/O 버퍼링 비활성화
		"-rtbufsize", "512M",                                                             // 실시간 입력 버퍼 크기
		"-af", "aresample=async=1:first_pts=0",                                           // 오디오 리샘플 (비동기 보정)
		"-use_wallclock_as_timestamps", "1",                                              // 시스템 시계 기준 타임스탬프 (립싱크 안정화)
		"-audio_buffer_size", "50",                                                       // 오디오 버퍼 크기 (ms)
		"-async", "1",                                                                    // 오디오 샘플 수 기준 동기화
		"-af", "aresample=resampler=soxr:osf=s32:async=1:min_comp=0.001:first_pts=0",    // soxr 리샘플러 (고품질 동기화)
		"-i", fmt.Sprintf("video=%s:audio=%s", video.Name, audio.Name),
		"-video_size", "1920x1080",                                                       // 입력 해상도
		"-vf", "scale=1920:1080,setdar=16/9",                                            // 출력 해상도 고정 + 화면비 16:9 강제
		"-window_title", "Preview",                                                       // 미리보기 창 제목
		"-x", "1920", "-y", "1080",                                                      // 창 크기
		"-af", "adelay=0|0",                                                              // 오디오 지연 없음 (좌/우 채널 0ms)
	)
	stderr, err := cmd.StderrPipe()
	if err != nil {
		return cmd.Run()
	}

	if err := cmd.Start(); err != nil {
		return err
	}

	scanner := bufio.NewScanner(stderr)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line != "" {
			logger.Log("ffplay", line)
		}
	}

	return cmd.Wait()
}

// resolveExe returns the system command if available, otherwise falls back to the APPDATA path.
func resolveExe(name, fallback string) string {
	if _, err := exec.LookPath(name); err == nil {
		return name
	}
	return fallback
}

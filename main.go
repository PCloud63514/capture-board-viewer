package main

import (
	"bufio"
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"strings"
	"time"

	"github.com/eiannone/keyboard"
)

func main() {
	clear()
	fmt.Println("===============================================")
	fmt.Println("   Capture Board Selector")
	fmt.Println("===============================================")

	if !commandExists("ffmpeg") || !commandExists("ffplay") {
		fmt.Println("[ERROR] ffmpeg 또는 ffplay를 찾을 수 없습니다.")
		pause()
		return
	}

	fmt.Println("[INFO] 장치 검색 중...")

	timestamp := time.Now().Format("20060102_150405")
	_ = os.MkdirAll("logs", 0755)
	logFile := fmt.Sprintf("logs/device_output_%s.txt", timestamp)

	cmd := exec.Command("ffmpeg", "-hide_banner", "-list_devices", "true", "-f", "dshow", "-i", "dummy")
	var stderr bytes.Buffer
	cmd.Stderr = &stderr
	_ = cmd.Run()

	output := stderr.String()
	_ = os.WriteFile(logFile, []byte(output), 0644)

	lines := strings.Split(output, "\n")
	var videos, audios []string

	for _, line := range lines {
		line = strings.TrimSpace(strings.ReplaceAll(line, "\uFEFF", ""))
		if strings.Contains(line, "(video)") {
			videos = append(videos, extractDeviceName(line))
		} else if strings.Contains(line, "(audio)") {
			audios = append(audios, extractDeviceName(line))
		}
	}

	if len(videos) == 0 || len(audios) == 0 {
		fmt.Println("[ERROR] 캡처보드 장치가 연결되어 있지 않습니다.")
		fmt.Println("[DEBUG] 로그 확인:", logFile)
		pause()
		return
	}

	video := selectDevice("비디오 장치 선택", videos, "")
	audio := selectDevice("오디오 장치 선택", audios, video)

	fmt.Printf("[INFO] VIDEO = %s\n", video)
	fmt.Printf("[INFO] AUDIO = %s\n", audio)
	fmt.Println("-----------------------------------------------")

	play := exec.Command("ffplay",
		"-hide_banner",
		"-loglevel", "error",
		"-f", "dshow",
		"-fflags nobuffer",
		"-flags", "low_delay",
		"-avioflags", "direct",
		"-rtbufsize", "512M",
		"-af", "aresample=resampler=soxr:osf=s32:async=1:min_comp=0.001:first_pts=0,adelay=0|0",
		"-use_wallclock_as_timestamps", "1",
		"-audio_buffer_size", "50",
		"-async", "1",
		"-i", fmt.Sprintf("video=%q:audio=%q", video, audio),
		"-video_size", "1920x1080",
		"-vf", "scale=1920:1080,setdar=16/9",
		"-window_title", "Preview",
		"-x", "1920", "-y", "1080",
	)
	play.Stdout = os.Stdout
	play.Stderr = os.Stderr
	_ = play.Run()

	fmt.Println("-----------------------------------------------")
	fmt.Println("[INFO] 프로그램이 종료되었습니다.")
	pause()
}

func selectDevice(title string, items []string, selectedVideo string) string {
	_ = keyboard.Open()
	defer keyboard.Close()

	index := 0

	for {
		clear()
		fmt.Println("===============================================")
		fmt.Println("   Capture Board Selector")
		fmt.Println("===============================================")
		fmt.Printf("▶ %s\n", title)
		if selectedVideo != "" {
			fmt.Printf("   선택된 비디오 장치: %s\n", selectedVideo)
		}
		fmt.Println("-----------------------------------------------")

		for i, item := range items {
			if i == index {
				fmt.Printf(" > %s\n", item)
			} else {
				fmt.Printf("   %s\n", item)
			}
		}

		fmt.Println("\n↑↓ 이동, Enter 선택, ESC 종료")

		char, key, _ := keyboard.GetKey()
		if key == keyboard.KeyArrowUp && index > 0 {
			index--
		} else if key == keyboard.KeyArrowDown && index < len(items)-1 {
			index++
		} else if key == keyboard.KeyEnter {
			clear()
			return items[index]
		} else if key == keyboard.KeyEsc || char == 'q' {
			fmt.Println("취소되었습니다.")
			os.Exit(0)
		}
	}
}

func extractDeviceName(line string) string {
	start := strings.Index(line, "\"")
	end := strings.LastIndex(line, "\"")
	if start >= 0 && end > start {
		return line[start+1 : end]
	}
	return line
}

func commandExists(name string) bool {
	_, err := exec.LookPath(name)
	return err == nil
}

func clear() {
	if runtime.GOOS == "windows" {
		cmd := exec.Command("cmd", "/c", "cls")
		cmd.Stdout = os.Stdout
		cmd.Run()
	} else {
		fmt.Print("\033[H\033[2J")
	}
}

func pause() {
	fmt.Println("Press Enter to continue . . .")
	bufio.NewReader(os.Stdin).ReadBytes('\n')
}

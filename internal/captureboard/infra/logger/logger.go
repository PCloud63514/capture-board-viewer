package logger

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"time"
)

const maxLogFiles = 15

func logDir() string {
	return filepath.Join(os.Getenv("APPDATA"), "capture-board-viewer", "logs")
}

func LogDir() string {
	return logDir()
}

func Write(name, content string) error {
	dir := logDir()
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("로그 디렉토리 생성 실패: %w", err)
	}

	timestamp := time.Now().Format("20060102_150405")
	filename := filepath.Join(dir, fmt.Sprintf("%s_%s.log", name, timestamp))
	if err := os.WriteFile(filename, []byte(content), 0644); err != nil {
		return fmt.Errorf("로그 파일 저장 실패: %w", err)
	}

	return rotate(dir)
}

func rotate(dir string) error {
	entries, err := os.ReadDir(dir)
	if err != nil {
		return err
	}

	var files []os.FileInfo
	for _, e := range entries {
		if e.IsDir() {
			continue
		}
		info, err := e.Info()
		if err != nil {
			continue
		}
		files = append(files, info)
	}

	if len(files) <= maxLogFiles {
		return nil
	}

	sort.Slice(files, func(i, j int) bool {
		return files[i].ModTime().Before(files[j].ModTime())
	})

	for _, f := range files[:len(files)-maxLogFiles] {
		os.Remove(filepath.Join(dir, f.Name()))
	}
	return nil
}

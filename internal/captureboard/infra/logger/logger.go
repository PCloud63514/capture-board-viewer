package logger

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"sync"
	"time"
)

const maxLogFiles = 15

var (
	mu   sync.Mutex
	file *os.File
)

func logDir() string {
	return filepath.Join(os.Getenv("APPDATA"), "capture-board-viewer", "logs")
}

func LogDir() string { return logDir() }

func Start(version string) {
	mu.Lock()
	defer mu.Unlock()

	if version == "" {
		version = "dev"
	}

	dir := logDir()
	if err := os.MkdirAll(dir, 0755); err != nil {
		return
	}

	ts := time.Now().Format("20060102_150405")
	path := filepath.Join(dir, fmt.Sprintf("session_%s.log", ts))
	f, err := os.Create(path)
	if err != nil {
		return
	}
	file = f

	fmt.Fprintf(file, "====================================\n")
	fmt.Fprintf(file, "Capture Board Viewer %s\n", version)
	fmt.Fprintf(file, "시작: %s\n", time.Now().Format("2006-01-02 15:04:05"))
	fmt.Fprintf(file, "====================================\n\n")

	rotate(dir)
}

func Log(section, msg string) {
	mu.Lock()
	defer mu.Unlock()
	if file == nil {
		return
	}
	fmt.Fprintf(file, "[%s] [%s] %s\n", time.Now().Format("15:04:05"), section, msg)
}

func Close() {
	mu.Lock()
	defer mu.Unlock()
	if file != nil {
		file.Close()
		file = nil
	}
}

func rotate(dir string) {
	entries, err := os.ReadDir(dir)
	if err != nil {
		return
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
		return
	}
	sort.Slice(files, func(i, j int) bool {
		return files[i].ModTime().Before(files[j].ModTime())
	})
	for _, f := range files[:len(files)-maxLogFiles] {
		os.Remove(filepath.Join(dir, f.Name()))
	}
}

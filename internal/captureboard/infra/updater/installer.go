package updater

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"

	"capture-board-selector/internal/captureboard/domain"
)

// updateScript는 exe 교체 후 재시작하는 배치 스크립트입니다. exe 안에 내장됩니다.
const updateScript = `@echo off
:wait
timeout /t 1 /nobreak >nul
move /y "%~1" "%~2" >nul 2>&1
if errorlevel 1 goto wait
start "" "%~2"
(goto) 2>nul & del "%~f0"
`

type Installer struct{}

func NewInstaller() domain.UpdateInstaller {
	return &Installer{}
}

func (i *Installer) Install(info domain.UpdateInfo, onProgress func(downloaded, total int64)) error {
	currentExe, err := os.Executable()
	if err != nil {
		return fmt.Errorf("현재 실행 파일 경로를 가져올 수 없습니다: %w", err)
	}

	newExe := currentExe + "_update.exe"
	if err = download(info.DownloadURL, newExe, onProgress); err != nil {
		return err
	}

	scriptPath := filepath.Join(os.TempDir(), "cbv_update.bat")
	if err = os.WriteFile(scriptPath, []byte(updateScript), 0644); err != nil {
		return fmt.Errorf("업데이트 스크립트 생성 실패: %w", err)
	}

	cmd := exec.Command("cmd", "/c", "start", "", "/min", scriptPath, newExe, currentExe)
	if err = cmd.Start(); err != nil {
		return fmt.Errorf("업데이트 실행 실패: %w", err)
	}

	os.Exit(0)
	return nil
}

func download(url, dest string, onProgress func(downloaded, total int64)) error {
	resp, err := http.Get(url)
	if err != nil {
		return fmt.Errorf("다운로드 실패: %w", err)
	}
	defer resp.Body.Close()

	total := resp.ContentLength
	buf := &bytes.Buffer{}
	if _, err = io.Copy(buf, &progressReader{r: resp.Body, total: total, onProgress: onProgress}); err != nil {
		return fmt.Errorf("다운로드 중 오류: %w", err)
	}

	return os.WriteFile(dest, buf.Bytes(), 0755)
}

type progressReader struct {
	r          io.Reader
	downloaded int64
	total      int64
	onProgress func(downloaded, total int64)
}

func (pr *progressReader) Read(p []byte) (int, error) {
	n, err := pr.r.Read(p)
	pr.downloaded += int64(n)
	if pr.onProgress != nil {
		pr.onProgress(pr.downloaded, pr.total)
	}
	return n, err
}

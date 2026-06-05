package ffmpeg

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"capture-board-selector/internal/captureboard/domain"
	"github.com/bodgit/sevenzip"
)

const downloadURL = "https://www.gyan.dev/ffmpeg/builds/ffmpeg-release-full.7z"

type Installer struct{}

func NewInstaller() domain.FFmpegInstaller {
	return &Installer{}
}

func (i *Installer) Install(onProgress func(downloaded, total int64)) error {
	resp, err := http.Get(downloadURL)
	if err != nil {
		return fmt.Errorf("다운로드 실패: %w", err)
	}
	defer resp.Body.Close()

	total := resp.ContentLength
	buf := &bytes.Buffer{}
	if _, err = io.Copy(buf, &progressReader{r: resp.Body, total: total, onProgress: onProgress}); err != nil {
		return fmt.Errorf("다운로드 중 오류: %w", err)
	}

	data := buf.Bytes()
	r, err := sevenzip.NewReader(bytes.NewReader(data), int64(len(data)))
	if err != nil {
		return fmt.Errorf("압축 해제 실패: %w", err)
	}

	binDir := filepath.Join(installDir(), "bin")
	if err = os.MkdirAll(binDir, 0755); err != nil {
		return fmt.Errorf("디렉토리 생성 실패: %w", err)
	}

	for _, f := range r.File {
		base := filepath.Base(f.Name)
		if !strings.EqualFold(base, "ffmpeg.exe") && !strings.EqualFold(base, "ffplay.exe") {
			continue
		}
		if err = extractFile(f, filepath.Join(binDir, base)); err != nil {
			return fmt.Errorf("%s 추출 실패: %w", base, err)
		}
	}
	return nil
}

func extractFile(f *sevenzip.File, dest string) error {
	rc, err := f.Open()
	if err != nil {
		return err
	}
	defer rc.Close()

	out, err := os.Create(dest)
	if err != nil {
		return err
	}
	defer out.Close()

	_, err = io.Copy(out, rc)
	return err
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

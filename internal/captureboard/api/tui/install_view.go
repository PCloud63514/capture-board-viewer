package tui

import (
	"fmt"
	"strings"

	"github.com/rivo/tview"
)

const approxSizeMB = 150

func (a *App) showInstallPrompt() {
	modal := tview.NewModal().
		SetText(fmt.Sprintf(
			"FFmpeg가 설치되어 있지 않습니다.\n\n약 %dMB를 다운로드합니다.\n계속하시겠습니까?",
			approxSizeMB,
		)).
		AddButtons([]string{"설치", "종료"}).
		SetDoneFunc(func(idx int, _ string) {
			if idx == 0 {
				a.startInstall()
			} else {
				a.tv.Stop()
			}
		})
	a.pages.AddAndSwitchToPage("install_prompt", modal, true)
}

func (a *App) startInstall() {
	statusText := tview.NewTextView().
		SetDynamicColors(true).
		SetTextAlign(tview.AlignCenter)

	box := tview.NewFlex().SetDirection(tview.FlexRow).
		AddItem(tview.NewTextView().SetTextAlign(tview.AlignCenter).SetText("FFmpeg 설치 중..."), 3, 0, false).
		AddItem(statusText, 3, 0, false)
	box.SetBorder(true).SetTitle(" FFmpeg 설치 ")

	a.pages.AddAndSwitchToPage("install_progress", box, true)

	go func() {
		err := a.installFFmpeg.Execute(func(downloaded, total int64) {
			a.tv.QueueUpdateDraw(func() {
				statusText.SetText(buildProgressText(downloaded, total))
			})
		})

		a.tv.QueueUpdateDraw(func() {
			if err != nil {
				a.showError("설치 실패", err.Error())
				return
			}
			a.showDeviceSelection()
		})
	}()
}

func (a *App) showError(title, message string) {
	modal := tview.NewModal().
		SetText(fmt.Sprintf("[red]%s[-]\n\n%s", title, message)).
		AddButtons([]string{"종료"}).
		SetDoneFunc(func(_ int, _ string) { a.tv.Stop() })
	a.pages.AddAndSwitchToPage("error", modal, true)
}

func (a *App) showDeviceSelection() {
	// 다음 단계: 장치 선택 (추후 구현)
	modal := tview.NewModal().
		SetText("FFmpeg 준비 완료!\n\n장치 선택 화면은 다음 단계에서 구현됩니다.").
		AddButtons([]string{"종료"}).
		SetDoneFunc(func(_ int, _ string) { a.tv.Stop() })
	a.pages.AddAndSwitchToPage("ready", modal, true)
}

func buildProgressText(downloaded, total int64) string {
	const barWidth = 30
	var filled int
	if total > 0 {
		filled = int(float64(downloaded) / float64(total) * barWidth)
	}
	bar := strings.Repeat("█", filled) + strings.Repeat("░", barWidth-filled)
	pct := 0.0
	if total > 0 {
		pct = float64(downloaded) / float64(total) * 100
	}
	return fmt.Sprintf("[%s] %.1f%%\n%s / %s",
		bar, pct,
		formatBytes(downloaded), formatBytes(total),
	)
}

func formatBytes(b int64) string {
	if b < 1024*1024 {
		return fmt.Sprintf("%.1f KB", float64(b)/1024)
	}
	return fmt.Sprintf("%.1f MB", float64(b)/(1024*1024))
}

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
				a.showError("설치 실패", err)
				return
			}
			a.showDeviceSelection()
		})
	}()
}

func (a *App) showError(title string, err error) {
	text := fmt.Sprintf("[red]%s[-]\n\n%s", title, err.Error())
	modal := tview.NewModal().SetText(text)

	if a.reporter != nil {
		modal.AddButtons([]string{"오류 전송 후 종료", "그냥 종료"}).
			SetDoneFunc(func(idx int, _ string) {
				if idx == 0 {
					a.sendErrorReport(err)
				} else {
					a.tv.Stop()
				}
			})
	} else {
		modal.AddButtons([]string{"종료"}).
			SetDoneFunc(func(_ int, _ string) { a.tv.Stop() })
	}

	a.pages.AddAndSwitchToPage("error", modal, true)
}

func (a *App) sendErrorReport(err error) {
	sending := tview.NewModal().
		SetText("오류 정보를 개발자에게 전송하고 있습니다...\n\n잠시만 기다려 주세요.")
	a.pages.AddAndSwitchToPage("sending", sending, true)

	go func() {
		a.reporter.Report(err)
		a.reporter.Flush()
		a.tv.QueueUpdateDraw(func() { a.tv.Stop() })
	}()
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

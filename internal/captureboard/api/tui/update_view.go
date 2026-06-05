package tui

import (
	"fmt"
	"strings"

	"capture-board-selector/internal/captureboard/domain"

	"github.com/rivo/tview"
)

func (a *App) checkForUpdate(onDone func()) {
	if a.currentVersion == "" {
		onDone()
		return
	}

	checking := tview.NewModal().SetText("업데이트 확인 중...")
	a.pages.AddAndSwitchToPage("update_check", checking, true)

	go func() {
		info, err := a.checkUpdate.Execute(a.currentVersion)
		a.tv.QueueUpdateDraw(func() {
			if err != nil || info == nil {
				onDone()
				return
			}
			a.showUpdatePrompt(info, onDone)
		})
	}()
}

func (a *App) showUpdatePrompt(info *domain.UpdateInfo, onSkip func()) {
	modal := tview.NewModal().
		SetText(fmt.Sprintf("새 버전이 있습니다!\n\n현재: %s  →  최신: %s\n\n업데이트하시겠습니까?",
			a.currentVersion, info.Version)).
		AddButtons([]string{"업데이트", "나중에"}).
		SetDoneFunc(func(idx int, _ string) {
			if idx == 0 {
				a.startUpdate(info)
			} else {
				onSkip()
			}
		})
	a.pages.AddAndSwitchToPage("update_prompt", modal, true)
}

func (a *App) startUpdate(info *domain.UpdateInfo) {
	statusText := tview.NewTextView().
		SetDynamicColors(true).
		SetTextAlign(tview.AlignCenter)

	box := tview.NewFlex().SetDirection(tview.FlexRow).
		AddItem(tview.NewTextView().SetTextAlign(tview.AlignCenter).
			SetText(fmt.Sprintf("%s 다운로드 중...", info.Version)), 3, 0, false).
		AddItem(statusText, 3, 0, false)
	box.SetBorder(true).SetTitle(" 업데이트 설치 ")

	a.pages.AddAndSwitchToPage("update_progress", box, true)

	go func() {
		err := a.installUpdate.Execute(*info, func(downloaded, total int64) {
			a.tv.QueueUpdateDraw(func() {
				statusText.SetText(buildUpdateProgressText(downloaded, total))
			})
		})

		if err != nil {
			a.tv.QueueUpdateDraw(func() {
				a.showError("업데이트 실패", err)
			})
		}
	}()
}

func buildUpdateProgressText(downloaded, total int64) string {
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

package tui

import (
	"fmt"
	"strings"

	"capture-board-selector/internal/captureboard/domain"
	"capture-board-selector/internal/captureboard/infra/logger"

	"github.com/rivo/tview"
)

func (a *App) checkForUpdate(onDone func()) {
	if a.currentVersion == "" {
		logger.Log("업데이트", "버전 정보 없음, 확인 건너뜀")
		onDone()
		return
	}

	logger.Log("업데이트", "최신 버전 확인 중...")

	checking := tview.NewModal().SetText("업데이트 확인 중...")
	a.pages.AddAndSwitchToPage("update_check", checking, true)

	go func() {
		info, err := a.checkUpdate.Execute(a.currentVersion)
		a.tv.QueueUpdateDraw(func() {
			if err != nil {
				logger.Log("업데이트", fmt.Sprintf("확인 실패: %v", err))
				onDone()
				return
			}
			if info == nil {
				logger.Log("업데이트", "최신 버전 사용 중")
				onDone()
				return
			}
			logger.Log("업데이트", fmt.Sprintf("새 버전 발견: %s", info.Version))
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
				logger.Log("업데이트", fmt.Sprintf("사용자가 업데이트 선택: %s", info.Version))
				a.startUpdate(info)
			} else {
				logger.Log("업데이트", "사용자가 업데이트 건너뜀")
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
			logger.Log("업데이트", fmt.Sprintf("다운로드 실패: %v", err))
			a.tv.QueueUpdateDraw(func() {
				a.showError("업데이트 실패", err)
			})
		} else {
			logger.Log("업데이트", "다운로드 완료, 재시작 중...")
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

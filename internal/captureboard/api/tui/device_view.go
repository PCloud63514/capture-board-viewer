package tui

import (
	"capture-board-selector/internal/captureboard/domain"
	"fmt"

	"github.com/rivo/tview"
)

func (a *App) showDeviceSelection() {
	loading := tview.NewModal().SetText("장치 검색 중...")
	a.pages.AddAndSwitchToPage("loading", loading, true)

	go func() {
		devices, err := a.listDevices.Execute()
		a.tv.QueueUpdateDraw(func() {
			if err != nil {
				a.showError("장치 오류", err)
				return
			}

			var videos, audios []domain.Device
			for _, d := range devices {
				switch d.Type {
				case domain.DeviceTypeVideo:
					videos = append(videos, d)
				case domain.DeviceTypeAudio:
					audios = append(audios, d)
				}
			}

			if len(videos) == 0 {
				a.showError("장치 없음", fmt.Errorf("연결된 비디오 장치를 찾을 수 없습니다.\n캡처보드가 연결되어 있는지 확인하세요."))
				return
			}
			if len(audios) == 0 {
				a.showError("장치 없음", fmt.Errorf("연결된 오디오 장치를 찾을 수 없습니다."))
				return
			}

			a.showVideoSelection(videos, audios)
		})
	}()
}

func (a *App) showVideoSelection(videos, audios []domain.Device) {
	list := tview.NewList()
	for _, v := range videos {
		list.AddItem(v.Name, "video", 0, nil)
	}

	list.SetSelectedFunc(func(idx int, _ string, _ string, _ rune) {
		a.showAudioSelection(videos[idx], audios)
	})
	list.SetDoneFunc(func() { a.tv.Stop() })

	list.SetBorder(true).
		SetTitle(fmt.Sprintf(" 비디오 장치 선택  |  ESC: 종료  |  %s ", a.currentVersion)).
		SetTitleAlign(tview.AlignLeft)

	a.pages.AddAndSwitchToPage("video_select", list, true)
}

func (a *App) showAudioSelection(video domain.Device, audios []domain.Device) {
	list := tview.NewList()
	for _, au := range audios {
		list.AddItem(au.Name, "audio", 0, nil)
	}

	list.SetSelectedFunc(func(idx int, _ string, _ string, _ rune) {
		a.selectedVideo = &video
		a.selectedAudio = &audios[idx]
		a.tv.Stop()
	})
	list.SetDoneFunc(func() {
		a.pages.SwitchToPage("video_select")
	})

	list.SetBorder(true).
		SetTitle(fmt.Sprintf(" 오디오 장치 선택  |  선택된 비디오: %s  |  ESC: 뒤로 ", video.Name)).
		SetTitleAlign(tview.AlignLeft)

	a.pages.AddAndSwitchToPage("audio_select", list, true)
}

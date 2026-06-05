package main

import (
	"github.com/rivo/tview"
)

func main() {
	app := tview.NewApplication()

	list := tview.NewList().
		AddItem("Capture Card A", "video device", 'a', nil).
		AddItem("Capture Card B", "video device", 'b', nil).
		AddItem("종료", "", 'q', func() { app.Stop() })

	list.SetBorder(true).SetTitle(" 비디오 장치 선택 ").SetTitleAlign(tview.AlignLeft)

	if err := app.SetRoot(list, true).Run(); err != nil {
		panic(err)
	}
}

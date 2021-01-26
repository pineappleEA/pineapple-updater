package main

import (
	"strconv"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
	"fyne.io/fyne/v2/container"

	"github.com/cavaliercoder/grab"
)


func aboutUI() {
	a := fyne.CurrentApp()
	w := a.NewWindow("About")
	w.Resize(fyne.NewSize(400, 300))
	logo := canvas.NewImageFromResource(resourceIconPng)
	logo.FillMode = canvas.ImageFillOriginal
	quitButton := widget.NewButton("close", func() { w.Close() })
	aboutText1 := widget.NewLabelWithStyle("Project Dëfënëstrëring", fyne.TextAlignCenter, fyne.TextStyle{Bold: true})
	aboutText2 := widget.NewLabelWithStyle("\nFrom EmuWorld with love\n2021", fyne.TextAlignCenter, fyne.TextStyle{Italic: true})
	aboutText3 := widget.NewLabelWithStyle("\n\n\nThis program is free software; you can redistribute it and/or modify\nit under the terms of the GNU General Public License as published by\nthe Free Software Foundation; either version 2 of the License, or\n(at your option) any later version.", fyne.TextAlignCenter, fyne.TextStyle{})
	ui := fyne.NewContainerWithLayout(layout.NewBorderLayout(logo, quitButton, nil, nil), logo, quitButton, aboutText1, aboutText2, aboutText3)
	w.SetIcon(resourceIconPng)
	w.SetContent(ui)
	w.SetFixedSize(true)
	w.Show()
}

func loadUI(versionSlice []int, linkMap map[int]string) fyne.CanvasObject {

	list := widget.NewList(
		func() int { return len(versionSlice) },
		func() fyne.CanvasObject {
			return widget.NewLabel("This is a test")
		},
		func(id int, obj fyne.CanvasObject) {
			obj.(*widget.Label).SetText("EA " + strconv.Itoa(versionSlice[id]))
		},
	)
	var selectedVersion int = 0
	list.OnSelected = func(id int) { selectedVersion = id }

	buttonSide := container.New(
		layout.NewVBoxLayout(),
		widget.NewButton("Install", func() { install(versionSlice, linkMap, selectedVersion) }),
		widget.NewButton("Uninstall", func() {}),
	)

	downloadProgress := widget.NewProgressBar()
	downloadProgress.Hide()
	buttonFooter := container.New(
		layout.NewHBoxLayout(),
		widget.NewButtonWithIcon("", resourceIconPng, func() { go aboutUI() }),
		widget.NewButton("Settings", func() {}),
		downloadProgress,
	)
	//combine three elements into one container/canvas
	ui := fyne.NewContainerWithLayout(layout.NewBorderLayout(nil, buttonFooter, nil, buttonSide), buttonFooter, buttonSide, list)
	return ui
}

func downloadUI(resp *grab.Response) {
	a := fyne.CurrentApp()
	w := a.NewWindow("Downloading...")
	downloadProgress := widget.NewProgressBar()
	downloadSpeed := widget.NewLabel("")
	w.Resize(fyne.NewSize(400, 200))
	w.SetIcon(resourceIconPng)
	w.SetFixedSize(true)
	ui := fyne.NewContainerWithLayout(layout.NewBorderLayout(nil, nil, nil, nil), downloadProgress, downloadSpeed)
	w.SetContent(ui)
	w.Show()
	go func() {
		for {
			time.Sleep(time.Millisecond * 250)
			downloadProgress.SetValue(resp.Progress())
			downloadSpeed.SetText("Download Speed: " + strconv.Itoa(int(resp.BytesPerSecond()/1000)) + "KByte/s")
			if int(resp.Progress()) == 1 {
				w.Close()
				break
			}
		}
	}()

}

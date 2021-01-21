package main

import (
	"bufio"
	"fmt"
	"net/http"
	"regexp"
	"strconv"

	"fyne.io/fyne"
	"fyne.io/fyne/app"
	"fyne.io/fyne/canvas"
	"fyne.io/fyne/layout"
	"fyne.io/fyne/widget"
)

const pineappleSite string = "https://raw.githubusercontent.com/pineappleEA/pineappleEA.github.io/master/index.html"

func aboutUI() {
	a := fyne.CurrentApp()
	w := a.NewWindow("About")
	w.Resize(fyne.NewSize(400, 400))
	logo := canvas.NewImageFromResource(resourceIconPng)
	logo.FillMode = canvas.ImageFillOriginal
	quitButton := widget.NewButton("close", func() { w.Close() })
	aboutText1 := widget.NewLabelWithStyle("Project Dëfënëstrëring", fyne.TextAlignCenter, fyne.TextStyle{Bold: true})
	aboutText2 := widget.NewLabelWithStyle("\nFrom EmuWorld with love\n2021", fyne.TextAlignCenter, fyne.TextStyle{Italic: true})
	aboutText3 := widget.NewLabelWithStyle("\n\n\nThis program is free software; you can redistribute it and/or modify\nit under the terms of the GNU General Public License as published by\nthe Free Software Foundation; either version 2 of the License, or\n(at your option) any later version.", fyne.TextAlignCenter, fyne.TextStyle{})
	ui := fyne.NewContainerWithLayout(layout.NewBorderLayout(logo, quitButton, nil, nil), logo, quitButton, aboutText1, aboutText2, aboutText3)
	w.SetIcon(resourceIconPng)
	w.SetContent(ui)
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

	buttonSide := widget.NewVBox(
		widget.NewButton("Install", func() { install(versionSlice, linkMap, selectedVersion) }),
		widget.NewButton("Uninstall", func() {}),
	)

	buttonFooter := widget.NewHBox(
		widget.NewButtonWithIcon("", resourceIconPng, func() { go aboutUI() }),
		widget.NewButton("Settings", func() {}),
	)

	//combine three elements into one container/canvas
	ui := fyne.NewContainerWithLayout(layout.NewBorderLayout(nil, buttonFooter, nil, buttonSide), buttonFooter, buttonSide, list)
	return ui
}

func main() {
	a := app.New()
	w := a.NewWindow("PinEApple Updater")
	w.SetIcon(resourceIconPng)
	versionSlice, linkMap := downloadList()
	w.SetContent(loadUI(versionSlice, linkMap))
	w.Resize(fyne.NewSize(500, 450))
	w.ShowAndRun()
}

func downloadList() ([]int, map[int]string) {
	//return variables
	linkMap := make(map[int]string)
	versionSlice := make([]int, 0)

	//download site into resp
	resp, err := http.Get(pineappleSite)
	if err != nil {
		// handle err
	}
	defer resp.Body.Close()

	//read response body through scanner
	scanner := bufio.NewScanner(resp.Body)
	for i := 0; scanner.Scan(); i++ {
		var line = scanner.Text()
		match, _ := regexp.MatchString("https://anonfiles.com", line)
		if match {
			// extract link
			linkPattern, _ := regexp.Compile("https://anonfiles.com/.*/YuzuEA-[0-9]*_7z")
			link := linkPattern.FindString(scanner.Text())

			// extract version number
			versionPattern, _ := regexp.Compile("EA [0-9]*")
			versionString := versionPattern.FindString(scanner.Text())
			numberPattern, _ := regexp.Compile("[0-9]*$")
			versionString = numberPattern.FindString(versionString)
			version, _ := strconv.Atoi(versionString)

			//save link in map
			linkMap[version] = link
			//add version number to slice
			versionSlice = append(versionSlice, version)

		} else if line == "</html>" {
			break
		}
	}
	return versionSlice, linkMap
}

func install(versionSlice []int, linkMap map[int]string, selectedVersion int) {
	fmt.Println(versionSlice[selectedVersion])
	fmt.Println(linkMap[versionSlice[selectedVersion]])
}
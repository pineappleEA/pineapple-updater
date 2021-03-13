package main

import (
	"bufio"
	"net/http"
	"os"
	"regexp"
	"strconv"
	"context"
	"log"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"

	"github.com/cavaliercoder/grab"
)

const pineappleSrc string = "https://github.com/pineappleEA/pineapple-src/"
const pineappleSite string = "https://pineappleEA.github.io/"
//TODO: set actually usable default install path
const defaultPath string = "C:/yuzu"

func main() {
	logfile, err := os.Create("log.txt")
	if err != nil {
		log.Fatal(err)
	}
	log.SetOutput(logfile)
	a := app.NewWithID("pinEApple updater")
	w := a.NewWindow("PinEApple Updater")
	w.SetIcon(resourceIconPng)
	log.Println("Downloading available versions from pineappleea.github.io")
	versionSlice, linkMap := downloadList()
	w.SetContent(mainUI(versionSlice, linkMap))
	w.Resize(fyne.NewSize(500, 450))
	w.Show()
	a.Run()
}

func downloadList() ([]int, map[int]string) {
	//return variables
	linkMap := make(map[int]string)
	versionSlice := make([]int, 0)

	//download site into resp
	resp, err := http.Get(pineappleSite)
	if err != nil {
		log.Fatal(os.Stderr, "Could not obtain list of versions!\n", err)
	}
	defer resp.Body.Close()

	//read response body through scanner
	scanner := bufio.NewScanner(resp.Body)
	for i := 0; scanner.Scan(); i++ {
		var line = scanner.Text()
		match, _ := regexp.MatchString("EA [0-9]", line)
		// extract version number
		versionPattern, _ := regexp.Compile("EA [0-9]*")
		versionString := versionPattern.FindString(scanner.Text())
		numberPattern, _ := regexp.Compile("[0-9]*$")
		versionString = numberPattern.FindString(versionString)
		version, _ := strconv.Atoi(versionString)
		if match {
			// extract link
			linkPattern, _ := regexp.Compile("https://anonfiles.com/.*/YuzuEA-[0-9]*_7z")
			link := linkPattern.FindString(scanner.Text())

			//save link in map
			linkMap[version] = link
			//add version number to slice
			versionSlice = append(versionSlice, version)

		} else if line == "</html>" {
			break
		}
	}
	if len(versionSlice) <= 1 {
		log.Fatal(os.Stderr, "Could not obtain list of files!\n")
	}
	log.Println("Found "+strconv.Itoa(len(versionSlice))+" versions")
	return versionSlice, linkMap
}

func install(versionSlice []int, linkMap map[int]string, selectedVersion int) {
	log.Println("Trying to install version "+strconv.Itoa(versionSlice[selectedVersion]))
	resp, _ := http.Get(pineappleSrc + "releases/download/EA-" + strconv.Itoa(versionSlice[selectedVersion]) + "/Windows-Yuzu-EA-" + strconv.Itoa(versionSlice[selectedVersion]) + ".7z")
	defer resp.Body.Close()
	var downloadLink string
	if resp.StatusCode == 200 {
		// Downloading from Github
		downloadLink = pineappleSrc + "releases/download/EA-" + strconv.Itoa(versionSlice[selectedVersion]) + "/Windows-Yuzu-EA-" + strconv.Itoa(versionSlice[selectedVersion]) + ".7z"
	} else {
		//Download from Anonfiles
		//Download Anonfiles page to grab direct download
		resp, err := http.Get(linkMap[versionSlice[selectedVersion]])
		if err != nil {
			log.Fatal(os.Stderr, "Neither GDrive nor Anonfiles responds! Exiting...\n", err)
		}
		//go line through line and search for direct download link with regex
		scanner := bufio.NewScanner(resp.Body)
		for scanner.Scan() {
			linkPattern, _ := regexp.Compile("https://cdn-.*anonfiles.*7z")
			if linkPattern.MatchString(scanner.Text()) {
				downloadLink = linkPattern.FindString(scanner.Text())
				break
			}
		}
		//exit if no download link found
		if downloadLink == "" {
			log.Fatal(os.Stderr, "No download link found, Anonfiles or Github seems to be having issues! Exiting...\n")
		}
		defer resp.Body.Close()
	}
	downloadFile(downloadLink)
}

//Downloads file from given link to set path
func downloadFile(link string) {
	log.Println("Downloading from "+link)
	//TODO: figure out proper way to set the path for windows
	req, _ := grab.NewRequest(fyne.CurrentApp().Preferences().StringWithFallback("path",defaultPath), link)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	
	req = req.WithContext(ctx)
	resp := grab.DefaultClient.Do(req)
	//TODO: figure out why the mainUI is unresponsive when the downloadUI is open
	go downloadUI(resp, cancel)

	// check for errors
	if err := resp.Err(); err != nil && err.Error() != "context canceled" {
		log.Fatal(os.Stderr, "Download failed: %v\n", err)
	}

}

package tabs

import (
	"fmt"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
	"net/url"
	"runtime"
	"strconv"
	"strings"
)

type downloadState = int

const (
	DownloadNotStarted downloadState = iota
	Downloading
	DownloadPaused
	DownloadFinished
	DownloadCanceled
	DownloadUnknown
)

type downloadItem struct {
	Name       string
	Icon       string
	Url        url.URL
	Size       uint
	TotalSize  uint
	RemainTime string
	Dir        string
	State      downloadState
	Err        error
	RowID      int
}

var (
	globalIndex           int
	isDownloading         = false
	Items                 []downloadItem
	list                  *widget.List
	downloadStateToString = map[int]string{
		DownloadNotStarted: "Not started",
		Downloading:        "Downloading...",
		DownloadPaused:     "Paused",
		DownloadFinished:   "Finished",
		DownloadCanceled:   "Canceled",
		DownloadUnknown:    "Unknown state",
	}
)

func newDownloadListArea() *widget.List {
	list = widget.NewList(countDownloadItems, newDownloadItemArea, updateDownloadItemArea)
	return list
}

func newDownloadItemArea() fyne.CanvasObject {
	if countDownloadItems() < 1 {
		return &widget.BaseWidget{}
	}
	defer func() {
		pc, _, _, _ := runtime.Caller(2)
		if strings.HasSuffix(runtime.FuncForPC(pc).Name(), "Refresh") {
			return
		}
		if globalIndex >= countDownloadItems() {
			globalIndex = 0
		} else {
			globalIndex += 1
		}
	}()
	currentIndex := globalIndex
	newItem := Items[countDownloadItems()-1]
	//icon := widget.NewLabel(newItem.Name)

	name := widget.NewLabel(newItem.Name)
	size := widget.NewLabel(fmt.Sprintf("%d/%d", newItem.Size, newItem.TotalSize))
	lVBox := container.NewVBox(name, size)

	time := widget.NewLabel(newItem.RemainTime)

	downloadProgressBar := widget.NewProgressBar()
	statusLabel := widget.NewLabel("downloading...")
	statusVBox := container.NewVBox(layout.NewSpacer(), downloadProgressBar, statusLabel, layout.NewSpacer())

	startButton := widget.NewButton("start", func() { startDownload(currentIndex) })
	cancelButton := widget.NewButton("cancel", func() { cancelDownload(currentIndex) })
	openDirButton := widget.NewButton("open", func() { openDir(currentIndex) })
	listControlBox := container.NewHBox(startButton, cancelButton, openDirButton)

	borderRightHBox := container.NewHBox(time, statusVBox, container.NewVBox(layout.NewSpacer(), listControlBox, layout.NewSpacer()))

	//return container.NewBorder(nil, nil, icon, borderRightHBox, lVBox)
	newItem.RowID = currentIndex
	return container.NewBorder(nil, nil, nil, borderRightHBox, lVBox)
}

func updateDownloadItemArea(id widget.ListItemID, item fyne.CanvasObject) {
	//downloadItem{
	//	Name:       "123",
	//	Icon:       "icon",
	//	Url:        url.URL{},
	//	Size:       1,
	//	TotalSize:  10,
	//	RemainTime: "--:--:--",
	//	Dir:        "dir",
	//	State:      0,
	//	Err:        nil,
	//})
	dataItem := Items[id]

	// Get widgets.
	border := item.(*fyne.Container)
	lvBox := border.Objects[0].(*fyne.Container)
	name := lvBox.Objects[0].(*widget.Label)
	size := lvBox.Objects[1].(*widget.Label)
	//icon := border.Objects[1]
	borderRightHBox := border.Objects[1].(*fyne.Container)
	timeLabel := borderRightHBox.Objects[0].(*widget.Label)
	statusVBox := borderRightHBox.Objects[1].(*fyne.Container)
	downloadProgressBar := statusVBox.Objects[1].(*widget.ProgressBar)
	statusLabel := statusVBox.Objects[2].(*widget.Label)
	//listControlBox := borderRightHBox.Objects[2].(*fyne.Container)

	// Set widget data.
	name.SetText(dataItem.Name)
	size.SetText(fmt.Sprintf("%d/%d", dataItem.Size, dataItem.TotalSize))
	downloadProgressBar.SetValue(float64(dataItem.Size) / float64(dataItem.TotalSize))
	statusLabel.SetText(downloadStateToString[dataItem.State])
	timeLabel.SetText(dataItem.RemainTime)
}

func newDownloadControlArea() *fyne.Container {
	addDownloadButton := widget.NewButton("Add", addDownload)
	label := widget.NewLabel("Downloading:")
	totalProgressBar := widget.NewProgressBar()
	totalProgressBar.Resize(fyne.NewSize(1000, 100))
	startButton := widget.NewButton("Start All", startDownloadAll)
	pauseButton := widget.NewButton("Pause All", pauseDownloadAll)
	cancelButton := widget.NewButton("Cancel ALl", cancelDownloadAll)
	controlHBox := container.NewBorder(nil, nil, container.NewHBox(addDownloadButton, label), container.NewHBox(startButton, pauseButton, cancelButton), totalProgressBar)
	return controlHBox
}

func NewDownloadTab() *container.TabItem {
	controlArea := newDownloadControlArea()
	listArea := newDownloadListArea()
	// TODO: Handle empty list.
	downloadTab := container.NewTabItem("Download", container.NewBorder(controlArea, nil, nil, nil, listArea))
	//downloadTab := container.NewTabItem("Download", container.NewVBox(controlArea))
	return downloadTab
}

var c = 1

func addDownload() {
	Items = append(Items, downloadItem{
		Name:       strconv.Itoa(c),
		Icon:       "icon",
		Url:        url.URL{},
		Size:       1,
		TotalSize:  10,
		RemainTime: "--:--:--",
		Dir:        "dir",
		State:      0,
		Err:        nil,
	})
	c++
	//globalIndex = 0
	Update()
}

func startDownloadAll() {

}

func pauseDownloadAll() {

}

func cancelDownloadAll() {

}

func pauseDownload(id int) {
	fmt.Println("Download paused", id)
}

func startDownload(id int) {
	fmt.Println("Download continued", id)
}

func cancelDownload(id int) {
	fmt.Println("Download canceled.", id)
}

func openDir(id int) {
	fmt.Println("Open dir", id)
}

func countDownloadItems() int {
	return len(Items)
}

func Update() {
	list.Refresh()
}

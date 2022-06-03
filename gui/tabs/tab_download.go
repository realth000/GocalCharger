package tabs

import (
	"fmt"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
	"net/url"
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

var (
	downloadStateToString = map[int]string{
		DownloadNotStarted: "Not started",
		Downloading:        "Downloading...",
		DownloadPaused:     "Paused",
		DownloadFinished:   "Finished",
		DownloadCanceled:   "Canceled",
		DownloadUnknown:    "Unknown state",
	}
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
}

var (
	isDownloading = false
	Items         []downloadItem
	list          *widget.List
)

func newDownloadListArea() *widget.List {
	list = widget.NewList(countDownloadItems, newDownloadItemArea, updateDownloadItemArea)
	return list
}

func newDownloadItemArea() fyne.CanvasObject {
	newItem := Items[countDownloadItems()-1]
	//icon := widget.NewLabel(newItem.Name)

	name := widget.NewLabel(newItem.Name)
	size := widget.NewLabel(fmt.Sprintf("%d/%d", newItem.Size, newItem.TotalSize))
	lVBox := container.NewVBox(name, size)

	time := widget.NewLabel(newItem.RemainTime)

	downloadProgressBar := widget.NewProgressBar()
	statusLabel := widget.NewLabel("downloading...")
	statusVBox := container.NewVBox(layout.NewSpacer(), downloadProgressBar, statusLabel, layout.NewSpacer())

	startButton := widget.NewButton("start", startDownload)
	cancelButton := widget.NewButton("cancel", cancelDownload)
	openDirButton := widget.NewButton("open", openDir)
	listControlBox := container.NewHBox(startButton, cancelButton, openDirButton)

	borderRightHBox := container.NewHBox(time, statusVBox, container.NewVBox(layout.NewSpacer(), listControlBox, layout.NewSpacer()))

	//return container.NewBorder(nil, nil, icon, borderRightHBox, lVBox)
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
	label := widget.NewLabel("Downloading:")
	totalProgressBar := widget.NewProgressBar()
	totalProgressBar.Resize(fyne.NewSize(1000, 100))
	startButton := widget.NewButton("Start All", startDownload)
	pauseButton := widget.NewButton("Pause All", pauseDownload)
	cancelButton := widget.NewButton("Cancel ALl", cancelDownload)
	controlHBox := container.NewBorder(nil, nil, label, container.NewHBox(startButton, pauseButton, cancelButton), totalProgressBar)
	return controlHBox
}

func NewDownloadTab() *container.TabItem {
	controlArea := newDownloadControlArea()
	Items = append(Items, downloadItem{
		Name:       "123",
		Icon:       "icon",
		Url:        url.URL{},
		Size:       1,
		TotalSize:  10,
		RemainTime: "--:--:--",
		Dir:        "dir",
		State:      0,
		Err:        nil,
	})
	listArea := newDownloadListArea()

	downloadTab := container.NewTabItem("Download", container.NewVBox(controlArea, listArea))
	fmt.Println(listArea.CreateItem())
	fmt.Println(listArea.CreateItem())
	return downloadTab
}

func pauseDownload() {
	fmt.Println("Download paused")
}

func startDownload() {
	fmt.Println("Download continued")
}

func cancelDownload() {
	fmt.Println("Download canceled.")
}

func openDir() {
	fmt.Println("Open dir")
}

func countDownloadItems() int {
	return len(Items)
}

func Update() {
	list.Refresh()
}

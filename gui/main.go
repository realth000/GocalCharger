package main

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"gocalcharger/gui/tabs"
	"time"
)

func updateTime() {
	if len(tabs.Items) > 0 {
		//tabs.Items[len(tabs.Items)-1].RemainTime = time.Now().Format("Time: 03:04:05")
		//tabs.Update()
	}
}

func main() {
	a := app.New()
	w := a.NewWindow("GocalChargerGui")
	w.Resize(fyne.NewSize(800, 600))
	downloadTab := tabs.NewDownloadTab()
	networkTab := tabs.NewNetworkTab()
	tab := container.NewAppTabs(downloadTab, networkTab)
	w.SetContent(tab)
	w.Show()

	go func() {
		for range time.Tick(time.Second) {
			updateTime()
		}
	}()
	a.Run()
}

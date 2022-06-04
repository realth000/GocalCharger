package tabs

import (
	"fmt"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
)

func makeServerSettingsArea() fyne.CanvasObject {
	configFileLabel := widget.NewLabel("Config file path")
	configFileEntry := widget.NewEntry()
	configFileEntry.SetPlaceHolder("*.toml")
	formLayout := container.New(layout.NewFormLayout(), configFileLabel, configFileEntry)
	loadConfigFileButton := widget.NewButton("Load config file", reloadServerConfigFile)
	buttonBox := container.NewHBox(container.NewVBox(loadConfigFileButton))
	return widget.NewCard("Server", "Server settings", container.NewVBox(formLayout, buttonBox))
}

func makeClientSettingsArea() fyne.CanvasObject {
	configFileLabel := widget.NewLabel("Config file path")
	configFileEntry := widget.NewEntry()
	configFileEntry.SetPlaceHolder("*.toml")
	formLayout := container.New(layout.NewFormLayout(), configFileLabel, configFileEntry)
	loadConfigFileButton := widget.NewButton("Load config file", reloadClientConfigFile)
	buttonBox := container.NewHBox(container.NewVBox(loadConfigFileButton))
	return widget.NewCard("Client", "Client settings", container.NewVBox(formLayout, buttonBox))
}

func reloadServerConfigFile() {
	fmt.Println("reload server config file")
}

func reloadClientConfigFile() {
	fmt.Println("reload server config file")
}
func NewSettingsTab() *container.TabItem {
	serverSettings := makeServerSettingsArea()
	clientSettings := makeClientSettingsArea()
	settingsTab := container.NewTabItem("Settins", container.NewVBox(serverSettings, clientSettings))
	return settingsTab
}

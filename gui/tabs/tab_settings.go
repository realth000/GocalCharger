package tabs

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

func makeServerSettingsArea() fyne.CanvasObject {
	return widget.NewCard("Server", "Server settings", nil)
}

func makeClientSettingsArea() fyne.CanvasObject {
	return widget.NewCard("Client", "Client settings", nil)
}

func NewSettingsTab() *container.TabItem {
	serverSettings := makeServerSettingsArea()
	clientSettings := makeClientSettingsArea()
	settingsTab := container.NewTabItem("Settins", container.NewVBox(serverSettings, clientSettings))
	return settingsTab
}

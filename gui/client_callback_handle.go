package main

import (
	"fmt"
	"fyne.io/fyne/v2/dialog"
	"gocalcharger/gui/action"
)

func handleClientSayHelloSuccess(callback action.ClientSayHelloCallbackArgs) {
	dialog.ShowInformation("Test connect", fmt.Sprintf("Test connect to server successful\n[Server: %s]", callback.ServerTarget), mainWindow)
}

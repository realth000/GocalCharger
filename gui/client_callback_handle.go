package main

import (
	"errors"
	"fmt"
	"fyne.io/fyne/v2/dialog"
	"gocalcharger/gui/action"
	"strings"
)

func handleClientSayHelloSuccess(callback action.ClientSayHelloCallbackArgs) {
	dialog.ShowInformation("Test connect", fmt.Sprintf("Test connect to server successful\n[Server: %s]", callback.ServerTarget), mainWindow)
}

func handleClientSayHelloFailed(callback action.ClientSayHelloCallbackArgs) {
	wrapString := strings.ReplaceAll(callback.Error.Error(), "error:", "error:\n")
	errString := fmt.Sprintf("%s\n[Server: %s]", wrapString, callback.ServerTarget)
	dialog.ShowError(errors.New(errString), mainWindow)
}

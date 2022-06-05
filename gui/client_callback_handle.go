package main

import (
	"errors"
	"fmt"
	"fyne.io/fyne/v2/dialog"
	"gocalcharger/gui/action"
	"gocalcharger/gui/tabs"
	"log"
	"strings"
)

func handleClientSayHelloSuccess(callback action.ClientSayHelloCallbackArgs) {
	dialog.ShowInformation("Test connect", fmt.Sprintf("Test connect to server successful\n[Server: %s]", callback.ServerTarget), mainWindow)
}

func handleClientSayHelloFailed(callback action.ClientSayHelloCallbackArgs) {
	wrapString := strings.ReplaceAll(callback.Error.Error(), "error:", "error:\n")
	errString := fmt.Sprintf("%s\n[Server: %s]", wrapString, callback.ServerTarget)
	dialog.ShowError(errors.New(errString), mainWindow)
	tabs.UpdateServerStatus(tabs.ServerClosed)
}

func handleServerStartGRPCSuccess(callback action.ServerStartGRPCArgs) {
	log.Printf("Start gRPC server successful[ServeTarget=%s]", callback.ServeTarget)
	tabs.UpdateServerStatus(tabs.ServerStarted)
}

func handleServerStartGRPCFailed(callback action.ServerStartGRPCArgs) {
	wrapString := strings.ReplaceAll(callback.Error.Error(), "error:", "error:\n")
	errString := fmt.Sprintf("%s\n[ServeTarget=%s]", wrapString, callback.ServeTarget)
	dialog.ShowError(errors.New(errString), mainWindow)
}

func handleServerStopGRPCSuccess() {
	tabs.UpdateServerStatus(tabs.ServerClosed)
}

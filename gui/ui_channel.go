package main

import (
	"gocalcharger/gui/action"
	"gocalcharger/gui/client"
	"gocalcharger/gui/server"
	"gocalcharger/gui/tabs"
)

// Channels
var (
	clientChannel         = &client.Channel
	clientCallbackChannel = &client.CallbackChannel
	serverCallbackChannel = &server.CallbackChannel
	uiTabsChannel         = &tabs.UITabsChannel
)

func StartReceivingChannels() {
	for {
		select {
		case x := <-*clientCallbackChannel:
			switch x.CallbackName {
			case action.ClientSayHelloSuccess:
				go handleClientSayHelloSuccess(x.CallbackArgs.(action.ClientSayHelloCallbackArgs))
			case action.ClientSayHelloFailed:
				go handleClientSayHelloFailed(x.CallbackArgs.(action.ClientSayHelloCallbackArgs))
			case action.ClientDownloadUpdate:
				go handleClientDownloadUpdate(x.CallbackArgs.(action.ClientDownloadUpdateArgs))
			}
		case x := <-*serverCallbackChannel:
			switch x.CallbackName {
			case action.ServerStartGRPCSuccess:
				go handleServerStartGRPCSuccess(x.CallbackArgs.(action.ServerStartGRPCArgs))
			case action.ServerStartGRPCFailed:
				go handleServerStartGRPCFailed(x.CallbackArgs.(action.ServerStartGRPCArgs))
			case action.ServerStopGRPCSuccess:
				go handleServerStopGRPCSuccess()
			}
		case x := <-*uiTabsChannel:
			switch x.ActionName {
			case action.UIStartServer:
				go server.StartServer()
			case action.UIStopServer:
				go server.StopServer()
			case action.ClientSayHello:
				go handleClientSayHello(x.ActionArgs.(action.ClientSayHelloArgs))
			case action.UIDownloadFile:
				go handleClientDownloadFile(x.ActionArgs.(action.UIDownloadFileArgs))
			}
		}
	}
}

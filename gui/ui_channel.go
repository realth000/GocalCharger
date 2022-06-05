package main

import (
	"gocalcharger/gui/action"
	"gocalcharger/gui/client"
	"gocalcharger/gui/server"
	"gocalcharger/gui/tabs"
)

// Channels
var (
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
			}

		}
	}
}

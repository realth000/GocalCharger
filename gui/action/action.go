package action

// Client Actions

type ClientActionName = string

const (
	ClientSayHello     ClientActionName = "ClientSayHello"
	ClientDownloadFile ClientActionName = "ClientDownloadFile"
)

type ClientAction struct {
	ActionName ClientActionName
	ActionArgs interface{}
}

type ClientSayHelloArgs struct {
	ClientName string
}

type ClientDownloadFileArgs struct {
	FilePath string
}

// Client callback Actions

type ClientActionCallbackName = string

const (
	ClientSayHelloSuccess ClientActionCallbackName = "SayHelloSuccess"
	ClientSayHelloFailed  ClientActionCallbackName = "SayHelloFailed"
)

type ClientActionCallback struct {
	CallbackName ClientActionCallbackName
	CallbackArgs interface{}
}

type ClientSayHelloCallbackArgs struct {
	ServerTarget string
	Error        error
}

// Server callback actions

type ServerActionCallbackName = string

const (
	ServerStartGRPCSuccess ServerActionCallbackName = "StartGRPCSuccess"
	ServerStartGRPCFailed  ServerActionCallbackName = "StartGRPCFailed"
	ServerStopGRPCSuccess  ServerActionCallbackName = "StopGRPCSuccess"
)

type ServerActionCallback struct {
	CallbackName ServerActionCallbackName
	CallbackArgs interface{}
}

type ServerStartGRPCArgs struct {
	ServeTarget string
	Error       error
}

// UI actions

type UIActionName = string

const (
	UIStartServer  UIActionName = "UIStartServer"
	UIStopServer   UIActionName = "UIStopServer"
	UIDownloadFile UIActionName = "UIDownloadFile"
)

type UIAction struct {
	ActionName UIActionName
	ActionArgs interface{}
}

type UIDownloadFileArgs struct {
	FilePath string
}

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
	ID       int
	FilePath string
}

// Client callback Actions

type ClientActionCallbackName = string

const (
	ClientSayHelloSuccess ClientActionCallbackName = "SayHelloSuccess"
	ClientSayHelloFailed  ClientActionCallbackName = "SayHelloFailed"
	ClientDownloadUpdate  ClientActionCallbackName = "ClientDownloadProgressUpdate"
)

type ClientActionCallback struct {
	CallbackName ClientActionCallbackName
	CallbackArgs interface{}
}

type ClientSayHelloCallbackArgs struct {
	ServerTarget string
	Error        error
}

type ClientDownloadUpdateArgs struct {
	ID        int
	FilePath  string
	Size      int
	TotalSize int
	Finished  bool
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
	UIStartServer        UIActionName = "UIStartServer"
	UIStopServer         UIActionName = "UIStopServer"
	UIDownloadFile       UIActionName = "UIDownloadFile"
	UIDownloadFileUpdate UIActionName = "UIDownloadFileUpdate"
)

type UIAction struct {
	ActionName UIActionName
	ActionArgs interface{}
}

type UIDownloadFileArgs struct {
	ID       int
	FilePath string
}

type UIDownloadFileUpdateArgs struct {
	ID        int
	FilePath  string
	Size      int
	TotalSize int
	Finished  bool
}

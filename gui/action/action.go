package action

// Client Actions

type ClientActionName = string

const (
	ClientSayHello     ClientActionName = "SayHello"
	ClientDownloadFile ClientActionName = "SayHello"
)

type ClientAction struct {
	ActionName ClientActionName
	ActionArgs interface{}
}

type ClientSayHelloArgs struct {
	ClientName string
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
)

type ServerActionCallback struct {
	CallbackName ServerActionCallbackName
	CallbackArgs interface{}
}

type ServerStartGRPCArgs struct {
	ServeTarget string
	Error       error
}

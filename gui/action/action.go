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
	ClientSayHelloTimeOut ClientActionCallbackName = "SayHelloTimeout"
)

type ClientActionCallback struct {
	CallbackName ClientActionCallbackName
	CallbackArgs interface{}
}

type ClientSayHelloCallbackArgs struct {
	ServerTarget string
}

package action

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

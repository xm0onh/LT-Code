package ConActAndNetInterface

type NetworkToConActInterface interface {
	//	GetIndex(string, []string) int
	//	SendDroplets(E.Request, int, string, *ConActor.ConActor)
	PassMsgToActor(event interface{}, committeeSize int, sourceIp string)
}

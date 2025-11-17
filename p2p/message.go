package p2p

const (
	IncomingMessage = 0x1
	IncomingStream = 0x2
)
// Message hold any aribitrary data that is being sent 
// each transport bw node and the server
type RPC struct{
	From 	string
	Payload []byte
	Stream	bool
}
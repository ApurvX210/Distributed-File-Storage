package p2p

// Message hold any aribitrary data that is being sent 
// each transport bw node and the server
type RPC struct{
	From 	Peer
	Payload []byte
}
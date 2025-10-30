package p2p

// Peers is the interface that represent the remote node
type Peer interface{
	Close() error
}

// Handle Communication between node in the network
// This can be of the form Tcp/Udp/WebSocket
type Transport interface{
	ListenAndAccept() error
	Consume() <- chan RPC
}
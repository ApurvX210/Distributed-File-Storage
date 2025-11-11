package p2p

import "net"

// Peers is the interface that represent the remote node
type Peer interface{
	Send([]byte) error
	RemoteAddr() 	net.Addr
	Close() 		error
}

// Handle Communication between node in the network
// This can be of the form Tcp/Udp/WebSocket
type Transport interface{
	ListenAndAccept() error
	Consume() <- chan RPC
	Close() error
	Dial(addr string) error
}
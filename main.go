package main

import (
	"Distributed-File-Storage/p2p"
	"log"
)

func main() {
	tcpOpts := p2p.TCPTransportOpts{
		ListenAddress: "localhost:5001",
		ShakeHand: 	  p2p.TCPHandShake,
		Decoder:	  p2p.DefaultDecoder{},
	}
	
	tcp := p2p.NewTcpTransport(tcpOpts)
	log.Fatal(tcp.ListenAndAccept())

	select{}
}
package main

import (
	"Distributed-File-Storage/p2p"
	"log"
)

func main() {
	tcpOpts := p2p.TCPTransportOpts{
		ListenAddress: ":5001",
		ShakeHand: 	  p2p.TCPHandShake,
		Decoder:	  p2p.DefaultDecoder{},
		// To Do: OnPeerFunc
	}

	storeOpts := StoreOpts{
		Root: "5001_DFA",
		PathTranformerFunc: CASPathTransformer,
	}

	fsOpts := FileServerOpts{
			StoreOpts: storeOpts,
			Transport: p2p.NewTcpTransport(tcpOpts),
		}
	fs := NewFileServer(fsOpts)
	
	log.Fatal(fs.Start())

	select{}
}
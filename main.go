package main

import (
	"Distributed-File-Storage/p2p"
	"log"
)

func makeServer(listenAddress string,nodes ...string) *FileServer{
	tcpOpts := p2p.TCPTransportOpts{
		ListenAddress: listenAddress,
		ShakeHand: 	  p2p.TCPHandShake,
		Decoder:	  p2p.DefaultDecoder{},
		// To Do: OnPeerFunc
	}

	storeOpts := StoreOpts{
		Root: listenAddress+"DFA",
		PathTranformerFunc: CASPathTransformer,
	}

	fsOpts := FileServerOpts{
			StoreOpts: storeOpts,
			Transport: p2p.NewTcpTransport(tcpOpts),
			quitch: 		make(chan struct{}),
			BootStrapNodes:	nodes,
		}
	return NewFileServer(fsOpts)
	
}

func main() {
	fs1 := makeServer(":5001")
	fs2 := makeServer(":4001",":5001")
	go func ()  {
		log.Fatal(fs1.Start())
	}()
	log.Fatal(fs2.Start())
	select{}
}
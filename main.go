package main

import (
	"Distributed-File-Storage/p2p"
	"bytes"
	"fmt"
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

	tcpTransport := p2p.NewTcpTransport(tcpOpts)

	fsOpts := FileServerOpts{
			StoreOpts: storeOpts,
			Transport: tcpTransport,
			quitch: 		make(chan struct{}),
			BootStrapNodes:	nodes,
		}

	fs := NewFileServer(fsOpts)

	tcpTransport.OnPeer = fs.OnPeer

	return fs
	
}

func main() {
	fs1 := makeServer(":5001")
	fs2 := makeServer(":4001",":5001")
	go func ()  {
		log.Fatal(fs1.Start())
	}()
	log.Fatal(fs2.Start())
	
	data := bytes.NewReader([]byte("File Stored Here"))
	fmt.Println(data)
	fs2.StoreFile("My Data",data)
}
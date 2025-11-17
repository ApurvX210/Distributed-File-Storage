package main

import (
	"Distributed-File-Storage/p2p"
	// "bytes"
	"fmt"
	"io"

	// "fmt"
	"log"
	"strings"
	"time"
)

func makeServer(listenAddress string,nodes ...string) *FileServer{
	tcpOpts := p2p.TCPTransportOpts{
		ListenAddress: listenAddress,
		ShakeHand: 	  p2p.TCPHandShake,
		Decoder:	  p2p.DefaultDecoder{},
		// To Do: OnPeerFunc
	}

	storeOpts := StoreOpts{
		Root: strings.Split(listenAddress, ":")[1]+"_"+"DFA",
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
	time.Sleep(time.Second)
	go func ()  {
		log.Fatal(fs2.Start())
	}()
	
	time.Sleep(time.Second * 2)
	// data := bytes.NewReader([]byte("File Stored hello my name is apurv lorem dcbjsdvbds sdcbsdcjsd sdbcsdbkc sbcjkdsck kjsbdckjdbskjcbsdcbkdsbckdskc sdjbckjdsbckjdbckdsb"))
	// fs2.StoreData("My Data",data)

	r, err := fs2.Get("My Data")
	if err != nil{
		log.Fatal(err)
	}
	b, err := io.ReadAll(r)
	fmt.Println(string(b))
	select{}
}
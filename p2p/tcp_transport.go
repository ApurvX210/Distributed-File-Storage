package p2p

import (
	"log"
	"net"
	"sync"
)

type TCPTransport struct {
	listenAddress string
	listener      net.Listener
	mu			  sync.RWMutex
	peers         map[net.Addr]Peer
}


func NewTcp(listenAddress string) *TCPTransport{
	return &TCPTransport{
		listenAddress: listenAddress,
	}
}

func handleConnection(conn net.Conn){

}

func (tcp *TCPTransport) listenAndAccept() error{
	var err error
	tcp.listener,err = net.Listen("tcp",tcp.listenAddress)

	if err != nil{
		log.Fatal("Error Occured while Initializing listener",err)
		return err
	}
	
	return tcp.acceptRequests()
}

func (tcp *TCPTransport) acceptRequests() error{
	for{
		con,err := tcp.listener.Accept()
		if err != nil{
			log.Fatal("Error Occured while Accepting Request",err)
			return err
		}
		go handleConnection(con)
	}
}
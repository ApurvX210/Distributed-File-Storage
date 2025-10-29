package p2p

import (
	"fmt"
	"log"
	"log/slog"
	"net"
	"sync"
)

// It represent the remote user connected via Tcp protocol
type TCPPeer struct{
	// Underlying connection of the peer
	conn net.Conn

	// If we dial and a connection => outbound - true
	// But if we accept and retrieve a connection => outbound - false
	outbound bool
}

func NewTCPPeer(conn net.Conn, outbound bool) *TCPPeer{
	return &TCPPeer{
		conn: conn,
		outbound: outbound,
	}
}

type TCPTransportOpts struct{
	ListenAddress string
	ShakeHand 	  HandShakerFunc
	Decoder		  Decoder
}
type TCPTransport struct {
	TCPTransportOpts
	listener      net.Listener
	mu			  sync.RWMutex
	peers         map[net.Addr]Peer
}


func NewTcpTransport(opts TCPTransportOpts) *TCPTransport{
	return &TCPTransport{
		TCPTransportOpts: opts,
	}
}

func (tcp *TCPTransport) ListenAndAccept() error{
	var err error
	tcp.listener,err = net.Listen("tcp",tcp.ListenAddress)

	if err != nil{
		log.Fatal("Error Occured while Initializing listener ",err)
		return err
	}
	slog.Info("Accepting Tcp connection on ","Address",tcp.ListenAddress)
	return tcp.acceptRequests()
}

func (tcp *TCPTransport) acceptRequests() error{
	for{
		conn,err := tcp.listener.Accept()
		if err != nil{
			log.Fatal("Error Occured while Accepting Request ",err)
			return err
		}
		go tcp.handleConnection(conn)
	}
}

func (tcp *TCPTransport) handleConnection(conn net.Conn){
	peer := NewTCPPeer(conn,false)
	
	if err:= tcp.ShakeHand(peer); err!=nil{
		slog.Error("Error occured while handshake with connection","Conn",conn)
		peer.conn.Close()
		return
	}
	fmt.Printf("New Incoming Connection %+v\n",peer)
	// Read Loop
	msg := &Message{}
	for{
		fmt.Println("hello")
		if err := tcp.Decoder.Decode(conn,msg); err != nil{
			fmt.Println("Error")
			slog.Error("Error occured while Reading the connection","Error",err)
			continue
		}

		fmt.Printf("%+v\n",msg)
	}
	
}
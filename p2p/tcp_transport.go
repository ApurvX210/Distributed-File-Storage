package p2p

import (
	"fmt"
	"io"
	"log"
	"log/slog"
	"net"
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

func (peer *TCPPeer) Close() error{
	return peer.conn.Close()
}

type TCPTransportOpts struct{
	ListenAddress string
	ShakeHand 	  HandShakerFunc
	Decoder		  Decoder
	OnPeer		  func(Peer) error
}
type TCPTransport struct {
	TCPTransportOpts
	listener      net.Listener
	rpcChan		  chan RPC
	// mu			  sync.RWMutex
	// peers         map[net.Addr]Peer
}


func NewTcpTransport(opts TCPTransportOpts) *TCPTransport{
	return &TCPTransport{
		TCPTransportOpts: opts,
		rpcChan: make(chan RPC),
	}
}

// func ()

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
	var err error

	peer := NewTCPPeer(conn,false)
	defer func(){
		peer.Close()
		fmt.Printf("Dropping peer connection : %s",err)
	}()
	
	if err = tcp.ShakeHand(peer); err!=nil{
		slog.Error("Error occured while handshake with connection","Conn",conn)
		peer.conn.Close()
		return
	}

	if tcp.OnPeer != nil{
		if err = tcp.OnPeer(peer); err != nil{
			return
		}
	}

	fmt.Printf("New Incoming Connection %+v\n",peer)
	// Read Loop
	rpc := &RPC{}
	for{
		fmt.Println("hello")
		if err = tcp.Decoder.Decode(conn,rpc); err != nil{
			if err == io.EOF{
				return
			}else{
				slog.Error("Error occured while Reading the connection","Error",err)
				continue
			}
		}
		rpc.From = peer
		fmt.Printf("%+v\n",rpc)
	}
}
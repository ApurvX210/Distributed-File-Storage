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
	net.Conn
	// If we dial and a connection => outbound - true
	// But if we accept and retrieve a connection => outbound - false
	outbound bool
}

func NewTCPPeer(conn net.Conn, outbound bool) *TCPPeer{
	return &TCPPeer{
		conn,
		outbound,
	}
}

func (peer *TCPPeer) Close() error{
	return peer.Close()
}

func (peer *TCPPeer) Send(data []byte) error{
	_,err := peer.Write(data)
	return err
}

// Remote address implement the peer interface
// Returns the Remote address of underlying connection
func (peer *TCPPeer) RemoteAddr() net.Addr{
	return peer.RemoteAddr()
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

func (tcp *TCPTransport) Consume() <- chan RPC{
	return tcp.rpcChan
}

// Dial Implements the Transport interface
func (tcp *TCPTransport) Dial(addr string) error{
	conn,err := net.Dial("tcp",addr)
	if err != nil{
		return err
	}
	go tcp.handleConnection(conn,true)
	return nil
}

func (tcp *TCPTransport) ListenAndAccept() error{
	var err error
	tcp.listener,err = net.Listen("tcp",tcp.ListenAddress)

	if err != nil{
		log.Fatal("Error Occured while Initializing listener ",err)
		return err
	}
	slog.Info("Accepting Tcp connection on ","Address",tcp.ListenAddress)
	go tcp.acceptRequests()
	return nil
}

func (tcp *TCPTransport) acceptRequests() error{
	for{
		conn,err := tcp.listener.Accept()
		if err != nil{
			log.Fatal("Error Occured while Accepting Request ",err)
			return err
		}
		go tcp.handleConnection(conn,false)
	}
}

// Close implements the Transport Interface CLose function
func (tcp *TCPTransport) Close() error{
	return tcp.listener.Close()
}

func (tcp *TCPTransport) handleConnection(conn net.Conn,outbound bool){
	var err error

	peer := NewTCPPeer(conn,outbound)
	defer func(){
		peer.Close()
		fmt.Printf("Dropping peer connection : %s",err)
	}()
	
	if err = tcp.ShakeHand(peer); err!=nil{
		slog.Error("Error occured while handshake with connection","Conn",conn)
		peer.Close()
		return
	}
	fmt.Println(tcp.OnPeer)
	if tcp.OnPeer != nil{
		if err = tcp.OnPeer(peer); err != nil{
			return
		}
	}
	if outbound{
		fmt.Printf("New Outgoing Connection %+v\n",peer)
	}else{
		fmt.Printf("New Incoming Connection %+v\n",peer)
	}
	
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
		tcp.rpcChan <- *rpc
		fmt.Printf("%+v\n",rpc)
	}
}
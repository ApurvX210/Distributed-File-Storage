package p2p

import (
	"testing"
	"github.com/stretchr/testify/assert"
)

func TestTcpConn(t *testing.T) {
	tcpOpts := TCPTransportOpts{
		ListenAddress: "localhost:5001",
		ShakeHand: 	  TCPHandShake,
		Decoder:	  DefaultDecoder{},
	}
	
	tcp := NewTcpTransport(tcpOpts)

	assert.Equal(t,tcp.ListenAddress,tcpOpts.ListenAddress)

	assert.Nil(t,tcp.ListenAndAccept())
}

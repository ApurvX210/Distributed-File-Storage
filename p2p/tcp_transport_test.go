package p2p

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTcpConn(t *testing.T) {
	listenAddress := "localhost:5001"
	tr := NewTcpTransport(listenAddress)

	assert.Equal(t,tr.listenAddress,listenAddress)

	assert.Nil(t,tr.ListenAndAccept())
}

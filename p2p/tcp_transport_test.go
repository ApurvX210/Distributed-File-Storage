package p2p

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTcpConn(t *testing.T) {
	listenAddress := "localhost:5001"
	tr := NewTcp(listenAddress)

	assert.Equal(t,tr.listenAddress,listenAddress)
}

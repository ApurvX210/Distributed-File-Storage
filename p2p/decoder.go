package p2p

import (
	"encoding/gob"
	"io"
)

type Decoder interface {
	Decode(io.Reader,*RPC) error
}

type GobDecoder struct{}

func (dec GobDecoder) Decode(r io.Reader, msg *RPC) error{
	return gob.NewDecoder(r).Decode(msg)
}

type DefaultDecoder struct{}

func (dec DefaultDecoder) Decode(r io.Reader, msg *RPC) error{
	peekBuf := make([]byte,1)

	_,err := r.Read(peekBuf)
	if err != nil{
		return err
	}

	// Incase of a stream we are not decoding what is being sent over the network
	// We are just setting stream true so we can handle that in our handle message
	switch peekBuf[0] {
	case IncomingMessage:
		msg.Stream = false
	case IncomingStream:
		msg.Stream = true
		return nil
	}

	buf := make([]byte,1028)
	n,err := r.Read(buf)

	if err != nil{
		return err
	}
	msg.Payload = buf[:n]
	return nil
}

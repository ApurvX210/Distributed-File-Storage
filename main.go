package main

import (
	"Distributed-File-Storage/p2p"
	"log"
)

func main() {
	listenAddr := "localhost:5001"
	tcp := p2p.NewTcpTransport(listenAddr)
	log.Fatal(tcp.ListenAndAccept())

	select{}
}
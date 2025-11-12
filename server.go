package main

import (
	"Distributed-File-Storage/p2p"
	"bytes"
	"encoding/gob"
	"fmt"
	"io"
	"log"
	"sync"
	// "io"
)

type FileServerOpts struct {
	StoreOpts     	StoreOpts
	Transport    	p2p.Transport
	quitch 			chan struct{}
	BootStrapNodes	[]string
}

type FileServer struct {
	FileServerOpts
	peerLock	sync.RWMutex
	peers		map[string] p2p.Peer
	store 		*Store
}

type Payload struct{
	key 	string
	Data 	[]byte
}

func NewFileServer(opts FileServerOpts) *FileServer {
	return &FileServer{
		FileServerOpts: opts,
		// peerLock: 		sync.RWMutex{},
		peers: 			make(map[string]p2p.Peer),
		store:          NewStore(opts.StoreOpts),
	}
}

func (fs *FileServer) Start() error{
	err := fs.Transport.ListenAndAccept()
	if err != nil{
		return err
	}
	fs.bootStrapNetwork()
	fs.loop()
	return nil
}

func (fs *FileServer) OnPeer(peer p2p.Peer) error{
	fs.peerLock.Lock()
	defer fs.peerLock.Unlock()
	fs.peers[peer.RemoteAddr().String()] = peer
	fmt.Printf("Coonnected with Remote %s",peer.RemoteAddr())
	return nil
}



func (fs *FileServer) Broadcast(p Payload) error{
	peers := []io.Writer{}
	for _,peer := range fs.peers{
		peers = append(peers, peer)
	}

	mu := io.MultiWriter(peers...)
	err := gob.NewEncoder(mu).Encode(p)
		if err != nil{
			return err
		}

	return nil
}

func (fs *FileServer) StoreFile(key string,file io.Reader) error{
	// Store this file to the disk
	// Broadcast this file to all known peer in the network
	err := fs.store.Write(key,file)
	if err != nil{
		return err
	}

	buf := bytes.Bufer
	p := Payload{
		key: key,
		Data: fs.store.BufferRead(key),
	}

	return nil
}

func (fs *FileServer) bootStrapNetwork() error{
	for _,addr := range fs.BootStrapNodes{
		go func(addr string){
			if err := fs.Transport.Dial(addr);err != nil{
				log.Println("Dial Error Occured",err)
			}
			
		}(addr)
	}
	return nil
}

func (fs *FileServer) loop(){
	defer func ()  {
		log.Panicln("File Server Stopped")
		fs.Transport.Close()
		fs.stop()
	}()
	for{
		select{
		case msg := <- fs.Transport.Consume():
			fmt.Println(msg)
		case <-fs.quitch:
			return 
		}
	}
}

func (fs *FileServer) stop(){
	close(fs.quitch)
}

// func (fs *FileServer) Store(key string,r io.Reader) error{
// 	return fs.store.writeStream(key,r)
// }
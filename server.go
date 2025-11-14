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

type Message struct{
	Payload any
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
	fmt.Printf("Connected with Remote %s",peer.RemoteAddr())
	return nil
}

func (fs *FileServer) Broadcast(msg Message) error{
	peers := []io.Writer{}
	for _,peer := range fs.peers{
		peers = append(peers, peer)
	}
	buf := new(bytes.Buffer)
	err := gob.NewEncoder(buf).Encode(msg)
		if err != nil{
			return err
		}

	mu := io.MultiWriter(peers...)
	_,err = mu.Write(buf.Bytes())

	return err
}

func (fs *FileServer) StoreData(key string,r io.Reader) error{
	// Store this file to the disk
	// Broadcast this file to all known peer in the network


	msg := Message{
		Payload: []byte("storagekey"),
	}
	buf := new(bytes.Buffer)
	if err := gob.NewEncoder(buf).Encode(msg); err != nil{
		return err
	}

	patload := []byte("This is my Large file")
	for _,peer := range fs.peers{
		if err := peer.Send(patload); err != nil{
			return err
		}
	}
	return nil
	// buf := new(bytes.Buffer)
	// tee := io.TeeReader(r,buf)
	// err := fs.store.Write(key,tee)
	// if err != nil{
	// 	return err
	// }
	
	// _,err = io.Copy(buf,r)
	// if err != nil{
	// 	return err
	// }

	// p := DataMessage{
	// 	key: key,
	// 	Data: buf.Bytes(),
	// }

	// return fs.Broadcast(Message{
	// 	from: fs.Transport.ListenAddr(),
	// 	Payload: p,
	// })
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
		case rpc := <- fs.Transport.Consume():
			println("Hello my name is Apurv")
			var msg Message
			err := gob.NewDecoder(bytes.NewReader(rpc.Payload)).Decode(&msg)
			if err != nil{
				log.Println("Error occured while decoding RPC Channel",err)
				continue
			}
			err = fs.handleMessage(&msg)
			if err != nil{
				log.Printf("Error occured while Storing data recieved from RPC Channel %s",err)
				continue
			}
		case <-fs.quitch:
			return 
		}
	}
}

func (fs *FileServer) handleMessage(msg *Message) error{

	switch msg.Payload.(type){
	case *DataMessage:
		payload := msg.Payload
		fs.store.Write(msg.Payload.key)
	}

	return nil
}

func (fs *FileServer) stop(){
	close(fs.quitch)
}

// func (fs *FileServer) Store(key string,r io.Reader) error{
// 	return fs.store.writeStream(key,r)
// }
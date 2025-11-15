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

type MessageStoreFile struct{
	Key 	string
	Size 	int64
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

	size,err := fs.store.Write(key,r)
	if err != nil{
		return err
	}

	msg := Message{
		Payload: MessageStoreFile{
			Key: key,
			Size: size,
		},
	}
	buf := new(bytes.Buffer)
	if err := gob.NewEncoder(buf).Encode(msg); err != nil{
		return err
	}
	for _,peer := range fs.peers{
		if err := peer.Send(buf.Bytes()); err != nil{
			fmt.Println("Error occurred while broadcasting Message to store file to peer ",peer)
		}
	}
	file,err := fs.store.BufferRead(key)
	if err != nil{
		return err
	}
	for _,peer := range fs.peers{
		if err := peer.Stream(file); err != nil{
			fmt.Println("Error occurred while broadcasting file to peer ",peer)
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
			var msg Message
			err := gob.NewDecoder(bytes.NewReader(rpc.Payload)).Decode(&msg)
			if err != nil{
				log.Println("Error occured while decoding RPC Channel",err)
				continue
			}
			fmt.Printf("Recieved %+v\n",msg.Payload)

			if err := fs.handleMessage(rpc.From,&msg);err != nil{
				log.Printf("Error occured while handling message ",msg)
				continue
			}
		case <-fs.quitch:
			return 
		}
	}
}

func (fs *FileServer) handleMessage(from string,msg *Message) error{

	switch v := msg.Payload.(type){
	case MessageStoreFile:
		return fs.handleMessageStoreFile(from,v)
	}

	return nil
}

func (fs *FileServer) handleMessageStoreFile(from string,msg MessageStoreFile) error{
	peer,ok := fs.peers[from]
	if !ok{
		return fmt.Errorf("peer not registered %+v",peer)
	}
	_,err := fs.store.Write(msg.Key,io.LimitReader(peer,msg.Size))
	if err != nil{
		return fmt.Errorf("error occured while Storing data recieved from RPC Channel %s",err)
	}
	peer.(*p2p.TCPPeer).Wg.Done()
	return nil
}

func (fs *FileServer) stop(){
	close(fs.quitch)
}

// func (fs *FileServer) Store(key string,r io.Reader) error{
// 	return fs.store.writeStream(key,r)
// }

func init(){
	gob.Register(MessageStoreFile{})
}
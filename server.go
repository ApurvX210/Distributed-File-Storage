package main

import (
	"Distributed-File-Storage/p2p"
	"fmt"
	"log"
	// "io"
)

type FileServerOpts struct {
	StoreOpts     StoreOpts
	Transport     p2p.Transport
}

type FileServer struct {
	FileServerOpts
	store *Store
	quitch chan struct{}
}

func NewFileServer(opts FileServerOpts) *FileServer {
	return &FileServer{
		FileServerOpts: opts,
		store:          NewStore(opts.StoreOpts),
		quitch: 		make(chan struct{}),
	}
}

func (fso *FileServer) Start() error{
	err := fso.Transport.ListenAndAccept()
	if err != nil{
		return err
	}
	fso.loop()
	return nil
}

func (fso *FileServer) loop(){
	defer func ()  {
		log.Panicln("File Server Stopped")
		fso.Transport.Close()
		fso.stop()
	}()
	for{
		select{
		case msg := <- fso.Transport.Consume():
			fmt.Println(msg)
		case <-fso.quitch:
			return 
		}
	}
}

func (fso *FileServer) stop(){
	close(fso.quitch)
}

// func (fso *FileServer) Store(key string,r io.Reader) error{
// 	return fso.store.writeStream(key,r)
// }
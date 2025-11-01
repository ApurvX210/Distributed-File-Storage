package main

import (
	"crypto/sha1"
	"encoding/hex"
	"io"
	"log"
	"os"
	"strings"
)

type PathTranformer func (string) PathKey

var DefaultPathTransformer = func(key string) PathKey{
	return PathKey{
		PathName: key,
		FileName: key,
	}
}

func CASPathTransformer(key string) PathKey{
	hash := sha1.Sum([]byte(key))
	hashStr := hex.EncodeToString(hash[:])

	blockSize := 5
	sliceLen := len(hashStr) / blockSize

	path := make([]string,sliceLen)

	for i:=0;i<sliceLen;i++{
		path[i] = hashStr[i*blockSize:i*blockSize+blockSize]
	}

	return PathKey{
		PathName: strings.Join(path,"/"),
		FileName: hashStr,
	} 
}

type PathKey struct{
	PathName string
	FileName string
}

func (p *PathKey) GenerateFilePath() string{
	return p.PathName+"/"+p.FileName
}

type StoreOpts struct {
	PathTranformerFunc PathTranformer
}

type Store struct {
	StoreOpts
}

func NewStore(storeOpts StoreOpts) *Store{
	return &Store{
		StoreOpts: storeOpts,
	}
}

func (s *Store) readStream(key string) (io.Reader,error){
	pathKey := s.PathTranformerFunc(key)
	filePath := pathKey.GenerateFilePath()

	f,err := os.Open(filePath)
	if err != nil{
		return nil,err
	}

	return f,nil
}

func (s *Store) Read(key string) (io.Reader,error){
	f,err := s.readStream(key)
	if err != nil{
		return nil,err
	}

}

func (s *Store) writeStream(key string, r io.Reader) error{
	pathKey := s.PathTranformerFunc(key)

	if err := os.MkdirAll(pathKey.PathName,os.ModePerm); err != nil{
		return err
	}

	filePath := pathKey.GenerateFilePath()

	f,err := os.Create(filePath)
	if err != nil{
		return err
	}

	n,err := io.Copy(f,r)
	if err != nil{
		return err
	}
	log.Printf("Written {%d} bytes to disk: %s",n,filePath)

	return  nil
}
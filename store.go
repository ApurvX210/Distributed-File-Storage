package main

import (
	"crypto/sha1"
	"encoding/hex"
	"io"
	"log"
	"os"
	"strings"
)

type PathTranformer func (string) string

var DefaultPathTransformer = func(key string)string{
	return key
}

func CASPathTransformer(key string) string{
	hash := sha1.Sum([]byte(key))
	hashStr := hex.EncodeToString(hash[:])

	blockSize := 5
	sliceLen := len(hashStr) / blockSize

	path := make([]string,sliceLen)

	for i:=0;i<sliceLen;i++{
		path[i] = hashStr[i*blockSize:i*blockSize+blockSize]
	}
	return strings.Join(path,"/")
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

func (s *Store) writeStream(key string, r io.Reader) error{
	pathName := s.PathTranformerFunc(key)

	if err := os.MkdirAll(pathName,os.ModePerm); err != nil{
		return err
	}

	fileName := "file.txt"
	pathFileName := pathName+"/"+fileName

	f,err := os.Create(pathFileName)
	if err != nil{
		return err
	}

	n,err := io.Copy(f,r)
	if err != nil{
		return err
	}
	log.Printf("Written {%d} bytes to disk: %s",n,pathFileName)

	return  nil
}
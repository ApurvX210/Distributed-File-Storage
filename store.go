package main

import (
	"bytes"
	"crypto/sha1"
	"encoding/hex"
	"io"
	"log"
	"os"
	"strings"
)

const DeaultRootFolder = "DFS"

type PathTranformer func(string, string) PathKey

var DefaultPathTransformer = func(key string, root string) PathKey {
	return PathKey{
		PathName: root + "/" + key,
		FileName: key,
	}
}

func CASPathTransformer(key string, root string) PathKey {
	hash := sha1.Sum([]byte(key))
	hashStr := hex.EncodeToString(hash[:])

	blockSize := 5
	sliceLen := len(hashStr) / blockSize

	path := make([]string, sliceLen)

	for i := 0; i < sliceLen; i++ {
		path[i] = hashStr[i*blockSize : i*blockSize+blockSize]
	}

	return PathKey{
		PathName: root + "/" + strings.Join(path, "/"),
		FileName: hashStr,
	}
}

type PathKey struct {
	PathName string
	FileName string
}

func (p PathKey) firstPathName() string {
	path := strings.Split(p.PathName, "/")
	return path[0] + "/" + path[1]
}

func (p *PathKey) GenerateFilePath() string {
	return p.PathName + "/" + p.FileName
}

type StoreOpts struct {
	// Root is the folder name in which all data will be stored
	Root               string
	PathTranformerFunc PathTranformer
}

type Store struct {
	StoreOpts
}

func NewStore(storeOpts StoreOpts) *Store {
	if storeOpts.PathTranformerFunc == nil {
		storeOpts.PathTranformerFunc = DefaultPathTransformer
	}
	if len(storeOpts.Root) == 0 {
		storeOpts.Root = DeaultRootFolder
	}
	return &Store{
		StoreOpts: storeOpts,
	}
}

func (s *Store) Delete(key string) error {
	pathKey := s.PathTranformerFunc(key, s.Root)

	defer func() {
		log.Printf("Deleted file [%s] from storage", pathKey.FileName)
	}()
	return os.RemoveAll(pathKey.firstPathName())
}

func (s *Store) Clear(key string) error {
	return os.RemoveAll(s.Root)
}

func (s *Store) Has(key string) bool {
	pathKey := s.PathTranformerFunc(key, s.Root)

	_, err := os.Stat(pathKey.GenerateFilePath())

	if err == nil {
		return true
	}
	return false
}

func (s *Store) readStream(key string) (io.ReadCloser, error) {
	pathKey := s.PathTranformerFunc(key, s.Root)
	filePath := pathKey.GenerateFilePath()

	f, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}

	return f, nil
}

func (s *Store) StreamRead(key string) (*bytes.Buffer, error) {
	f, err := s.readStream(key)

	if err != nil {
		return nil, err
	}
	defer f.Close()
	buff := new(bytes.Buffer)
	_, err = io.Copy(buff, f)

	return buff, err
}

func (s *Store) BufferRead(key string) (io.ReadCloser, error) {
	return s.readStream(key)
}

func (s *Store) Write(key string, r io.Reader) error {
	return s.writeStream(key, r)
}

func (s *Store) writeStream(key string, r io.Reader) error {
	pathKey := s.PathTranformerFunc(key, s.Root)
	if err := os.MkdirAll(pathKey.PathName, os.ModePerm); err != nil {
		return err
	}
	filePath := pathKey.GenerateFilePath()

	f, err := os.Create(filePath)
	if err != nil {
		return err
	}
	n, err := io.Copy(f, r)
	if err != nil {
		return err
	}
	log.Printf("Written {%d} bytes to disk: %s", n, filePath)

	return nil
}

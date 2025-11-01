package main

import (
	"bytes"
	"fmt"
	"testing"
	// "github.com/stretchr/testify/assert"
)

func TestPathTransformer(t *testing.T){
	key := "Yash"
	path := CASPathTransformer(key)
	fmt.Println(path)
}
func TestStore(t *testing.T){
	opts := StoreOpts{
		PathTranformerFunc: CASPathTransformer,
	}

	store := NewStore(opts)

	data := bytes.NewReader([]byte("Hello my name is Apurv"))

	err := store.writeStream("specialPicture",data)
	fmt.Println(err)
	// assert.Equal(t,store.writeStream("specialPicture",data),nil)
}
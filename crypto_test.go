package main

import (
	"bytes"
	"fmt"
	"log"
	"testing"
)

func TestCrypto(t *testing.T) {
	src := bytes.NewReader([]byte("Hello my name is Apurv"))
	dst := new(bytes.Buffer)
	key := newEncrptionKey()
	_, err := copyEncrypt(key,src,dst)
	if err != nil{
		log.Fatal(err)
	}
	
	out := new(bytes.Buffer)
	if _,err := copyDcrypt(key,dst,out); err != nil{
		log.Fatal(err)
	}
	fmt.Println(string(out.Bytes()))
}
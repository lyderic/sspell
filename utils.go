package main

import (
	"crypto/md5"
	"encoding/hex"
	"io/ioutil"
	"log"
)

func md5sum(filepath string) string {
	content, err := ioutil.ReadFile(filepath)
	if err != nil {
		log.Fatal(err)
	}
	sum16Byte := md5.Sum(content)
	return hex.EncodeToString(sum16Byte[:])
}

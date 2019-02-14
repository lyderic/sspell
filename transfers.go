package main

import (
	"fmt"
	"log"
	"path/filepath"

	minio "github.com/minio/minio-go"
)

const (
	TO   = 1
	FROM = 0
)

func transfer(direction int, file string, minioClient *minio.Client) {
	var n int64
	var err error
	switch direction {
	case TO:
		n, err = minioClient.FPutObject(
			config.Bucket,
			file,
			filepath.Join(config.SpellDir, file),
			minio.PutObjectOptions{})
	case FROM:
		err = minioClient.FGetObject(
			config.Bucket,
			file,
			filepath.Join(config.SpellDir, file),
			minio.GetObjectOptions{})
	default:
		log.Fatal("Invalid transfer direction")
	}
	if err == nil {
		fmt.Printf("> %-16.16s : ok (%d bytes)\n", file, n)
	} else {
		log.Fatalln(file, ": failed!", err)
	}
}

func getFromRemote(minioClient *minio.Client, files []string) {
	fmt.Println("Getting from remote:")
	for _, file := range files {
		transfer(FROM, file, minioClient)
	}
}

func putToRemote(minioClient *minio.Client, files []string) {
	fmt.Println("Pushing to remote:")
	for _, file := range files {
		transfer(TO, file, minioClient)
	}
}

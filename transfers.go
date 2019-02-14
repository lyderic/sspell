package main

import (
	"fmt"
	"log"
	"path/filepath"

	minio "github.com/minio/minio-go"
)

func getFromRemote(minioClient *minio.Client, files []string) {
	fmt.Println("Getting from remote:")
	for _, file := range files {
		err := minioClient.FGetObject(
			config.Bucket,
			file,
			filepath.Join(config.SpellDir, file),
			minio.GetObjectOptions{})
		if err == nil {
			fmt.Printf("> %-16.16s : ok\n", file)
		} else {
			log.Fatalln(file, ": failed!", err)
		}
	}
}

func putToRemote(minioClient *minio.Client, files []string) {
	fmt.Println("Pushing to remote:")
	for _, file := range files {
		n, err := minioClient.FPutObject(
			config.Bucket,
			file,
			filepath.Join(config.SpellDir, file),
			minio.PutObjectOptions{})
		if err == nil {
			fmt.Printf("> %-16.16s : ok (%d bytes)\n", file, n)
		} else {
			log.Fatalln(file, ": failed!", err)
		}
	}
}

package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"testing"

	minio "github.com/minio/minio-go"
)

func TestTransfer(t *testing.T) {
	fmt.Println("Removing remote files")
	var err error
	err = rclone("delete", "-v", "minio:sspell/fr.utf-8.add")
	if err != nil {
		panic(err)
	}
	err = rclone("delete", "-v", "minio:sspell/fr.utf-8.add.spl")
	if err != nil {
		panic(err)
	}
	configPath := filepath.Join(os.Getenv("HOME"), "mondes", "sspell.conf")
	configure(configPath)
	minioClient, err := minio.New(config.Endpoint, config.KeyId, config.KeySecret, config.UseSSL)
	if err != nil {
		panic(err)
	}
	transfer(TO, config.UTF8File, minioClient)
	//transfer(TO, config.BinaryFile, minioClient)
	var n int64
	n, err = minioClient.FPutObject(config.Bucket, "fr.utf-8.add.spl", filepath.Join(config.SpellDir, "fr.utf-8.add.spl"), minio.PutObjectOptions{ContentType: "application/octet-stream"})
	if err != nil {
		panic(err)
	}
	fmt.Println("TRANSFERRED", n, "bytes")
}

func rclone(args ...string) error {
	cmd := exec.Command("lprclone", args...)
	cmd.Stdin, cmd.Stdout, cmd.Stderr = os.Stdin, os.Stdout, os.Stderr
	return cmd.Run()
}

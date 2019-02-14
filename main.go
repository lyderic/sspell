package main

import (
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"

	minio "github.com/minio/minio-go"
)

func init() {
	log.SetFlags(log.Lshortfile)
}

func main() {

	configPath := filepath.Join(os.Getenv("HOME"), "mondes", "sspell.conf")
	if len(os.Args) > 1 {
		if os.Args[1] == "-f" {
			configPath = os.Args[2]
		}
	}

	// Initialise default values
	config.Endpoint = "minio.cybermoped.com:9000"
	config.UseSSL = true
	config.Bucket = "sspell"
	config.UTF8File = "fr.utf-8.add"

	// Values above will be replaced by the ones on the config file, if found
	configFile, err := os.Open(configPath)
	if err != nil {
		log.Fatal(err)
	}
	defer configFile.Close()
	decoder := json.NewDecoder(configFile)
	err = decoder.Decode(&config)
	if err != nil {
		log.Fatalln("Decoding configuration file", configPath, "failed:", err)
	}
	config.DataFile = config.UTF8File + ".spl"
	config.SpellDir = filepath.Join(os.Getenv("HOME"), ".vim", "spell")

	minioClient, err := minio.New(config.Endpoint, config.KeyId, config.KeySecret, config.UseSSL)
	if err != nil {
		log.Fatal(err)
	}

	exists, err := minioClient.BucketExists(config.Bucket)
	if err != nil {
		log.Fatal(err)
	}
	if !exists {
		log.Fatalln("Bucket", config.Bucket, "not found!")
	}

	remoteInfo, err := minioClient.StatObject(config.Bucket, config.DataFile, minio.StatObjectOptions{})
	if err != nil {
		log.Fatal(err)
	}
	remoteSum := remoteInfo.ETag
	fmt.Println("[R] MD5 sum:", remoteSum)

	localFile := filepath.Join(config.SpellDir, config.DataFile)
	if _, err = os.Stat(localFile); os.IsNotExist(err) {
		fmt.Println("No local file:", localFile)
		files := []string{config.UTF8File, config.DataFile}
		getFromRemote(minioClient, files)
	}
	localBytes, err := ioutil.ReadFile(localFile)
	if err != nil {
		log.Fatal(err)
	}
	localSum16Byte := md5.Sum(localBytes)
	localSum := hex.EncodeToString(localSum16Byte[:])
	fmt.Println("[L] MD5 Sum:", localSum)

	if remoteSum == localSum {
		fmt.Println("No change.")
		return
	}

	remoteLastModified := remoteInfo.LastModified
	remoteUnix := remoteLastModified.Unix()
	fmt.Printf("[R] ModTime: %-19.19s [UNIX:%d]\n", remoteLastModified, remoteUnix)

	localInfo, err := os.Stat(localFile)
	if err != nil {
		log.Fatal(err)
	}
	localLastModified := localInfo.ModTime()
	localUnix := localLastModified.Unix()
	fmt.Printf("[L] ModTime: %-19.19s [UNIX:%d]\n", localLastModified, localUnix)

	files := []string{config.UTF8File, config.DataFile}

	if localUnix < remoteUnix {
		fmt.Println("Remote most recent.")
		getFromRemote(minioClient, files)
	} else if localUnix > remoteUnix {
		fmt.Println("Local most recent.")
		putToRemote(minioClient, files)
	} else {
		log.Fatal("ERROR: cannot evaluate time difference!")
	}
}

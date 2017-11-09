package main

import (
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"github.com/minio/minio-go"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
)

type Configuration struct {
	Endpoint  string
	UseSSL    bool
	Bucket    string
	UTF8File  string
	DataFile  string
	SpellDir  string
	KeyId     string
	KeySecret string
}

var conf Configuration

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
	conf.Endpoint = "minio.cybermoped.com:9000"
	conf.UseSSL = true
	conf.Bucket = "sspell"
	conf.UTF8File = "fr.utf-8.add"

	// Values above will be replaced by the ones on the config file, if found
	configFile, err := os.Open(configPath)
	if err != nil {
		log.Fatal(err)
	}
	defer configFile.Close()
	decoder := json.NewDecoder(configFile)
	err = decoder.Decode(&conf)
	if err != nil {
		log.Fatalln("Decoding configuration file", configPath, "failed:", err)
	}
	conf.DataFile = conf.UTF8File + ".spl"
	conf.SpellDir = filepath.Join(os.Getenv("HOME"), ".vim", "spell")

	minioClient, err := minio.New(conf.Endpoint, conf.KeyId, conf.KeySecret, conf.UseSSL)
	if err != nil {
		log.Fatal(err)
	}

	exists, err := minioClient.BucketExists(conf.Bucket)
	if err != nil {
		log.Fatal(err)
	}
	if !exists {
		log.Fatalln("Bucket", conf.Bucket, "not found!")
	}

	remoteInfo, err := minioClient.StatObject(conf.Bucket, conf.DataFile, minio.StatObjectOptions{})
	if err != nil {
		log.Fatal(err)
	}
	remoteSum := remoteInfo.ETag
	fmt.Println("[R] MD5 sum:", remoteSum)

	localFile := filepath.Join(conf.SpellDir, conf.DataFile)
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

	files := []string{conf.UTF8File, conf.DataFile}

	if localUnix < remoteUnix {
		fmt.Println("Remote most recent.")
    fmt.Println("Getting from remote:")
		for _, file := range files {
			err := minioClient.FGetObject(
				conf.Bucket,
				file,
				filepath.Join(conf.SpellDir, file),
				minio.GetObjectOptions{})
			if err == nil {
				fmt.Printf("> %-14.14s : ok\n", file)
			} else {
				log.Fatalln(file, ": failed!", err)
			}
		}
	} else if localUnix > remoteUnix {
		fmt.Println("Local most recent.")
    fmt.Println("Pushing to remote:")
		for _, file := range files {
			n, err := minioClient.FPutObject(
				conf.Bucket,
				file,
				filepath.Join(conf.SpellDir, file),
				minio.PutObjectOptions{})
			if err == nil {
				fmt.Printf("> %-14.14s : ok (%d bytes)\n", file, n)
			} else {
				log.Fatalln(file, ": failed!", err)
			}
		}
	} else {
		log.Fatal("ERROR: cannot evaluate time difference!")
	}
}

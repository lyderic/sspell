package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"time"

	minio "github.com/minio/minio-go"
)

const (
	VERSION = "0.1.1"
)

var (
	debug bool
)

func init() {
	log.SetFlags(log.Lshortfile)
}

func main() {

	var showversion bool
	configPath := filepath.Join(os.Getenv("HOME"), "mondes", "sspell.conf")
	flag.StringVar(&configPath, "f", configPath, "configuration file")
	flag.BoolVar(&debug, "debug", false, "show debug information")
	flag.BoolVar(&showversion, "version", false, "show version")
	flag.Usage = usage
	flag.Parse()

	if len(flag.Args()) > 0 {
		switch flag.Args()[0] {
		case "version":
			fmt.Println(VERSION)
		default:
			usage()
		}
		return
	}

	if showversion {
		fmt.Println(VERSION)
		return
	}

	if debug {
		fmt.Println("Configuration:", configPath)
	}

	configure(configPath)

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

	files := []string{config.UTF8File, config.BinaryFile}

	remoteInfo, err := minioClient.StatObject(config.Bucket, config.BinaryFile, minio.StatObjectOptions{})
	if err != nil {
		fmt.Println("Missing files on remote!!!")
		transfer(TO, config.UTF8File, minioClient)
		transfer(TO, config.BinaryFile, minioClient)
		return
	}
	fmt.Println(">>>>>", remoteInfo)
	remote.Sum = remoteInfo.ETag
	fmt.Println("[R] MD5 sum:", remote.Sum)

	localFile := filepath.Join(config.SpellDir, config.BinaryFile)
	if _, err = os.Stat(localFile); os.IsNotExist(err) {
		fmt.Println("No local file:", localFile)
		getFromRemote(minioClient, files)
	}
	local.Sum = md5sum(localFile)
	fmt.Println("[L] MD5 Sum:", local.Sum)

	if remote.Sum == local.Sum {
		fmt.Println("No change.")
		return
	}
	remote.LastModified = remoteInfo.LastModified
	remote.Unix = remote.LastModified.Unix()
	displayTime("R", remote.LastModified, remote.Unix)
	localInfo, err := os.Stat(localFile)
	if err != nil {
		log.Fatal(err)
	}
	local.LastModified = localInfo.ModTime()
	local.Unix = local.LastModified.Unix()
	displayTime("L", local.LastModified, local.Unix)
	if local.Unix < remote.Unix {
		fmt.Println("Remote most recent.")
		getFromRemote(minioClient, files)
	} else if local.Unix > remote.Unix {
		fmt.Println("Local most recent.")
		putToRemote(minioClient, files)
	} else {
		log.Fatal("ERROR: cannot evaluate time difference!")
	}
}

func displayTime(direction string, lastModified time.Time, unix int64) {
	fmt.Printf("[%s] ModTime: %-19.19s [UNIX:%d]\n", direction, lastModified, unix)
}

func usage() {
	fmt.Println("sspell, version", VERSION, "(c) Lyderic Landry, London 2019")
	fmt.Println("Usage: sspell [-h] [-f filepath]")
	flag.PrintDefaults()
}

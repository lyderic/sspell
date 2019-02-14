package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path/filepath"
)

type Configuration struct {
	Endpoint   string
	UseSSL     bool
	Bucket     string
	UTF8File   string
	BinaryFile string
	SpellDir   string
	KeyId      string
	KeySecret  string
}

var config Configuration

func configure(configPath string) {

	// Initialise default values
	config.Endpoint = "minio.cybermoped.com:9000"
	config.UseSSL = true
	config.Bucket = "sspell"
	config.UTF8File = "fr.utf-8.add"

	// Default values above will be replaced by the ones on the config file, if found
	fh, err := os.Open(configPath)
	if err != nil {
		log.Fatal(err)
	}
	defer fh.Close()
	decoder := json.NewDecoder(fh)
	err = decoder.Decode(&config)
	if err != nil {
		log.Fatalln("Decoding configuration file", configPath, "failed:", err)
	}
	config.BinaryFile = config.UTF8File + ".spl"
	config.SpellDir = filepath.Join(os.Getenv("HOME"), ".vim", "spell")

	if debug {
		fmt.Printf("%#v\n", config)
	}

	if len(config.KeyId) == 0 {
		log.Fatal("Key Id not found!")
	}

	if len(config.KeySecret) == 0 {
		log.Fatal("Key Secret not found!")
	}
}

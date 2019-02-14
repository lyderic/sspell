package main

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

var config Configuration

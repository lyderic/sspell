package main

import "time"

type Metadata struct {
	Name         string
	Sum          string
	LastModified time.Time
	Unix         int64
}

var local, remote Metadata

func init() {
	local.Name = "local"
	remote.Name = "remote"
}

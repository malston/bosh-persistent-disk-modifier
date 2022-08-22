package main

import (
	"github.com/malston/bosh-persistent-disk-modifier/bosh"
	"log"
)

const (
	host           = "127.0.0.1"
	user           = "vcap"
	password       = ""
	sshHost        = "192.168.10.20"
	sshUsername    = "jumpbox"
	sshPassword    = ""
	sshPrivateKey  = ""
	tunnelRequired = false
)

func main() {
	db, err := bosh.NewDatabase(
		host,
		user,
		password,
		sshHost,
		sshUsername,
		sshPassword,
		sshPrivateKey,
		tunnelRequired,
	)
	if err != nil {
		log.Fatalf("failed to connect to bosh database: %v", err)
	}

	b := &bosh.BOSH{
		DB: db,
	}

	_, err = b.GetPersistentDiskMappings()
	if err != nil {
		return
	}
}

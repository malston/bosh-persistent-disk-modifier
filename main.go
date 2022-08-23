package main

import (
	"log"
	"os"

	goflags "github.com/jessevdk/go-flags"
	"github.com/malston/bosh-persistent-disk-modifier/bosh"
)

const (
	host = "127.0.0.1"
	user = "vcap"
)

var opts struct {
	Deployment string `short:"n" long:"deployment" description:"A deployment name" required:"true"`
}

func main() {
	_, err := goflags.Parse(&opts)
	if err != nil {
		os.Exit(1)
	}

	db, err := bosh.NewDatabase(
		host,
		user,
	)
	if err != nil {
		log.Fatalf("failed to connect to bosh database: %v", err)
	}

	b := &bosh.BOSH{
		DB: db,
	}

	err = b.UpdatePersistentDiskCIDs(opts.Deployment)
	if err != nil {
		log.Fatalf("%v\n", err)
	}
}

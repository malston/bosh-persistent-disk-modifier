package main

import (
	"log"
	"net/url"
	"os"

	goflags "github.com/jessevdk/go-flags"
	"github.com/malston/bosh-persistent-disk-modifier/bosh"
	"github.com/vmware/govmomi/vim25"
)

const (
	host = "127.0.0.1"
	user = "vcap"
)

var opts struct {
	VCenterHostname string `long:"vcenter" description:"The vcenter hostname" required:"true"`
	VCenterUsername string `short:"u" long:"username" description:"The vcenter username" required:"true"`
	VCenterPassword string `short:"p" long:"password" description:"The vcenter password" required:"true"`
	Deployment string `short:"n" long:"deployment" description:"The bosh deployment name" required:"true"`
}

func main() {
	_, err := goflags.Parse(&opts)
	if err != nil {
		os.Exit(1)
	}

	db, err := bosh.NewDBConnection(
		host,
		user,
	)
	if err != nil {
		log.Fatalf("failed to connect to bosh database: %v", err)
	}

	r := &bosh.Repository{
		DB: db,
	}
	u, err := url.Parse("https://"+opts.VCenterHostname)
	if err != nil {
		log.Fatalf("unable to parse url: %v", err)
	}

	u.Path = vim25.Path
	u.User = url.UserPassword(opts.VCenterUsername, opts.VCenterPassword)

	err = r.UpdatePersistentDiskCIDs(opts.Deployment, u, true)
	if err != nil {
		log.Fatalf("%v\n", err)
	}
}

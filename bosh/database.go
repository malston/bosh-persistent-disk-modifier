package bosh

import (
	"fmt"
	"github.com/jmoiron/sqlx"
	"github.com/malston/bosh-persistent-disk-modifier/ssh"
	"sync"

	_ "github.com/lib/pq"
)

func NewDatabase(
	host string,
	user string,
	password string,
	tunnelHost string,
	tunnelUser string,
	tunnelPassword string,
	tunnelPrivateKey string,
	tunnelRequired bool,
) (*sqlx.DB, error) {
	sshTunnel, err := ssh.NewTunnel(host, tunnelHost, tunnelUser, tunnelPassword, tunnelPrivateKey, tunnelRequired)

	if err != nil {
		return nil, err
	}

	host = fmt.Sprintf("%s:5432", host)
	wg := &sync.WaitGroup{}
	if sshTunnel != nil {
		wg.Add(1)
		go func() {
			err := sshTunnel.Start(wg)
			if err != nil {
				fmt.Println("failed to start ssh tunnel")
				panic(err)
			}
		}()
		host = fmt.Sprintf("localhost:%d", sshTunnel.Local.Port)
	}
	wg.Wait()

	conn := fmt.Sprintf("postgres://%s@%s/bosh?sslmode=disable", user, host)
	fmt.Printf("connecting to bosh db: %s\n", conn)
	db, err := sqlx.Connect("postgres", conn)
	if err != nil {
		fmt.Println("failed to connect to bosh db")
		return db, err
	}

	fmt.Println("connected to bosh db")
	return db, err
}

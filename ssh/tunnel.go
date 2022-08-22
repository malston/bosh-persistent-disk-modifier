package ssh

import (
	"fmt"
	"io"
	"log"
	"net"
	"sync"

	"golang.org/x/crypto/ssh"
)

// Endpoint represents a Host/port combo
type Endpoint struct {
	// Server host address
	Host string
	// Server port
	Port int
}

func (endpoint *Endpoint) String() string {
	return fmt.Sprintf("%s:%d", endpoint.Host, endpoint.Port)
}

// SSH represents a local/remote/tunnel combination
type SSH struct {
	Local  *Endpoint
	Server *Endpoint
	Remote *Endpoint

	Config *ssh.ClientConfig
}

// Start opens a persistent tunnel
func (tunnel *SSH) Start(wg *sync.WaitGroup) error {
	log.Printf("starting ssh tunnel: %+v", tunnel)
	listener, err := net.Listen("tcp", tunnel.Local.String())
	if err != nil {
		return err
	}
	defer listener.Close()
	wg.Done()
	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Printf("failed to connect to ssh tunnel: %#v", tunnel)
			return err
		}
		go tunnel.forward(conn)
	}
}

func (tunnel *SSH) forward(localConn net.Conn) {
	serverConn, err := ssh.Dial("tcp", tunnel.Server.String(), tunnel.Config)
	if err != nil {
		log.Fatalf("Server dial error: %s\n", err)
	}

	remoteConn, err := serverConn.Dial("tcp", tunnel.Remote.String())
	if err != nil {
		log.Fatalf("Remote dial error: %s\n", err)
	}

	copyConn := func(writer, reader net.Conn) {
		_, _ = io.Copy(writer, reader)
	}

	go copyConn(localConn, remoteConn)
	go copyConn(remoteConn, localConn)
}

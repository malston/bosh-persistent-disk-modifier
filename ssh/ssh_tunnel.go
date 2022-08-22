package ssh

import (
	"errors"
	"io/ioutil"
	"net"
	"os"
	"time"

	"golang.org/x/crypto/ssh"
	"golang.org/x/crypto/ssh/agent"
)

func NewTunnel(
	host,
	tunnelHost,
	tunnelUser,
	tunnelPassword,
	tunnelPrivateKey string,
	tunnelRequired bool,
) (*SSH, error) {
	if !tunnelRequired {
		// log.Println("ssh tunneling not required, skipping tunnel creation")
		return nil, nil
	}

	if tunnelHost == "" || tunnelUser == "" || (tunnelPassword == "" && tunnelPrivateKey == "") {
		return nil, errors.New("tunneling is required, but the tunnel information was not specified")
	}

	localPort := findLocalPort()
	if localPort < 0 {
		return nil, errors.New("no available ports on localhost")
	}

	localServer := Endpoint{
		Host: "localhost",
		Port: localPort,
	}

	remoteServer := Endpoint{
		Host: host,
		Port: 5432,
	}

	tunnelServer := Endpoint{
		Host: tunnelHost,
		Port: 22,
	}

	var authMethods []ssh.AuthMethod
	if tunnelPrivateKey != "" {
		pem, err := ioutil.ReadFile(tunnelPrivateKey)
		if err != nil {
			return nil, err
		}

		signer, err := ssh.ParsePrivateKey(pem)
		if err != nil {
			return nil, err
		}

		authMethods = append(authMethods, ssh.PublicKeys(signer))
	}

	if tunnelPassword != "" {
		authMethods = append(authMethods, ssh.Password(tunnelPassword))
	}

	if sshAgent, err := net.Dial("unix", os.Getenv("SSH_AUTH_SOCK")); err == nil {
		authMethods = append(authMethods, ssh.PublicKeysCallback(agent.NewClient(sshAgent).Signers))
	}

	config := ssh.ClientConfig{
		Timeout: 5 * time.Second,
		User:    tunnelUser,
		Auth:    authMethods,
		HostKeyCallback: func(hostname string, remote net.Addr, key ssh.PublicKey) error {
			return nil
		},
	}

	return &SSH{
		Config: &config,
		Local:  &localServer,
		Server: &tunnelServer,
		Remote: &remoteServer,
	}, nil
}

func findLocalPort() int {
	addr, err := net.ResolveTCPAddr("tcp", "localhost:0")
	if err != nil {
		return -1
	}

	l, err := net.ListenTCP("tcp", addr)
	if err != nil {
		return -1
	}
	defer l.Close()
	return l.Addr().(*net.TCPAddr).Port
}

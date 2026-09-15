package engine

import (
	"fmt"
	"io"
	"net"
	"time"

	"golang.org/x/crypto/ssh"
)

type SSHTunnel struct {
	LocalHost  string
	LocalPort  int
	listener   net.Listener
	sshClient  *ssh.Client
	RemoteHost string
	RemotePort int
}

func NewSSHTunnel(sshHost string, sshPort int, sshUser, sshPass, sshKey string, remoteHost string, remotePort int) (*SSHTunnel, error) {
	var authMethods []ssh.AuthMethod

	if sshKey != "" {
		signer, err := ssh.ParsePrivateKey([]byte(sshKey))
		if err != nil {
			return nil, fmt.Errorf("invalid ssh private key: %w", err)
		}
		authMethods = append(authMethods, ssh.PublicKeys(signer))
	} else if sshPass != "" {
		authMethods = append(authMethods, ssh.Password(sshPass))
	} else {
		return nil, fmt.Errorf("ssh password or key is required")
	}

	config := &ssh.ClientConfig{
		User:            sshUser,
		Auth:            authMethods,
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
		Timeout:         10 * time.Second,
	}

	sshAddr := fmt.Sprintf("%s:%d", sshHost, sshPort)
	client, err := ssh.Dial("tcp", sshAddr, config)
	if err != nil {
		return nil, fmt.Errorf("ssh dial failed: %w", err)
	}

	// Listen on ephemeral port
	listener, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		client.Close()
		return nil, fmt.Errorf("ssh local listen failed: %w", err)
	}

	localPort := listener.Addr().(*net.TCPAddr).Port

	tunnel := &SSHTunnel{
		LocalHost:  "127.0.0.1",
		LocalPort:  localPort,
		listener:   listener,
		sshClient:  client,
		RemoteHost: remoteHost,
		RemotePort: remotePort,
	}

	go tunnel.acceptLoop()

	return tunnel, nil
}

func (t *SSHTunnel) acceptLoop() {
	remoteAddr := fmt.Sprintf("%s:%d", t.RemoteHost, t.RemotePort)
	for {
		localConn, err := t.listener.Accept()
		if err != nil {
			break
		}

		go func(lConn net.Conn) {
			defer lConn.Close()
			remoteConn, err := t.sshClient.Dial("tcp", remoteAddr)
			if err != nil {
				return
			}
			defer remoteConn.Close()

			errc := make(chan error, 2)
			go func() {
				_, err := io.Copy(remoteConn, lConn)
				errc <- err
			}()
			go func() {
				_, err := io.Copy(lConn, remoteConn)
				errc <- err
			}()
			<-errc
		}(localConn)
	}
}

func (t *SSHTunnel) Close() error {
	var errs []error
	if err := t.listener.Close(); err != nil {
		errs = append(errs, err)
	}
	if err := t.sshClient.Close(); err != nil {
		errs = append(errs, err)
	}
	if len(errs) > 0 {
		return fmt.Errorf("ssh tunnel close errors: %v", errs)
	}
	return nil
}

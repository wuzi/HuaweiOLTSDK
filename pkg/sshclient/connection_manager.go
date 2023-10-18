package sshclient

import (
	"errors"
	"golang.org/x/crypto/ssh"
	"io"
)

type ConnectionManager struct {
	Host       string
	Port       string
	JumpClient *ConnectionManager
	SSHConfig  *ssh.ClientConfig
	Connection *ssh.Client
	Session    *ssh.Session
	Stdout     io.Reader
	Stdin      io.WriteCloser
}

func NewClient(user, password, host, port string) *ConnectionManager {
	sshConfig := &ssh.ClientConfig{
		User:            user,
		Auth:            []ssh.AuthMethod{ssh.Password(password)},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
		Config: ssh.Config{
			KeyExchanges: []string{"diffie-hellman-group-exchange-sha256"},
		},
	}
	return &ConnectionManager{
		SSHConfig: sshConfig,
		Host:      host,
		Port:      port,
	}
}

func (c *ConnectionManager) Connect() error {
	conn, err := ssh.Dial("tcp", c.Host+":"+c.Port, c.SSHConfig)
	if err != nil {
		return err
	}
	c.Connection = conn

	err = c.createSession()
	if err != nil {
		return err
	}
	return nil
}

func (c *ConnectionManager) ConnectThroughJumpHost(user, password, host, port string) error {
	var err error

	c.JumpClient = NewClient(user, password, host, port)
	c.JumpClient.Connection, err = ssh.Dial("tcp", c.JumpClient.Host+":"+c.JumpClient.Port, c.JumpClient.SSHConfig)
	if err != nil {
		return err
	}

	conn, err := c.JumpClient.Connection.Dial("tcp", c.Host+":"+c.Port)
	if err != nil {
		return err
	}

	connSSH, conSSHChan, connSSHReq, err := ssh.NewClientConn(conn, c.Host, c.SSHConfig)
	if err != nil {
		return err
	}

	c.Connection = ssh.NewClient(connSSH, conSSHChan, connSSHReq)

	err = c.createSession()
	if err != nil {
		return err
	}

	return nil
}

func (c *ConnectionManager) Close() error {
	err := c.Connection.Close()

	if err != nil {
		return err
	}

	if c.JumpClient != nil {
		return c.JumpClient.Close()
	}

	return err
}

func (c *ConnectionManager) Wait() error {
	err := c.Session.Wait()
	if err != nil {
		var exitMissingError *ssh.ExitMissingError
		if errors.As(err, &exitMissingError) {
			return nil
		}
		return err
	}
	return nil
}

func (c *ConnectionManager) createSession() error {
	session, err := c.Connection.NewSession()
	if err != nil {
		return err
	}
	c.Session = session
	return nil
}

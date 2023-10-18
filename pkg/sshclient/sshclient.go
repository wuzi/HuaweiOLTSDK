package sshclient

import (
	"errors"
	"fmt"
	"golang.org/x/crypto/ssh"
	"io"
	"regexp"
	"strings"
)

type Context struct {
	Level int
	Frame string
	Slot  string
}

type Client struct {
	Host       string
	Port       string
	JumpClient *Client
	SSHConfig  *ssh.ClientConfig
	Connection *ssh.Client
	Session    *ssh.Session
	Stdout     io.Reader
	Stdin      io.WriteCloser
	Context    Context
}

func NewClient(user, password, host, port string) *Client {
	sshConfig := &ssh.ClientConfig{
		User:            user,
		Auth:            []ssh.AuthMethod{ssh.Password(password)},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
		Config: ssh.Config{
			KeyExchanges: []string{"diffie-hellman-group-exchange-sha256"},
		},
	}
	return &Client{
		SSHConfig: sshConfig,
		Host:      host,
		Port:      port,
		Context:   Context{},
	}
}

func (c *Client) Connect() error {
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

func (c *Client) ConnectThroughJumpHost(user, password, host, port string) error {
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

func (c *Client) Close() error {
	err := c.Connection.Close()

	if err != nil {
		return err
	}

	if c.JumpClient != nil {
		return c.JumpClient.Close()
	}

	return err
}

func (c *Client) ReadUntilPrompt(prompt string) (string, error) {
	output := make([]byte, 4096)
	var accumulatedOutput []byte
	for {
		n, err := c.Stdout.Read(output)
		if err != nil {
			if err == io.EOF {
				break
			}
			return "", fmt.Errorf("failed to read output: %v", err)
		}

		buffer := output[:n]
		text := string(buffer)

		if strings.Contains(text, "---- More ( Press 'Q' to break ) ----") {
			_, err := c.Stdin.Write([]byte("\n"))
			if err != nil {
				return "", err
			}
		}

		text = strings.Replace(text, "---- More ( Press 'Q' to break ) ----", "", -1)
		re := regexp.MustCompile(`.\[37D`)
		text = re.ReplaceAllString(text, "")

		accumulatedOutput = append(accumulatedOutput, []byte(text)...)
		if strings.Contains(string(accumulatedOutput), prompt) {
			break
		}
	}
	return string(accumulatedOutput), nil
}

func (c *Client) RunCommand(command, prompt string) (string, error) {
	_, err := c.Stdin.Write([]byte(command + "\n"))
	if err != nil {
		return "", err
	}

	output, err := c.ReadUntilPrompt(prompt)
	if err != nil {
		return "", err
	}

	return output, nil
}

func (c *Client) Wait() error {
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

func (c *Client) createSession() error {
	session, err := c.Connection.NewSession()
	if err != nil {
		return err
	}

	stdout, err := session.StdoutPipe()
	if err != nil {
		return err
	}

	c.Stdout = stdout
	stdin, err := session.StdinPipe()
	if err != nil {
		return err
	}

	c.Stdin = stdin
	err = session.Shell()
	if err != nil {
		return err
	}

	_, err = c.ReadUntilPrompt("MA5683T>")
	if err != nil {
		return err
	}

	c.Session = session
	return nil
}

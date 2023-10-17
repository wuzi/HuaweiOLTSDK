package sshclient

import (
	"errors"
	"fmt"
	"golang.org/x/crypto/ssh"
	"io"
	"regexp"
	"strings"
)

type Client struct {
	SSHConfig  *ssh.ClientConfig
	Connection *ssh.Client
	Session    *ssh.Session
	JumpClient *ssh.Client
	Host       string
	Port       string
	Stdout     io.Reader
	Stdin      io.WriteCloser
	Level      int
}

func NewClient(user, password, host, port string) *Client {
	sshConfig := createSSHClientConfig(user, password)

	return &Client{
		SSHConfig: sshConfig,
		Host:      host,
		Port:      port,
	}
}

func (c *Client) Connect() error {
	conn, err := ssh.Dial("tcp", c.Host+":"+c.Port, c.SSHConfig)
	if err != nil {
		return err
	}
	return c.createSession(conn)
}

func (c *Client) ConnectThroughJumpHost(jumpHost, jumpUser, jumpPassword string) error {
	config := createSSHClientConfig(jumpUser, jumpPassword)

	jumpClient, err := ssh.Dial("tcp", jumpHost+":22", config)
	if err != nil {
		return err
	}

	targetConn, err := jumpClient.Dial("tcp", c.Host+":22")
	if err != nil {
		return err
	}

	targetSSH, targetSSHChans, targetSSHReqs, err := ssh.NewClientConn(targetConn, c.Host, c.SSHConfig)
	if err != nil {
		return err
	}

	targetClient := ssh.NewClient(targetSSH, targetSSHChans, targetSSHReqs)
	err = c.createSession(targetClient)
	if err != nil {
		return err
	}

	c.JumpClient = jumpClient
	return nil
}

func (c *Client) Close() error {
	err1 := c.Connection.Close()
	err2 := c.JumpClient.Close()
	if err1 != nil {
		return err1
	}
	return err2
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

func (c *Client) createSession(conn *ssh.Client) error {
	session, err := conn.NewSession()
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
	c.Connection = conn
	return nil
}

func createSSHClientConfig(user, password string) *ssh.ClientConfig {
	return &ssh.ClientConfig{
		User:            user,
		Auth:            []ssh.AuthMethod{ssh.Password(password)},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
		Config: ssh.Config{
			KeyExchanges: []string{"diffie-hellman-group-exchange-sha256"},
		},
	}
}

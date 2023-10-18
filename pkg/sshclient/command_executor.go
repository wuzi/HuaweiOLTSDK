package sshclient

import (
	"fmt"
	"io"
	"regexp"
	"strings"
)

type ExecutorContext struct {
	Level int
	Frame string
	Slot  string
}

type CommandExecutor struct {
	Stdout          io.Reader
	Stdin           io.WriteCloser
	ExecutorContext ExecutorContext
}

func NewCommandExecutor(connManager *ConnectionManager) (*CommandExecutor, error) {
	stdout, err := connManager.Session.StdoutPipe()
	if err != nil {
		return nil, err
	}

	stdin, err := connManager.Session.StdinPipe()
	if err != nil {
		return nil, err
	}

	err = connManager.Session.Shell()
	if err != nil {
		return nil, err
	}

	commExecutor := &CommandExecutor{
		Stdout:          stdout,
		Stdin:           stdin,
		ExecutorContext: ExecutorContext{},
	}

	_, err = commExecutor.ReadUntilPrompt("MA5683T>")
	if err != nil {
		return nil, err
	}

	return commExecutor, nil
}

func (c *CommandExecutor) RunCommand(command, prompt string) (string, error) {
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

func (c *CommandExecutor) ReadUntilPrompt(prompt string) (string, error) {
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

func (c *CommandExecutor) Enable() error {
	if c.ExecutorContext.Level != 0 {
		return fmt.Errorf("not in root mode")
	}
	output, err := c.RunCommand("enable", "MA5683T#")
	if err != nil {
		return fmt.Errorf("failed to run command enable: %v", err)
	}
	c.ExecutorContext.Level = 1
	fmt.Print(output)
	return nil
}

func (c *CommandExecutor) Config() error {
	if c.ExecutorContext.Level != 1 {
		return fmt.Errorf("not in enable mode")
	}
	output, err := c.RunCommand("config", "MA5683T(config)#")
	if err != nil {
		return fmt.Errorf("failed to run command config: %v", err)
	}
	c.ExecutorContext.Level = 2
	fmt.Print(output)
	return nil
}

func (c *CommandExecutor) Quit(exit bool) error {
	var output string
	var err error

	if c.ExecutorContext.Level >= 3 {
		output, err = c.RunCommand("quit", "MA5683T(config)#")
		if err != nil {
			return fmt.Errorf("failed to run command: %v", err)
		}
		fmt.Print(output)
		c.ExecutorContext.Level = 2
		if !exit {
			return nil
		}
	}

	if c.ExecutorContext.Level >= 2 {
		output, err = c.RunCommand("quit", "MA5683T#")
		if err != nil {
			return fmt.Errorf("failed to run command: %v", err)
		}
		fmt.Print(output)
		c.ExecutorContext.Level = 1
		if !exit {
			return nil
		}
	}

	output, err = c.RunCommand("quit", "before logout")
	if err != nil {
		return fmt.Errorf("failed to run command: %v", err)
	}
	fmt.Print(output)

	output, err = c.RunCommand("y", "to log on")
	if err != nil {
		return fmt.Errorf("failed to run command: %v", err)
	}
	fmt.Print(output)
	return nil
}

func (c *CommandExecutor) DisplayUnmanagedOnt() ([]UnmanagedONT, error) {
	if c.ExecutorContext.Level != 2 {
		return nil, fmt.Errorf("not in config mode")
	}
	output, err := c.RunCommand("display ont autofind all", "MA5683T(config)#")
	if err != nil {
		return nil, fmt.Errorf("failed to run command: %v", err)
	}
	fmt.Print(output)
	return ParseUnmanagedONT(output)
}

func (c *CommandExecutor) InterfaceGPON(frame string, slot string) error {
	if c.ExecutorContext.Level != 2 {
		return fmt.Errorf("not in config mode")
	}
	output, err := c.RunCommand(fmt.Sprintf("interface gpon %s/%s", frame, slot), fmt.Sprintf("MA5683T(config-if-gpon-%s/%s)#", frame, slot))
	if err != nil {
		return fmt.Errorf("failed to run command: %v", err)
	}
	c.ExecutorContext.Level = 3
	c.ExecutorContext.Frame = frame
	c.ExecutorContext.Slot = slot
	fmt.Print(output)
	return nil
}

func (c *CommandExecutor) AddOnt(port string, serialNumber string, description string) (string, error) {
	if c.ExecutorContext.Level != 3 {
		return "", fmt.Errorf("not in interface gpon mode")
	}

	output, err := c.RunCommand(fmt.Sprintf("ont add %s sn-auth %s omci ont-lineprofile-id 60 ont-srvprofile-id 35 desc %s",
		port,
		serialNumber,
		description,
	), fmt.Sprintf("MA5683T(config-if-gpon-%s/%s)#", c.ExecutorContext.Frame, c.ExecutorContext.Slot))

	if err != nil {
		return "", fmt.Errorf("failed to run command: %v", err)
	}
	fmt.Print(output)

	if strings.Contains(output, "Failure: SN already exists") {
		return "", fmt.Errorf("serial number already exists")
	}

	re := regexp.MustCompile(`ONTID :(\d+)`)
	match := re.FindStringSubmatch(output)
	if len(match) < 2 {
		return "", fmt.Errorf("ONTID not found in command output")
	}

	return match[1], nil
}

func (c *CommandExecutor) DeleteOnt(port string) error {
	if c.ExecutorContext.Level != 3 {
		return fmt.Errorf("not in interface gpon mode")
	}

	output, err := c.RunCommand(fmt.Sprintf("ont delete %s all", port), "(y/n)[n]:")
	if err != nil {
		return fmt.Errorf("failed to run command: %v", err)
	}
	fmt.Print(output)

	output, err = c.RunCommand("y", fmt.Sprintf("MA5683T(config-if-gpon-%s/%s)#", c.ExecutorContext.Frame, c.ExecutorContext.Slot))
	if err != nil {
		return fmt.Errorf("failed to run command: %v", err)
	}
	fmt.Print(output)
	return nil
}

func (c *CommandExecutor) AddNativeVlan(port string, ontID string) error {
	if c.ExecutorContext.Level != 3 {
		return fmt.Errorf("not in interface gpon mode")
	}

	output, err := c.RunCommand(fmt.Sprintf("ont port native-vlan %s %s eth 1 vlan 20 priority 0", port, ontID), fmt.Sprintf("MA5683T(config-if-gpon-%s/%s)#", c.ExecutorContext.Frame, c.ExecutorContext.Slot))
	if err != nil {
		return fmt.Errorf("failed to run command: %v", err)
	}

	if strings.Contains(output, "Failure: Make configuration repeatedly") {
		return fmt.Errorf("make configuration repeatedly")
	}

	fmt.Print(output)
	return nil
}

func (c *CommandExecutor) AddServicePort(vlan string, frame string, slot string, port string, ontID string) error {
	if c.ExecutorContext.Level != 2 {
		return fmt.Errorf("not in config mode")
	}

	output, err := c.RunCommand(fmt.Sprintf("service-port vlan %s gpon %s/%s/%s ont %s gemport 20 multi-service user-vlan 20 tag-transform translate inbound traffic-table index 10 outbound traffic-table index 10", vlan, frame, slot, port, ontID), "MA5683T(config)#")
	if err != nil {
		return fmt.Errorf("failed to run command: %v", err)
	}

	if strings.Contains(output, "Failure: VLAN does not exist") {
		return fmt.Errorf("VLAN does not exist")
	}

	fmt.Print(output)
	return nil
}

func (c *CommandExecutor) GetOntData(frame string, slot string, port string, ontID string) (*ONT, error) {
	if c.ExecutorContext.Level != 2 {
		return nil, fmt.Errorf("not in config mode")
	}

	output, err := c.RunCommand(fmt.Sprintf("display current-configuration ont %s/%s/%s %s", frame, slot, port, ontID), "MA5683T(config)#")
	if err != nil {
		return nil, fmt.Errorf("failed to run command: %v", err)
	}

	if strings.Contains(output, "Failure") || strings.Contains(output, "Error") {
		return nil, fmt.Errorf("could not get service port")
	}

	fmt.Print(output)

	ont := &ONT{
		Frame: frame,
		Slot:  slot,
		Port:  port,
		ID:    ontID,
	}

	if match := regexp.MustCompile(`sn-auth "(.*?)"`).FindStringSubmatch(output); len(match) > 1 {
		ont.SerialNumber = match[1]
	}

	if match := regexp.MustCompile(`desc "(.*?)"`).FindStringSubmatch(output); len(match) > 1 {
		ont.Description = match[1]
	}

	if match := regexp.MustCompile(`service-port (\d+)`).FindStringSubmatch(output); len(match) > 1 {
		ont.ServicePort = match[1]
	}

	if match := regexp.MustCompile(`vlan (\d+) gpon`).FindStringSubmatch(output); len(match) > 1 {
		ont.VlanID = match[1]
	}

	return ont, nil
}

func (c *CommandExecutor) UndoServicePort(id string) error {
	if c.ExecutorContext.Level != 2 {
		return fmt.Errorf("not in config mode")
	}

	output, err := c.RunCommand(fmt.Sprintf("undo service-port %s", id), "MA5683T(config)#")
	if err != nil {
		return fmt.Errorf("failed to run command: %v", err)
	}

	if strings.Contains(output, "Failure") {
		return fmt.Errorf("could not undo service port")
	}

	fmt.Print(output)
	return nil
}

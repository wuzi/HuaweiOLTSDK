package sshclient

import (
	"fmt"
	"regexp"
	"strings"
)

func (c *Client) Enable() error {
	if c.Context.Level != 0 {
		return fmt.Errorf("not in root mode")
	}
	output, err := c.RunCommand("enable", "MA5683T#")
	if err != nil {
		return fmt.Errorf("failed to run command enable: %v", err)
	}
	c.Context.Level = 1
	fmt.Print(output)
	return nil
}

func (c *Client) Config() error {
	if c.Context.Level != 1 {
		return fmt.Errorf("not in enable mode")
	}
	output, err := c.RunCommand("config", "MA5683T(config)#")
	if err != nil {
		return fmt.Errorf("failed to run command config: %v", err)
	}
	c.Context.Level = 2
	fmt.Print(output)
	return nil
}

func (c *Client) Quit() error {
	var output string
	var err error

	if c.Context.Level >= 3 {
		output, err = c.RunCommand("quit", "MA5683T(config)#")
		if err != nil {
			return fmt.Errorf("failed to run command: %v", err)
		}
		fmt.Print(output)
	}

	if c.Context.Level >= 2 {
		output, err = c.RunCommand("quit", "MA5683T#")
		if err != nil {
			return fmt.Errorf("failed to run command: %v", err)
		}
		fmt.Print(output)
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

func (c *Client) DisplayUnmanagedOnt() ([]ONT, error) {
	if c.Context.Level != 2 {
		return nil, fmt.Errorf("not in config mode")
	}
	output, err := c.RunCommand("display ont autofind all", "MA5683T(config)#")
	if err != nil {
		return nil, fmt.Errorf("failed to run command: %v", err)
	}
	fmt.Print(output)
	return ParseUnmanagedONT(output)
}

func (c *Client) InterfaceGPON(frame int, slot int) error {
	if c.Context.Level != 2 {
		return fmt.Errorf("not in config mode")
	}
	output, err := c.RunCommand(fmt.Sprintf("interface gpon %d/%d", frame, slot), fmt.Sprintf("MA5683T(config-if-gpon-%d/%d)#", frame, slot))
	if err != nil {
		return fmt.Errorf("failed to run command: %v", err)
	}
	c.Context.Level = 3
	c.Context.Frame = frame
	c.Context.Slot = slot
	fmt.Print(output)
	return nil
}

func (c *Client) AddOnt(port int, serialNumber string, description string) (string, error) {
	if c.Context.Level != 3 {
		return "", fmt.Errorf("not in interface gpon mode")
	}

	output, err := c.RunCommand(fmt.Sprintf("ont add %d sn-auth %s omci ont-lineprofile-id 60 ont-srvprofile-id 35 desc %s",
		port,
		serialNumber,
		description,
	), fmt.Sprintf("MA5683T(config-if-gpon-%d/%d)#", c.Context.Frame, c.Context.Slot))

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

func (c *Client) DeleteOnt(port int) error {
	if c.Context.Level != 3 {
		return fmt.Errorf("not in interface gpon mode")
	}

	output, err := c.RunCommand(fmt.Sprintf("ont delete %d all", port), "(y/n)[n]:")
	if err != nil {
		return fmt.Errorf("failed to run command: %v", err)
	}
	fmt.Print(output)

	output, err = c.RunCommand("y", fmt.Sprintf("MA5683T(config-if-gpon-%d/%d)#", c.Context.Frame, c.Context.Slot))
	if err != nil {
		return fmt.Errorf("failed to run command: %v", err)
	}
	fmt.Print(output)
	return nil
}

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

func (c *Client) Quit(exit bool) error {
	var output string
	var err error

	if c.Context.Level >= 3 {
		output, err = c.RunCommand("quit", "MA5683T(config)#")
		if err != nil {
			return fmt.Errorf("failed to run command: %v", err)
		}
		fmt.Print(output)
		c.Context.Level = 2
		if !exit {
			return nil
		}
	}

	if c.Context.Level >= 2 {
		output, err = c.RunCommand("quit", "MA5683T#")
		if err != nil {
			return fmt.Errorf("failed to run command: %v", err)
		}
		fmt.Print(output)
		c.Context.Level = 1
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

func (c *Client) DisplayUnmanagedOnt() ([]UnmanagedONT, error) {
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

func (c *Client) InterfaceGPON(frame string, slot string) error {
	if c.Context.Level != 2 {
		return fmt.Errorf("not in config mode")
	}
	output, err := c.RunCommand(fmt.Sprintf("interface gpon %s/%s", frame, slot), fmt.Sprintf("MA5683T(config-if-gpon-%s/%s)#", frame, slot))
	if err != nil {
		return fmt.Errorf("failed to run command: %v", err)
	}
	c.Context.Level = 3
	c.Context.Frame = frame
	c.Context.Slot = slot
	fmt.Print(output)
	return nil
}

func (c *Client) AddOnt(port string, serialNumber string, description string) (string, error) {
	if c.Context.Level != 3 {
		return "", fmt.Errorf("not in interface gpon mode")
	}

	output, err := c.RunCommand(fmt.Sprintf("ont add %s sn-auth %s omci ont-lineprofile-id 60 ont-srvprofile-id 35 desc %s",
		port,
		serialNumber,
		description,
	), fmt.Sprintf("MA5683T(config-if-gpon-%s/%s)#", c.Context.Frame, c.Context.Slot))

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

func (c *Client) DeleteOnt(port string) error {
	if c.Context.Level != 3 {
		return fmt.Errorf("not in interface gpon mode")
	}

	output, err := c.RunCommand(fmt.Sprintf("ont delete %s all", port), "(y/n)[n]:")
	if err != nil {
		return fmt.Errorf("failed to run command: %v", err)
	}
	fmt.Print(output)

	output, err = c.RunCommand("y", fmt.Sprintf("MA5683T(config-if-gpon-%s/%s)#", c.Context.Frame, c.Context.Slot))
	if err != nil {
		return fmt.Errorf("failed to run command: %v", err)
	}
	fmt.Print(output)
	return nil
}

func (c *Client) AddNativeVlan(port string, ontID string) error {
	if c.Context.Level != 3 {
		return fmt.Errorf("not in interface gpon mode")
	}

	output, err := c.RunCommand(fmt.Sprintf("ont port native-vlan %s %s eth 1 vlan 20 priority 0", port, ontID), fmt.Sprintf("MA5683T(config-if-gpon-%s/%s)#", c.Context.Frame, c.Context.Slot))
	if err != nil {
		return fmt.Errorf("failed to run command: %v", err)
	}

	if strings.Contains(output, "Failure: Make configuration repeatedly") {
		return fmt.Errorf("make configuration repeatedly")
	}

	fmt.Print(output)
	return nil
}

func (c *Client) AddServicePort(vlan string, frame string, slot string, port string, ontID string) error {
	if c.Context.Level != 2 {
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

func (c *Client) GetOntData(frame string, slot string, port string, ontID string) (*ONT, error) {
	if c.Context.Level != 2 {
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

func (c *Client) UndoServicePort(id string) error {
	if c.Context.Level != 2 {
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

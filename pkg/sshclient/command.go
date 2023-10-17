package sshclient

import "fmt"

func (c *Client) Enable() error {
	output, err := c.RunCommand("enable", "MA5683T#")
	if err != nil {
		return fmt.Errorf("failed to run command enable: %v", err)
	}
	c.Level = 0
	fmt.Print(output)
	return nil
}

func (c *Client) Config() error {
	output, err := c.RunCommand("config", "MA5683T(config)#")
	if err != nil {
		return fmt.Errorf("failed to run command config: %v", err)
	}
	c.Level = 1
	fmt.Print(output)
	return nil
}

func (c *Client) Quit() error {
	var output string
	var err error

	if c.Level >= 2 {
		output, err = c.RunCommand("quit", "MA5683T(config)#")
		if err != nil {
			return fmt.Errorf("failed to run command: %v", err)
		}
		fmt.Print(output)
	}

	if c.Level >= 1 {
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

func (c *Client) DisplayUnmanagedONT() ([]ONT, error) {
	output, err := c.RunCommand("display ont autofind all", "MA5683T(config)#")
	if err != nil {
		return nil, fmt.Errorf("failed to run command: %v", err)
	}
	fmt.Print(output)
	return ParseUnmanagedONT(output)
}

func (c *Client) InterfaceGPON(frame int, slot int) error {
	output, err := c.RunCommand(fmt.Sprintf("interface gpon %d/%d", frame, slot), fmt.Sprintf("MA5683T(config-if-gpon-%d/%d)#", frame, slot))
	if err != nil {
		return fmt.Errorf("failed to run command: %v", err)
	}
	c.Level = 2
	fmt.Print(output)
	return nil
}

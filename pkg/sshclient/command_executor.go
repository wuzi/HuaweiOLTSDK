package sshclient

import (
	"fmt"
	"io"
	"regexp"
	"strconv"
	"strings"
)

type ExecutorContext struct {
	Level int
	Frame int
	Slot  int
}

type CommandExecutor struct {
	Verbose         bool
	Stdout          io.Reader
	Stdin           io.WriteCloser
	ExecutorContext ExecutorContext
}

type CommandExecutorOptions struct {
	Verbose bool
}

func NewCommandExecutor(connManager *ConnectionManager, options CommandExecutorOptions) (*CommandExecutor, error) {
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
		Verbose:         options.Verbose,
	}

	_, err = commExecutor.readOutputUntilPrompt("MA5683T>")
	if err != nil {
		return nil, err
	}

	err = commExecutor.enable()
	if err != nil {
		return nil, err
	}

	err = commExecutor.config()
	if err != nil {
		return nil, err
	}

	return commExecutor, nil
}

func (c *CommandExecutor) ExecuteCommand(command, prompt string) (string, error) {
	_, err := c.Stdin.Write([]byte(command + "\n"))
	if err != nil {
		return "", err
	}

	output, err := c.readOutputUntilPrompt(prompt)
	if err != nil {
		return "", err
	}

	if c.Verbose {
		fmt.Print(output)
	}

	return output, nil
}

func (c *CommandExecutor) ExitCommandLevel() error {
	return c.quit(false)
}

func (c *CommandExecutor) ExitCommandSession() error {
	return c.quit(true)
}

func (c *CommandExecutor) GetUnmanagedOpticalNetworkTerminals() ([]ONTDetail, error) {
	if c.ExecutorContext.Level != 2 {
		return nil, fmt.Errorf("not in config mode")
	}
	output, err := c.ExecuteCommand("display ont autofind all", "MA5683T(config)#")
	if err != nil {
		return nil, fmt.Errorf("failed to run command: %v", err)
	}
	return ParseUnmanagedONT(output)
}

func (c *CommandExecutor) GetOpticalInfo(port, ontID int) (*OnuOpticalInfo, error) {
	if c.ExecutorContext.Level != 3 {
		return nil, fmt.Errorf("not in config mode")
	}
	output, err := c.ExecuteCommand(fmt.Sprintf("display ont optical-info %d %d", port, ontID), fmt.Sprintf("MA5683T(config-if-gpon-%d/%d)#", c.ExecutorContext.Frame, c.ExecutorContext.Slot))
	if err != nil {
		return nil, fmt.Errorf("failed to run command: %v", err)
	}
	return ParseOnuOpticalInfo(output), nil
}

func (c *CommandExecutor) EnterInterfaceGPONMode(frame int, slot int) error {
	if c.ExecutorContext.Level != 2 {
		return fmt.Errorf("not in config mode")
	}
	_, err := c.ExecuteCommand(fmt.Sprintf("interface gpon %d/%d", frame, slot), fmt.Sprintf("MA5683T(config-if-gpon-%d/%d)#", frame, slot))
	if err != nil {
		return fmt.Errorf("failed to run command: %v", err)
	}
	c.ExecutorContext.Level = 3
	c.ExecutorContext.Frame = frame
	c.ExecutorContext.Slot = slot
	return nil
}

func (c *CommandExecutor) AddOpticalNetworkTerminal(port int, serialNumber string, description string) (int, error) {
	if c.ExecutorContext.Level != 3 {
		return 0, fmt.Errorf("not in interface gpon mode")
	}

	serialParts := strings.Split(serialNumber, " ")
	output, err := c.ExecuteCommand(fmt.Sprintf("ont add %d sn-auth %s omci ont-lineprofile-id 60 ont-srvprofile-id 35 desc %s",
		port,
		serialParts[0],
		description,
	), fmt.Sprintf("MA5683T(config-if-gpon-%d/%d)#", c.ExecutorContext.Frame, c.ExecutorContext.Slot))

	if err != nil {
		return 0, fmt.Errorf("failed to run command: %v", err)
	}

	err = c.checkOutputFailure(output)
	if err != nil {
		return 0, err
	}

	re := regexp.MustCompile(`ONTID :(\d+)`)
	match := re.FindStringSubmatch(output)
	if len(match) < 2 {
		return 0, fmt.Errorf("ONTID not found in command output")
	}

	ontID, err := strconv.Atoi(match[1])
	if err != nil {
		return 0, fmt.Errorf("failed to parse ONTID: %v", err)
	}

	return ontID, nil
}

func (c *CommandExecutor) DeleteOpticalNetworkTerminal(port int) error {
	if c.ExecutorContext.Level != 3 {
		return fmt.Errorf("not in interface gpon mode")
	}

	_, err := c.ExecuteCommand(fmt.Sprintf("ont delete %d all", port), "(y/n)[n]:")
	if err != nil {
		return fmt.Errorf("failed to run command: %v", err)
	}

	_, err = c.ExecuteCommand("y", fmt.Sprintf("MA5683T(config-if-gpon-%d/%d)#", c.ExecutorContext.Frame, c.ExecutorContext.Slot))
	if err != nil {
		return fmt.Errorf("failed to run command: %v", err)
	}
	return nil
}

func (c *CommandExecutor) AddNativeVirtualLan(port, ontID int) error {
	if c.ExecutorContext.Level != 3 {
		return fmt.Errorf("not in interface gpon mode")
	}

	output, err := c.ExecuteCommand(fmt.Sprintf("ont port native-vlan %d %d eth 1 vlan 20 priority 0", port, ontID), fmt.Sprintf("MA5683T(config-if-gpon-%d/%d)#", c.ExecutorContext.Frame, c.ExecutorContext.Slot))
	if err != nil {
		return fmt.Errorf("failed to run command: %v", err)
	}

	err = c.checkOutputFailure(output)
	if err != nil {
		return err
	}

	return nil
}

func (c *CommandExecutor) AddServicePort(vlan, frame, slot, port, ontID int) error {
	if c.ExecutorContext.Level != 2 {
		return fmt.Errorf("not in config mode")
	}

	output, err := c.ExecuteCommand(fmt.Sprintf("service-port vlan %d gpon %d/%d/%d ont %d gemport 20 multi-service user-vlan 20 tag-transform translate inbound traffic-table index 10 outbound traffic-table index 10", vlan, frame, slot, port, ontID), "MA5683T(config)#")
	if err != nil {
		return fmt.Errorf("failed to run command: %v", err)
	}

	err = c.checkOutputFailure(output)
	if err != nil {
		return err
	}

	return nil
}

func (c *CommandExecutor) GetOpticalNetworkTerminal(frame, slot, port, ontID int) (*ONT, error) {
	if c.ExecutorContext.Level != 2 {
		return nil, fmt.Errorf("not in config mode")
	}

	output, err := c.ExecuteCommand(fmt.Sprintf("display current-configuration ont %d/%d/%d %d", frame, slot, port, ontID), "MA5683T(config)#")
	if err != nil {
		return nil, fmt.Errorf("failed to run command: %v", err)
	}

	err = c.checkOutputFailure(output)
	if err != nil {
		return nil, err
	}

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
		ont.ServicePort, err = strconv.Atoi(match[1])
		if err != nil {
			return nil, fmt.Errorf("failed to parse service port: %v", err)
		}
	}

	if match := regexp.MustCompile(`vlan (\d+) gpon`).FindStringSubmatch(output); len(match) > 1 {
		ont.VlanID, err = strconv.Atoi(match[1])
		if err != nil {
			return nil, fmt.Errorf("failed to parse vlan id: %v", err)
		}
	}

	return ont, nil
}

func (c *CommandExecutor) UndoServicePort(id int) error {
	if c.ExecutorContext.Level != 2 {
		return fmt.Errorf("not in config mode")
	}

	output, err := c.ExecuteCommand(fmt.Sprintf("undo service-port %d", id), "MA5683T(config)#")
	if err != nil {
		return fmt.Errorf("failed to run command: %v", err)
	}

	err = c.checkOutputFailure(output)
	if err != nil {
		return err
	}

	return nil
}

func (c *CommandExecutor) quit(exit bool) error {
	var err error

	if c.ExecutorContext.Level >= 3 {
		_, err = c.ExecuteCommand("quit", "MA5683T(config)#")
		if err != nil {
			return fmt.Errorf("failed to run command: %v", err)
		}
		c.ExecutorContext.Level = 2
		if !exit {
			return nil
		}
	}

	if c.ExecutorContext.Level >= 2 {
		_, err = c.ExecuteCommand("quit", "MA5683T#")
		if err != nil {
			return fmt.Errorf("failed to run command: %v", err)
		}
		c.ExecutorContext.Level = 1
		if !exit {
			return nil
		}
	}

	_, err = c.ExecuteCommand("quit", "before logout")
	if err != nil {
		return fmt.Errorf("failed to run command: %v", err)
	}

	_, err = c.ExecuteCommand("y", "to log on")
	if err != nil {
		return fmt.Errorf("failed to run command: %v", err)
	}
	return nil
}

func (c *CommandExecutor) readOutputUntilPrompt(prompt string) (string, error) {
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

func (c *CommandExecutor) enable() error {
	if c.ExecutorContext.Level != 0 {
		return fmt.Errorf("not in root mode")
	}
	_, err := c.ExecuteCommand("enable", "MA5683T#")
	if err != nil {
		return fmt.Errorf("failed to run command enable: %v", err)
	}
	c.ExecutorContext.Level = 1
	return nil
}

func (c *CommandExecutor) config() error {
	if c.ExecutorContext.Level != 1 {
		return fmt.Errorf("not in enable mode")
	}
	_, err := c.ExecuteCommand("config", "MA5683T(config)#")
	if err != nil {
		return fmt.Errorf("failed to run command config: %v", err)
	}
	c.ExecutorContext.Level = 2
	return nil
}

func (c *CommandExecutor) checkOutputFailure(output string) error {
	if strings.Contains(output, "Failure:") {
		lines := strings.Split(output, "\n")
		for _, line := range lines {
			if strings.Contains(line, "Failure:") {
				return fmt.Errorf(line)
			}
		}
	}
	return nil
}

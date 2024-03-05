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
	Verbose           bool
	Stdout            io.Reader
	Stdin             io.WriteCloser
	ExecutorContext   ExecutorContext
	ConnectionManager *ConnectionManager
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
		Stdout:            stdout,
		Stdin:             stdin,
		ExecutorContext:   ExecutorContext{},
		Verbose:           options.Verbose,
		ConnectionManager: connManager,
	}

	_, err = commExecutor.readOutputUntilPrompt(">")
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

func (c *CommandExecutor) GetUnmanagedOpticalNetworkTerminals() ([]UnmanagedONT, error) {
	if c.ExecutorContext.Level != 2 {
		return nil, fmt.Errorf("not in config mode")
	}
	output, err := c.ExecuteCommand("display ont autofind all\n", "(config)#")
	if err != nil {
		return nil, fmt.Errorf("failed to run command: %v", err)
	}
	return ParseUnmanagedONT(output)
}

func (c *CommandExecutor) GetOpticalInfo(port, ontID int) (*OpticalInfo, error) {
	if c.ExecutorContext.Level != 3 {
		return nil, fmt.Errorf("not in config mode")
	}
	output, err := c.ExecuteCommand(fmt.Sprintf("display ont optical-info %d %d", port, ontID), fmt.Sprintf("(config-if-gpon-%d/%d)#", c.ExecutorContext.Frame, c.ExecutorContext.Slot))
	if err != nil {
		return nil, fmt.Errorf("failed to run command: %v", err)
	}
	return ParseOpticalInfo(output)
}

func (c *CommandExecutor) GetGeneralInfoBySn(sn string) (*GeneralInfo, error) {
	if c.ExecutorContext.Level != 2 {
		return nil, fmt.Errorf("not in config mode")
	}
	output, err := c.ExecuteCommand(fmt.Sprintf("display ont info by-sn %s", strings.Split(sn, " ")[0]), "(config)#")
	if err != nil {
		return nil, fmt.Errorf("failed to run command: %v", err)
	}
	return ParseGeneralInfoBySn(output)
}

func (c *CommandExecutor) GetServicePorts(frame, slot, port, ontID int) ([]ServicePort, error) {
	if c.ExecutorContext.Level != 2 {
		return nil, fmt.Errorf("not in config mode")
	}
	output, err := c.ExecuteCommand(fmt.Sprintf("display service-port port %d/%d/%d ont %d\n", frame, slot, port, ontID), "(config)#")
	if err != nil {
		return nil, fmt.Errorf("failed to run command: %v", err)
	}
	return ParseServicePorts(output)
}

func (c *CommandExecutor) EnterInterfaceGPONMode(frame int, slot int) error {
	if c.ExecutorContext.Level != 2 {
		return fmt.Errorf("not in config mode")
	}
	_, err := c.ExecuteCommand(fmt.Sprintf("interface gpon %d/%d", frame, slot), fmt.Sprintf("(config-if-gpon-%d/%d)#", frame, slot))
	if err != nil {
		return fmt.Errorf("failed to run command: %v", err)
	}
	c.ExecutorContext.Level = 3
	c.ExecutorContext.Frame = frame
	c.ExecutorContext.Slot = slot
	return nil
}

func (c *CommandExecutor) AddOpticalNetworkTerminal(port int, sn string, description string) (int, error) {
	if c.ExecutorContext.Level != 3 {
		return 0, fmt.Errorf("not in interface gpon mode")
	}

	output, err := c.ExecuteCommand(fmt.Sprintf("ont add %d sn-auth %s omci ont-lineprofile-id 60 ont-srvprofile-id 35 desc %s",
		port,
		strings.Split(sn, " ")[0],
		description,
	), fmt.Sprintf("(config-if-gpon-%d/%d)#", c.ExecutorContext.Frame, c.ExecutorContext.Slot))

	if err != nil {
		return 0, fmt.Errorf("failed to run command: %v", err)
	}

	lines := strings.Split(output, "\n")
	err = parseLinesFailure(lines)
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

	_, err = c.ExecuteCommand("y", fmt.Sprintf("(config-if-gpon-%d/%d)#", c.ExecutorContext.Frame, c.ExecutorContext.Slot))
	if err != nil {
		return fmt.Errorf("failed to run command: %v", err)
	}
	return nil
}

func (c *CommandExecutor) AddNativeVirtualLan(port, ontID int) error {
	if c.ExecutorContext.Level != 3 {
		return fmt.Errorf("not in interface gpon mode")
	}

	output, err := c.ExecuteCommand(fmt.Sprintf("ont port native-vlan %d %d eth 1 vlan 20 priority 0", port, ontID), fmt.Sprintf("(config-if-gpon-%d/%d)#", c.ExecutorContext.Frame, c.ExecutorContext.Slot))
	if err != nil {
		return fmt.Errorf("failed to run command: %v", err)
	}

	lines := strings.Split(output, "\n")
	err = parseLinesFailure(lines)
	if err != nil {
		return err
	}

	return nil
}

func (c *CommandExecutor) AddServicePort(vlan, frame, slot, port, ontID int) error {
	if c.ExecutorContext.Level != 2 {
		return fmt.Errorf("not in config mode")
	}

	output, err := c.ExecuteCommand(fmt.Sprintf("service-port vlan %d gpon %d/%d/%d ont %d gemport 20 multi-service user-vlan 20 tag-transform translate inbound traffic-table index 10 outbound traffic-table index 10", vlan, frame, slot, port, ontID), "(config)#")
	if err != nil {
		return fmt.Errorf("failed to run command: %v", err)
	}

	lines := strings.Split(output, "\n")
	err = parseLinesFailure(lines)
	if err != nil {
		return err
	}

	return nil
}

func (c *CommandExecutor) UndoServicePort(id int) error {
	if c.ExecutorContext.Level != 2 {
		return fmt.Errorf("not in config mode")
	}

	output, err := c.ExecuteCommand(fmt.Sprintf("undo service-port %d", id), "(config)#")
	if err != nil {
		return fmt.Errorf("failed to run command: %v", err)
	}

	lines := strings.Split(output, "\n")
	err = parseLinesFailure(lines)
	if err != nil {
		return err
	}

	return nil
}

func (c *CommandExecutor) quit(exit bool) error {
	var err error

	if c.ExecutorContext.Level >= 3 {
		_, err = c.ExecuteCommand("quit", "(config)#")
		if err != nil {
			return fmt.Errorf("failed to run command: %v", err)
		}
		c.ExecutorContext.Level = 2
		if !exit {
			return nil
		}
	}

	if c.ExecutorContext.Level >= 2 {
		_, err = c.ExecuteCommand("quit", "#")
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

	if err != nil {
		err := c.ConnectionManager.Close()
		if err != nil {
			fmt.Println("Failed to close connection: ", err)
		}
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
	_, err := c.ExecuteCommand("enable", "#")
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
	_, err := c.ExecuteCommand("config", "(config)#")
	if err != nil {
		return fmt.Errorf("failed to run command config: %v", err)
	}
	c.ExecutorContext.Level = 2
	return nil
}

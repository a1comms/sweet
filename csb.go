package sweet

import (
	"fmt"
	"strings"
)
// Cisco SBS (Built for Cisco SG300)
type Csb struct {
}

func newCsbCollector() Collector {
	return Csb{}
}

func (collector Csb) Collect(device DeviceConfig) (map[string]string, error) {
	result := make(map[string]string)

	c, err := newSSHCollector(device)
	if err != nil {
		return result, fmt.Errorf("Error connecting to %s: %s", device.Hostname, err.Error())
	}
    if err := expect("ser Name:", c.Receive); err != nil {
        return result, fmt.Errorf("Missing Username prompt: %s", err.Error())
    }
    c.Send <- device.Config["user"] + "\n"
	if err := expect("assword:", c.Receive); err != nil {
		return result, fmt.Errorf("Missing password prompt: %s", err.Error())
	}
	c.Send <- device.Config["pass"] + "\n"
	multi := []string{"#", ">", "assword:"}
	m, err := expectMulti(multi, c.Receive)
	if err != nil {
		return result, fmt.Errorf("Invalid response to password: %s", err.Error())
	}
	if m == "assword:" {
		return result, fmt.Errorf("Bad username or password.")
    }    

// Turn off terminal paging	

	c.Send <- "terminal datadump\n"
	if err := expect("#", c.Receive); err != nil {
		return result, fmt.Errorf("Command 'terminal datadump' failed: %s", err.Error())
	}
// Dump config
    c.Send <- "show running-config\n"
	result["config"], err = expectSaveTimeout("#", c.Receive, device.CommandTimeout)
	if err != nil {
		return result, fmt.Errorf("Command 'show running-config' failed: %s", err.Error())
	}
// Dump Version    
	c.Send <- "show version\n"
	result["version"], err = expectSaveTimeout("#", c.Receive, device.CommandTimeout)
	if err != nil {
		return result, fmt.Errorf("Command 'show version' failed: %s", err.Error())
	}

	// cleanup config results
	result["config"] = strings.TrimSpace(strings.TrimPrefix(result["config"], "show running-config"))
	result["config"] = strings.TrimSpace(strings.TrimPrefix(result["config"], "config-file-header"))

	c.Send <- "exit\n"

	return result, nil
}

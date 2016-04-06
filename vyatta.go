package sweet

import (
	"fmt"
)
// Vyatta - build around Ubiquiti EdgeOS
type Vyatta struct {
}

func newVyattaCollector() Collector {
	return Vyatta{}
}

func (collector Vyatta) Collect(device DeviceConfig) (map[string]string, error) {
	result := make(map[string]string)

	c, err := newSSHCollector(device)
	if err != nil {
		return result, fmt.Errorf("Error connecting to %s: %s", device.Hostname, err.Error())
	}
	if err := expect("assword:", c.Receive); err != nil {
		return result, fmt.Errorf("Missing password prompt: %s", err.Error())
	}
	c.Send <- device.Config["pass"] + "\n"
	multi := []string{"$", "assword:"}
	m, err := expectMulti(multi, c.Receive)
	if err != nil {
		return result, fmt.Errorf("Invalid response to password: %s", err.Error())
	}
	if m == "assword:" {
		return result, fmt.Errorf("Bad username or password.")
    }    

// Disable terminal pager
	c.Send <- "terminal length 0\n"
	if err := expect("$", c.Receive); err != nil {
		return result, fmt.Errorf("Command 'terminal length 0' failed: %s", err.Error())
	}
// Dump config
    c.Send <- "show configuration\n"
	result["config"], err = expectSaveTimeout("$", c.Receive, device.CommandTimeout)
	if err != nil {
		return result, fmt.Errorf("Command 'show configuration' failed: %s", err.Error())
	}

	c.Send <- "exit\n"

	return result, nil
}

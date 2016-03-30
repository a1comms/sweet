package sweet

import (
	"fmt"
//	"strings"
)
// HP/Comware 
type Cmw struct {
}

func newCmwCollector() Collector {
	return Cmw{}
}

func (collector Cmw) Collect(device DeviceConfig) (map[string]string, error) {
	result := make(map[string]string)

	c, err := newSSHCollector(device)
	if err != nil {
		return result, fmt.Errorf("Error connecting to %s: %s", device.Hostname, err.Error())
	}
	if err := expect("assword:", c.Receive); err != nil {
		return result, fmt.Errorf("Missing password prompt: %s", err.Error())
	}
	c.Send <- device.Config["pass"] + "\n"
	multi := []string{"#", ">", "assword:", "[Y/N]:", "osition."}
	m, err := expectMulti(multi, c.Receive)
	if err != nil {
		return result, fmt.Errorf("Invalid response to password: %s", err.Error())
	}
	if m == "assword:" {
		return result, fmt.Errorf("Bad username or password.")
    } else if m == ">" {
// Check for comware unlock pass - unlock if present    
        if len(device.Config["comware-unlock-pass"]) > 0 {
            c.Send <- "xtd-cli-mode\n"
            if err := expect("Y/N]:", c.Receive); err != nil {
                return result, fmt.Errorf("No Y/N prompt")
            }
            c.Send <- "y\n"
            if err := expect("assword:", c.Receive); err !=nil{
                return result, fmt.Errorf("Missing password prompt for Extented CLI")
            }
            c.Send <- device.Config["comware-unlock-pass"] + "\n"
            if err := expect(">", c.Receive); err !=nil{
                return result, fmt.Errorf("Problem with CLI unlock password")
            }
        }
// Turn off terminal paging
        c.Send <- "screen-length disable\n"
        if err := expect(">", c.Receive); err !=nil {
            return result, fmt.Errorf("Unable to disable pager")
        }    
        c.Send <- "system-view\n"
        if err := expect("]", c.Receive); err !=nil{
            return result, fmt.Errorf("Unable to enter system-view")
        }
// Dump config
        c.Send <- "display current-configuration\n"
        result["config"], err = expectSaveTimeout("]", c.Receive, device.CommandTimeout)
        if err != nil {
            return result, fmt.Errorf("Command 'display current-configuration' failed: %s", err.Error())
        }
        
   }

	c.Send <- "quit\n"

	return result, nil
}

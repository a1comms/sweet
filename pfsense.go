package sweet

import (
	"fmt"
    "regexp"
    "strings"
    "time"
)
// pfSense 
type Pfs struct {
}

func newPfsCollector() Collector {
	return Pfs{}
}

func (collector Pfs) Collect(device DeviceConfig) (map[string]string, error) {
	result := make(map[string]string)
    tail := regexp.MustCompile("(?m)[\r\n]+^.*-RELEASE.*$")
	c, err := newSSHCollector(device)
	if err != nil {
		return result, fmt.Errorf("Error connecting to %s: %s", device.Hostname, err.Error())
	}else {
    }
	if err := expect("Password", c.Receive); err != nil {
		return result, fmt.Errorf("Missing password prompt: %s", err.Error())
	}else {
    }    
	c.Send <- device.Config["pass"] + "\n"
	multi := []string{"option:", "root:", "Password", }
	m, err := expectMulti(multi, c.Receive)
	if err != nil {
		return result, fmt.Errorf("Invalid response to password: %s", err.Error())
	}
	if m == "Password" {
		return result, fmt.Errorf("Bad username or password.")
    }else if m == "option:" {
        c.Send <- "8\n"
        if err := expect("-RELEASE", c.Receive); err !=nil {
            return result, fmt.Errorf("Unable to activate Shell")
        }
// Dump config
        var conf string
        conf = "cat /conf/config.xml"
        c.Send <- conf + "\n" 
        result["config"], err = expectSaveTimeout("-RELEASE", c.Receive, device.CommandTimeout)
        if err != nil {
            return result, fmt.Errorf("Unable to dump config.xml", err.Error())
        }
        
   }
    result["config"] = strings.TrimSpace(strings.TrimPrefix(result["config"], "cat /conf/config.xml\r"))
    result["config"] = tail.ReplaceAllString(result["config"], "")
	c.Send <- "exit\n"
    time.Sleep(1 * time.Second)
    c.Send <- "0\n"
	return result, nil
}

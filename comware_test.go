package sweet

import (
	"os"
	"strings"
	"testing"
	"time"
)

func TestCmwGood(t *testing.T) {
	d := new(DeviceConfig)
	d.Config = make(map[string]string)
	d.Method = "cmw"
	d.Timeout = 10 * time.Second

	if os.Getenv("SWEET_COMWARE_HOST") == "" {
		t.Error("Test requries SWEET_COMWARE_HOST environment variable")
		return
	}
	if os.Getenv("SWEET_COMWARE_USER") == "" {
		t.Error("Test requries SWEET_COMWARE_USER environment variable")
		return
	}
	if os.Getenv("SWEET_COMWARE_PASS") == "" {
		t.Error("Test requries SWEET_COMWARE_PASS environment variable")
		return
	}

	d.Hostname = os.Getenv("SWEET_COMWARE_HOST")
	d.Config["user"] = os.Getenv("SWEET_COMWARE_USER")
	d.Config["pass"] = os.Getenv("SWEET_COMWARE_PASS")

	d.Target = d.Hostname

	s := CollectCmw(*d)
	if !strings.Contains(s["config"], "aaa authorization commands") {
		t.Errorf("Config missing aaa line")
	}
	if !strings.Contains(s["config"], "ntp access-group peer") {
		t.Errorf("Config missing ntp line close to end")
	}

}

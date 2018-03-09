package sweet

import (
	"os"
	"strings"
	"testing"
	"time"
)

func TestPfsGood(t *testing.T) {
	d := new(DeviceConfig)
	d.Config = make(map[string]string)
	d.Method = "pfsense"
	d.Timeout = 10 * time.Second

	if os.Getenv("SWEET_PFSENSE_HOST") == "" {
		t.Error("Test requries SWEET_PFSENSE_HOST environment variable")
		return
	}
	if os.Getenv("SWEET_PFSENSE_USER") == "" {
		t.Error("Test requries SWEET_PFSENSE_USER environment variable")
		return
	}
	if os.Getenv("SWEET_PFSENSE_PASS") == "" {
		t.Error("Test requries SWEET_PFSENSE_PASS environment variable")
		return
	}

	d.Hostname = os.Getenv("SWEET_PFSENSE_HOST")
	d.Config["user"] = os.Getenv("SWEET_PFSENSE_USER")
	d.Config["pass"] = os.Getenv("SWEET_PFSENSE_PASS")

	d.Target = d.Hostname

	s := CollectPfs(*d)
	if !strings.Contains(s["config"], "aaa authorization commands") {
		t.Errorf("Config missing aaa line")
	}
	if !strings.Contains(s["config"], "ntp access-group peer") {
		t.Errorf("Config missing ntp line close to end")
	}

}

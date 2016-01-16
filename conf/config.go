package conf

import (
	"encoding/json"
	"fmt"
	"github.com/golang/glog"
	"github.com/samitpal/goProbe/modules"
	"github.com/samitpal/goProbe/modules/http"
	"github.com/samitpal/goProbe/modules/ping_port"
)

/* Example json config
[
    {
        "probe_type": "http",
        "probe_config": {
            "probe_name": "example_http",
            "probe_url": "http://example.com"

        }
    },
    {
	    "probe_type": "ping_port",
	    "probe_config": {
	        "probe_name": "example_ping_port_80",
	        "probe_host_name": "example.com",
	        "probe_host_port": 80
	        }
    }
]
*/

type Probes struct {
	ProbeType   string          `json:"probe_type"`
	ProbeConfig json.RawMessage `json:"probe_config"` // Here we branch to the respective probe type config.
}

func SetupConfig(config []byte) ([]modules.Prober, error) {
	var p []Probes
	err := json.Unmarshal(config, &p)
	if err != nil {
		return nil, err
	}

	var probes []modules.Prober
	// TODO: It is a bit ugly here. A new probe module will need its own 'if/else if' block. Make it better.
	for _, c := range p {
		if c.ProbeType == "http" {
			t := http.NewHttpProbe()
			err := json.Unmarshal(c.ProbeConfig, t)
			if err != nil {
				return nil, err
			}
			// Call the module's Prepare method which should do its own initialization (if any).
			err = t.Prepare()
			if err == nil {
				probes = append(probes, t)
			} else {
				glog.Errorf("Error in config: %v", err)
			}
		} else if c.ProbeType == "ping_port" {
			t := ping_port.NewPingPortProbe()
			err := json.Unmarshal(c.ProbeConfig, t)
			if err != nil {
				return nil, err
			}
			// Call the module's Prepare method which should do its own initialization (if any).
			err = t.Prepare()
			if err == nil {
				probes = append(probes, t)
			} else {
				glog.Errorf("Error in config: %v", err)
			}
		}
		// Add a new 'else if' statement here for a new probe type.
	}
	if err = checkDuplicateProbeNames(probes); err != nil {
		return nil, err
	}
	return probes, nil
}

// checkDuplicateProbeNames checks for duplicate probe names.
func checkDuplicateProbeNames(pms []modules.Prober) error {
	probeCount := make(map[string]int)
	for _, pm := range pms {
		probeCount[*pm.Name()]++
		if probeCount[*pm.Name()] > 1 {
			return fmt.Errorf("Duplicate probe name '%s' found.", *pm.Name())
		}
	}
	return nil
}

func GetProbeNames(pms []modules.Prober) []string {
	var probeNames []string
	for _, pm := range pms {
		probeNames = append(probeNames, *pm.Name())
	}
	return probeNames
}

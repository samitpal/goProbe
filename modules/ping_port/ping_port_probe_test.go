package ping_port

import "testing"

func TestCheckConfig(t *testing.T) {

	pn := "probe1"
	phn := "example.com"
	phr := 80
	pnw := "invalid"

	// test without probe name.
	pm1 := NewPingPortProbe()
	pm1.ProbeHostName = &phn
	pm1.ProbeHostPort = &phr
	err := pm1.checkConfig()
	if err == nil {
		t.Errorf("Probe name is mandatory. Test expected to fail but is passing")
	}
	// test without probe host name.
	pm2 := NewPingPortProbe()
	pm2.ProbeName = &pn
	pm2.ProbeHostPort = &phr
	err = pm2.checkConfig()
	if err == nil {
		t.Errorf("Probe host name is mandatory. Test expected to fail but is passing")
	}

	// test without probe host port.
	pm3 := NewPingPortProbe()
	pm3.ProbeName = &pn
	pm3.ProbeHostName = &phn
	err = pm3.checkConfig()
	if err == nil {
		t.Errorf("Probe port is mandatory. Test expected to fail but is passing")
	}

	// test with invalid network.
	pm4 := NewPingPortProbe()
	pm4.ProbeName = &pn
	pm4.ProbeHostName = &phn
	pm4.ProbeHostPort = &phr
	pm4.ProbeNetwork = &pnw
	err = pm4.checkConfig()
	if err == nil {
		t.Errorf("Probe network is invalid. Test expected to fail but is passing")
	}

}

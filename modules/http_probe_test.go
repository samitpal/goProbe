package modules

import "testing"

func TestCheckConfig(t *testing.T) {

	pn := "probe1"
	pu := "http://example.com"
	pa := "check_match_payload"
	pm := "invalid_method"

	// test without probe url.
	hm1 := NewHttpProbe()
	hm1.ProbeName = &pn
	err := hm1.checkConfig()
	if err == nil {
		t.Errorf("Probe url is mandatory. Test expected to fail but is passing")
	}
	// test without probe name.
	hm2 := NewHttpProbe()
	hm2.ProbeURL = &pu
	err = hm2.checkConfig()
	if err == nil {
		t.Errorf("Probe name is mandatory. Test expected to fail but is passing")
	}

	// test with action set to check_match_payload but not match string.
	hm3 := NewHttpProbe()
	hm3.ProbeName = &pn
	hm3.ProbeURL = &pu
	hm3.ProbeAction = &pa
	err = hm3.checkConfig()
	if err == nil {
		t.Errorf("Probe match string is mandatory. Test expected to fail but is passing")
	}

	// test with invalid http method.
	hm4 := NewHttpProbe()
	hm4.ProbeName = &pn
	hm4.ProbeURL = &pu
	hm4.ProbeHttpMethod = &pm
	err = hm4.checkConfig()
	if err == nil {
		t.Errorf("Probe http method is invalid. Test expected to fail but is passing")
	}

}

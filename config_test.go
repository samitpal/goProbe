package main

import "testing"

func TestSetupConfig(t *testing.T) {
	// config with 2 valid and one invalid probes.
	config := []byte(`
		[
    	{
        "probe_type": "http",
        "probe_config": {
            "probe_name": "probe1",
            "probe_url": "http://example.com",
            "probe_timeout": 20,
            "probe_interval": 30

        }	
    	},
    	{
        "probe_type": "http",
        "probe_config": {
            "probe_name": "probe2",
            "probe_url": "https://example.com",
            "probe_action": "check_match_payload",
            "probe_match_string": "match_me"
        }
    	},
    	{
        "probe_type": "http",
        "probe_config": {
            "probe_name": "probe3",
            "probe_url": "https://example.com",
  			"probe_http_method": "INVALID"
        }
    	}
		]`)

	got, _ := SetupConfig(config)

	// Test that there are two elements.
	if len(got) != 2 {
		t.Error("Element lentgth should be two")
	}

	//Test 1
	b1 := []byte(
		`{
 "probe_name": "probe1",
 "probe_url": "http://example.com",
 "probe_http_method": "GET",
 "probe_action": "check_ret_200",
 "probe_match_string": null,
 "probe_interval": 30,
 "probe_timeout": 20
}`)
	b2 := []byte(
		`{
 "probe_name": "probe2",
 "probe_url": "https://example.com",
 "probe_http_method": "GET",
 "probe_action": "check_match_payload",
 "probe_match_string": "match_me",
 "probe_interval": 60,
 "probe_timeout": 40
}`)
	want1 := []string{string(b1), string(b2)}
	for i, _ := range got {
		if want1[i] != got[i].RetConfig() {
			t.Errorf("Got: \n%v\n Want: \n%v", got[i].RetConfig(), want1[i])
		}
	}

	//Test 2
	p1 := "probe1"
	p2 := "probe2"
	want2 := []*string{&p1, &p2}
	for i, _ := range got {
		if *want2[i] != *got[i].Name() {
			t.Errorf("Got: %v\n Want: %v", *got[i].Name(), *want2[i])
		}
	}

	//Test 3
	t1 := 20
	t2 := 40
	want3 := []*int{&t1, &t2}
	for i, _ := range got {
		if *want3[i] != *got[i].TimeoutSecs() {
			t.Errorf("Got: %v\n Want: %v", *got[i].TimeoutSecs(), *want3[i])
		}
	}

	// Test 4
	i1 := 30
	i2 := 60
	want4 := []*int{&i1, &i2}
	for i, _ := range got {
		if *want4[i] != *got[i].RunIntervalSecs() {
			t.Errorf("Got: %v\n Want: %v", *got[i].RunIntervalSecs(), *want4[i])
		}
	}
}

package main

import "testing"

func TestSetupConfig(t *testing.T) {
	// config with 2 valid and one invalid http probes. It also has a valid ping_port probe.
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
            "probe_match_string": "match_me",
			"probe_http_headers": {
                "host": "blah",
                "user_agent": "test-agent"

            }            
        }
    	},
    	{
        "probe_type": "http",
        "probe_config": {
            "probe_name": "probe3",
            "probe_url": "https://example.com",
  			"probe_http_method": "INVALID"
        }
        },
        {
        "probe_type": "ping_port",
        "probe_config": {
            "probe_name": "probe4",
            "probe_host_name": "example.com",
            "probe_host_port": 22
        }
    	}
		]`)

	got, _ := setupConfig(config)

	// Test that there are three elements.
	if len(got) != 3 {
		t.Error("Element lentgth should be three")
	}

	//Test 1
	b1 := []byte(
		`{
 "probe_name": "probe1",
 "probe_url": "http://example.com",
 "probe_http_method": "GET",
 "probe_action": "check_ret_200",
 "probe_match_string": null,
 "probe_http_headers": null,
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
 "probe_http_headers": {
  "host": "blah",
  "user_agent": "test-agent"
 },
 "probe_interval": 60,
 "probe_timeout": 40
}`)
	b3 := []byte(
		`{
 "probe_name": "probe4",
 "probe_interval": 60,
 "probe_timeout": 10,
 "probe_host_name": "example.com",
 "probe_host_port": 22,
 "probe_network": "tcp"
}`)

	want1 := []string{string(b1), string(b2), string(b3)}
	for i, _ := range got {
		if want1[i] != got[i].RetConfig() {
			t.Errorf("Got: \n%v\n Want: \n%v", got[i].RetConfig(), want1[i])
		}
	}

	//Test 2
	p1 := "probe1"
	p2 := "probe2"
	p4 := "probe4"
	want2 := []*string{&p1, &p2, &p4}
	for i, _ := range got {
		if *want2[i] != *got[i].Name() {
			t.Errorf("Got: %v\n Want: %v", *got[i].Name(), *want2[i])
		}
	}

	//Test 3
	t1 := 20
	t2 := 40
	t4 := 10
	want3 := []*int{&t1, &t2, &t4}
	for i, _ := range got {
		if *want3[i] != *got[i].TimeoutSecs() {
			t.Errorf("Got: %v\n Want: %v", *got[i].TimeoutSecs(), *want3[i])
		}
	}

	// Test 4
	i1 := 30
	i2 := 60
	i4 := 60
	want4 := []*int{&i1, &i2, &i4}
	for i, _ := range got {
		if *want4[i] != *got[i].RunIntervalSecs() {
			t.Errorf("Got: %v\n Want: %v", *got[i].RunIntervalSecs(), *want4[i])
		}
	}
}

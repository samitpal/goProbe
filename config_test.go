package main

import (
	"reflect"
	"testing"
)

func TestSetupConfig(t *testing.T) {
	// config with 3 valid and 1 invalid http probes. It also has a valid ping_port probe.
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
            "probe_action": "check_sslcert_expiry"
        }	
    	},
    	{
        "probe_type": "http",
        "probe_config": {
            "probe_name": "probe4",
            "probe_url": "https://example.com",
  			"probe_http_method": "INVALID"
        }
        },
        {
        "probe_type": "ping_port",
        "probe_config": {
            "probe_name": "probe5",
            "probe_host_name": "example.com",
            "probe_host_port": 22
        }
    	}
		]`)

	got, _ := setupConfig(config)

	// Test that there are three elements.
	if len(got) != 4 {
		t.Error("Element lentgth should be four")
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
 "probe_sslcert_expires_in_days": null,
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
 "probe_sslcert_expires_in_days": null,
 "probe_interval": 60,
 "probe_timeout": 40
}`)
	b3 := []byte(
		`{
 "probe_name": "probe3",
 "probe_url": "https://example.com",
 "probe_http_method": "GET",
 "probe_action": "check_sslcert_expiry",
 "probe_match_string": null,
 "probe_http_headers": null,
 "probe_sslcert_expires_in_days": 30,
 "probe_interval": 60,
 "probe_timeout": 40
}`)
	b5 := []byte(
		`{
 "probe_name": "probe5",
 "probe_interval": 60,
 "probe_timeout": 10,
 "probe_host_name": "example.com",
 "probe_host_port": 22,
 "probe_network": "tcp"
}`)

	want1 := []string{string(b1), string(b2), string(b3), string(b5)}
	for i, _ := range got {
		if want1[i] != got[i].RetConfig() {
			t.Errorf("Got: \n%v\n Want: \n%v", got[i].RetConfig(), want1[i])
		}
	}

	//Test 2
	p1 := "probe1"
	p2 := "probe2"
	p3 := "probe3"
	p5 := "probe5"
	want2 := []*string{&p1, &p2, &p3, &p5}
	for i, _ := range got {
		if *want2[i] != *got[i].Name() {
			t.Errorf("Got: %v\n Want: %v", *got[i].Name(), *want2[i])
		}
	}

	//Test 3
	t1 := 20
	t2 := 40
	t3 := 40
	t5 := 10
	want3 := []*int{&t1, &t2, &t3, &t5}
	for i, _ := range got {
		if *want3[i] != *got[i].TimeoutSecs() {
			t.Errorf("Got: %v\n Want: %v", *got[i].TimeoutSecs(), *want3[i])
		}
	}

	// Test 4
	i1 := 30
	i2 := 60
	i3 := 60
	i5 := 60
	want4 := []*int{&i1, &i2, &i3, &i5}
	for i, _ := range got {
		if *want4[i] != *got[i].RunIntervalSecs() {
			t.Errorf("Got: %v\n Want: %v", *got[i].RunIntervalSecs(), *want4[i])
		}
	}
}

func TestCheckDuplicateProbeNames(t *testing.T) {
	// config with duplicate probe names
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
            "probe_name": "probe1",
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
            "probe_action": "check_sslcert_expiry"
        }	
    	}]`)

	_, err := setupConfig(config)
	if err == nil {
		t.Error("Expecting error due to duplicate probe names, but test is passing")
	}
}

func TestGetProbeNames(t *testing.T) {
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
        }
        ]`)
	probes, _ := setupConfig(config)
	probeNames := getProbeNames(probes)
	if !reflect.DeepEqual(probeNames, []string{"probe1", "probe2"}) {
		t.Errorf("Got: %v\n Want: %v", probeNames, []string{"probe1", "probe2"})
	}
}

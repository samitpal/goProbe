package misc

import (
	"github.com/samitpal/goProbe/modules"
	"testing"
)

func TestCheckProbeConfig(t *testing.T) {
	p := []modules.Prober{}
	err := CheckProbeConfig(p)
	if err == nil {
		t.Error("Expected an error to be returned")
	}
}

func TestCheckProbeData(t *testing.T) {
	up := float64(0)
	latency := float64(2)
	startTime := int64(1234)
	endTime := int64(5678)

	// Test 1
	pd := &modules.ProbeData{Latency: &latency, StartTime: &startTime, EndTime: &endTime}
	err := CheckProbeData(pd)
	if err == nil {
		t.Errorf("Expected an error to be returned %v", err)

	}

	// Test 2
	pd = &modules.ProbeData{IsUp: &up, StartTime: &startTime, EndTime: &endTime}
	err = CheckProbeData(pd)
	if err == nil {
		t.Errorf("Expected an error to be returned %v", err)

	}

	// Test 3
	pd = &modules.ProbeData{IsUp: &up, Latency: &latency, EndTime: &endTime}
	err = CheckProbeData(pd)
	if err == nil {
		t.Errorf("Expected an error to be returned %v", err)

	}

	// Test 4
	pd = &modules.ProbeData{IsUp: &up, Latency: &latency, StartTime: &startTime}
	err = CheckProbeData(pd)
	if err == nil {
		t.Errorf("Expected an error to be returned %v", err)

	}
}

package metric_export

import (
	"github.com/samitpal/goProbe/modules"
	"reflect"
	"testing"
	"time"
)

func TestSetFieldValues(t *testing.T) {
	up := float64(1)
	ps := float64(45)
	st := int64(567)
	et := int64(890)
	lt := float64(123)

	pd := modules.ProbeData{
		IsUp:        &up,
		Latency:     &lt,
		StartTime:   &st,
		EndTime:     &et,
		PayloadSize: &ps,
	}
	pn := "probe1"

	je := NewJSONExport()
	epochTime := time.Now().Unix()
	je.SetFieldValues(pn, &pd, epochTime)

	probe1Up := map[string]TimeValue{"probe1": TimeValue{1, epochTime}}
	if !reflect.DeepEqual(probe1Up, je.ProbeIsUp.Up) {
		t.Errorf("Got: %v\n Want: %v", je.ProbeIsUp.Up, probe1Up)
	}

	probe1Latency := map[string]TimeValue{"probe1": TimeValue{123, epochTime}}
	if !reflect.DeepEqual(probe1Latency, je.ProbeLatency.Latency) {
		t.Errorf("Got: %v\n Want: %v", je.ProbeLatency.Latency, probe1Latency)
	}

	probe1PayloadSize := map[string]TimeValue{"probe1": TimeValue{45, epochTime}}
	if !reflect.DeepEqual(probe1PayloadSize, je.ProbePayloadSize.Payload) {
		t.Errorf("Got: %v\n Want: %v", je.ProbePayloadSize, probe1PayloadSize)
	}

}

func TestSetFieldValuesUnexpected(t *testing.T) {
	pn := "probe1"
	je := NewJSONExport()
	epochTime := time.Now().Unix()
	je.SetFieldValuesUnexpected(pn, epochTime)

	probe1Up := map[string]TimeValue{"probe1": TimeValue{-1, epochTime}}
	if !reflect.DeepEqual(probe1Up, je.ProbeIsUp.Up) {
		t.Errorf("Got: %v\n Want: %v", je.ProbeIsUp.Up, probe1Up)
	}

	probe1Latency := map[string]TimeValue{"probe1": TimeValue{-1, epochTime}}
	if !reflect.DeepEqual(probe1Latency, je.ProbeLatency.Latency) {
		t.Errorf("Got: %v\n Want: %v", je.ProbeLatency.Latency, probe1Latency)
	}

	probe1PayloadSize := map[string]TimeValue{"probe1": TimeValue{-1, epochTime}}
	if !reflect.DeepEqual(probe1PayloadSize, je.ProbePayloadSize.Payload) {
		t.Errorf("Got: %v\n Want: %v", je.ProbePayloadSize, probe1PayloadSize)
	}
}

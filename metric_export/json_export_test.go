package metric_export

import (
	"github.com/samitpal/goProbe/modules"
	"reflect"
	"testing"
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
	je.SetFieldValues(pn, &pd)

	probe1Up := map[string]float64{"probe1": 1}
	if !reflect.DeepEqual(probe1Up, je.ProbeIsUp.Up) {
		t.Errorf("Got: %v\n Want: %v", je.ProbeIsUp.Up, probe1Up)
	}

	probe1Latency := map[string]float64{"probe1": 123}
	if !reflect.DeepEqual(probe1Latency, je.ProbeLatency.Latency) {
		t.Errorf("Got: %v\n Want: %v", je.ProbeLatency.Latency, probe1Latency)
	}

	probe1PayloadSize := map[string]float64{"probe1": 45}
	if !reflect.DeepEqual(probe1PayloadSize, je.ProbePayloadSize.Payload) {
		t.Errorf("Got: %v\n Want: %v", je.ProbePayloadSize, probe1PayloadSize)
	}

}

func TestSetFieldValuesUnexpected(t *testing.T) {
	pn := "probe1"
	je := NewJSONExport()
	je.SetFieldValuesUnexpected(pn)

	probe1Up := map[string]float64{"probe1": -1}
	if !reflect.DeepEqual(probe1Up, je.ProbeIsUp.Up) {
		t.Errorf("Got: %v\n Want: %v", je.ProbeIsUp.Up, probe1Up)
	}

	probe1Latency := map[string]float64{"probe1": -1}
	if !reflect.DeepEqual(probe1Latency, je.ProbeLatency.Latency) {
		t.Errorf("Got: %v\n Want: %v", je.ProbeLatency.Latency, probe1Latency)
	}

	probe1PayloadSize := map[string]float64{"probe1": -1}
	if !reflect.DeepEqual(probe1PayloadSize, je.ProbePayloadSize.Payload) {
		t.Errorf("Got: %v\n Want: %v", je.ProbePayloadSize, probe1PayloadSize)
	}
}

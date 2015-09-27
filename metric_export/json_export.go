package metric_export

import (
	"encoding/json"
	"github.com/samitpal/goProbe/modules"
	"net/http"
	"sync"
)

type TimeValue struct {
	Value float64 `json:"value"`
	Time  int64   `json:"time"`
}

type ProbeCount struct {
	sync.RWMutex
	Count map[string]TimeValue `json:"probe_count"`
}

type ProbeErrorCount struct {
	sync.RWMutex
	ErrorCount map[string]TimeValue `json:"probe_error_count"`
}

type ProbeTimeoutCount struct {
	sync.RWMutex
	TimeoutCount map[string]TimeValue `json:"probe_timeout_count"`
}

type ProbeIsUp struct {
	sync.RWMutex
	Up map[string]TimeValue `json:"probe_is_up"` // value of 1 is a success while 0 is a failure.
}

type ProbeLatency struct {
	sync.RWMutex
	Latency map[string]TimeValue `json:"probe_latency"`
}

type ProbePayloadSize struct {
	sync.RWMutex
	Payload map[string]TimeValue `json:"probe_payload_size"`
}

type jsonExport struct {
	ProbeCount
	ProbeErrorCount   // error count indicates error in probe module.
	ProbeTimeoutCount // timeout count increases when a probe times out.
	ProbeIsUp         // value of 1 is success, 0 is failure. -value of -1 could be because of probe module failure/timeout.
	ProbeLatency      // latency in milli seconds.
	ProbePayloadSize  // size of the response payload.
}

func NewJSONExport() *jsonExport {
	return &jsonExport{
		ProbeCount:        ProbeCount{Count: make(map[string]TimeValue)},
		ProbeErrorCount:   ProbeErrorCount{ErrorCount: make(map[string]TimeValue)},
		ProbeTimeoutCount: ProbeTimeoutCount{TimeoutCount: make(map[string]TimeValue)},
		ProbeIsUp:         ProbeIsUp{Up: make(map[string]TimeValue)},
		ProbeLatency:      ProbeLatency{Latency: make(map[string]TimeValue)},
		ProbePayloadSize:  ProbePayloadSize{Payload: make(map[string]TimeValue)},
	}

}

func (pm *jsonExport) Prepare() {
	// Nothing to do.
}

func (pm *jsonExport) IncProbeCount(s string, t int64) {
	pm.ProbeCount.Lock()
	var val float64
	_, ok := pm.ProbeCount.Count[s]
	if ok {
		val = pm.ProbeCount.Count[s].Value + 1
	} else {
		val = 1
	}
	pm.ProbeCount.Count[s] = TimeValue{Value: val, Time: t}
	pm.ProbeCount.Unlock()
}

func (pm *jsonExport) IncProbeErrorCount(s string, t int64) {
	pm.ProbeErrorCount.Lock()
	var val float64
	_, ok := pm.ProbeErrorCount.ErrorCount[s]
	if ok {
		val = pm.ProbeErrorCount.ErrorCount[s].Value + 1
	} else {
		val = 1
	}
	pm.ProbeErrorCount.ErrorCount[s] = TimeValue{Value: val, Time: t}
	pm.ProbeErrorCount.Unlock()
}

func (pm *jsonExport) IncProbeTimeoutCount(s string, t int64) {
	pm.ProbeTimeoutCount.Lock()
	var val float64
	_, ok := pm.ProbeTimeoutCount.TimeoutCount[s]
	if ok {
		val = pm.ProbeTimeoutCount.TimeoutCount[s].Value + 1
	} else {
		val = 1
	}
	pm.ProbeTimeoutCount.TimeoutCount[s] = TimeValue{Value: val, Time: t}
	pm.ProbeTimeoutCount.Unlock()
}

func (pm *jsonExport) SetFieldValues(s string, pd *modules.ProbeData, t int64) {
	pm.ProbeIsUp.Lock()
	pm.ProbeIsUp.Up[s] = TimeValue{Value: *pd.IsUp, Time: t}
	pm.ProbeIsUp.Unlock()

	pm.ProbeLatency.Lock()
	pm.ProbeLatency.Latency[s] = TimeValue{Value: *pd.Latency, Time: t}
	pm.ProbeLatency.Unlock()

	if pd.PayloadSize != nil {
		pm.ProbePayloadSize.Lock()
		pm.ProbePayloadSize.Payload[s] = TimeValue{Value: *pd.PayloadSize, Time: t}
		pm.ProbePayloadSize.Unlock()
	}
}

// SetFieldValuesUnexpected sets values to the fields to -1 to indicate a probe module error/timeout.
func (pm *jsonExport) SetFieldValuesUnexpected(s string, t int64) {
	pm.ProbeIsUp.Lock()
	pm.ProbeIsUp.Up[s] = TimeValue{Value: -1, Time: t}
	pm.ProbeIsUp.Unlock()

	pm.ProbeLatency.Lock()
	pm.ProbeLatency.Latency[s] = TimeValue{Value: -1, Time: t}
	pm.ProbeLatency.Unlock()

	pm.ProbePayloadSize.Lock()
	pm.ProbePayloadSize.Payload[s] = TimeValue{Value: -1, Time: t}
	pm.ProbePayloadSize.Unlock()
}

func jsonHttpHandler(pm *jsonExport) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		dst, err := json.MarshalIndent(pm, "", " ")
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(dst))
	}
	return http.HandlerFunc(fn)
}

func (pm *jsonExport) MetricHttpHandler() http.Handler {
	return jsonHttpHandler(pm)

}

func (pm *jsonExport) MarshalJSON() ([]byte, error) {
	m := make(map[string]interface{})

	pm.ProbeCount.RLock()
	m["probe_count"] = pm.ProbeCount.Count
	pm.ProbeCount.RUnlock()

	pm.ProbeErrorCount.RLock()
	m["probe_error_count"] = pm.ProbeErrorCount.ErrorCount
	pm.ProbeErrorCount.RUnlock()

	pm.ProbeTimeoutCount.RLock()
	m["probe_timeout_count"] = pm.ProbeTimeoutCount.TimeoutCount
	pm.ProbeTimeoutCount.RUnlock()

	pm.ProbeIsUp.RLock()
	m["probe_up"] = pm.ProbeIsUp.Up
	pm.ProbeIsUp.RUnlock()

	pm.ProbeLatency.RLock()
	m["probe_latency"] = pm.ProbeLatency.Latency
	pm.ProbeLatency.RUnlock()

	pm.ProbePayloadSize.RLock()
	m["probe_payload_size"] = pm.ProbePayloadSize.Payload
	pm.ProbePayloadSize.RUnlock()

	return json.Marshal(m)
}

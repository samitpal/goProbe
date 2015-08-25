package metric_export

import (
	"encoding/json"
	"flag"
	"fmt"
	"github.com/samitpal/goProbe/modules"
	"net/http"
	"sync"
)

var (
	jsonMetricsPath = flag.String("json_metrics_path", "/metrics", "Metric exposition path.")
)

type ProbeCount struct {
	sync.RWMutex
	Count map[string]int64 `json:"probe_count"`
}

type ProbeErrorCount struct {
	sync.RWMutex
	ErrorCount map[string]int64 `json:"probe_error_count"`
}

type ProbeTimeoutCount struct {
	sync.RWMutex
	TimeoutCount map[string]int64 `json:"probe_timeout_count"`
}

type ProbeIsUp struct {
	sync.RWMutex
	Up map[string]float64 `json:"probe_is_up"` // value of 1 is a success while 0 is a failure.
}

type ProbeLatency struct {
	sync.RWMutex
	Latency map[string]float64 `json:"probe_latency"`
}

type ProbePayloadSize struct {
	sync.RWMutex
	Payload map[string]float64 `json:"probe_payload_size"`
}

type jsonExport struct {
	ProbeCount
	ProbeErrorCount
	ProbeTimeoutCount
	ProbeIsUp
	ProbeLatency
	ProbePayloadSize
}

func NewJSONExport() *jsonExport {
	return &jsonExport{
		ProbeCount:        ProbeCount{Count: make(map[string]int64)},
		ProbeErrorCount:   ProbeErrorCount{ErrorCount: make(map[string]int64)},
		ProbeTimeoutCount: ProbeTimeoutCount{TimeoutCount: make(map[string]int64)},
		ProbeIsUp:         ProbeIsUp{Up: make(map[string]float64)},
		ProbeLatency:      ProbeLatency{Latency: make(map[string]float64)},
		ProbePayloadSize:  ProbePayloadSize{Payload: make(map[string]float64)},
	}

}

func (pm *jsonExport) Prepare() {
	// Nothing to do.
}

func (pm *jsonExport) IncProbeCount(s string) {
	pm.ProbeCount.Lock()
	pm.ProbeCount.Count[s]++
	pm.ProbeCount.Unlock()
}

func (pm *jsonExport) IncProbeErrorCount(s string) {
	pm.ProbeErrorCount.Lock()
	pm.ProbeErrorCount.ErrorCount[s]++
	pm.ProbeErrorCount.Unlock()
}

func (pm *jsonExport) IncProbeTimeoutCount(s string) {
	pm.ProbeTimeoutCount.Lock()
	pm.ProbeTimeoutCount.TimeoutCount[s]++
	pm.ProbeTimeoutCount.Unlock()
}

func (pm *jsonExport) SetFieldValues(s string, pd *modules.ProbeData) {
	pm.ProbeIsUp.Lock()
	pm.ProbeIsUp.Up[s] = *pd.IsUp
	pm.ProbeIsUp.Unlock()

	pm.ProbeLatency.Lock()
	pm.ProbeLatency.Latency[s] = *pd.Latency
	pm.ProbeLatency.Unlock()

	pm.ProbePayloadSize.Lock()
	pm.ProbePayloadSize.Payload[s] = *pd.PayloadSize
	pm.ProbePayloadSize.Unlock()
}

// SetFieldValuesUnexpected sets values to the fields to -1 to indicate a probe error/timeout.
func (pm *jsonExport) SetFieldValuesUnexpected(s string) {
	pm.ProbeIsUp.Lock()
	pm.ProbeIsUp.Up[s] = -1
	pm.ProbeIsUp.Unlock()

	pm.ProbeLatency.Lock()
	pm.ProbeLatency.Latency[s] = -1
	pm.ProbeLatency.Unlock()

	pm.ProbePayloadSize.Lock()
	pm.ProbePayloadSize.Payload[s] = -1
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

func (pm *jsonExport) RegisterHttpHandler() {
	http.Handle(*jsonMetricsPath, jsonHttpHandler(pm))

}

func (pm *jsonExport) MarshalJSON() ([]byte, error) {
	m := make(map[string]interface{})
	fmt.Println("Inside marshaljson")

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

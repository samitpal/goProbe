package modules

import (
	"encoding/json"
	"fmt"
)

type TestProbe struct {
	ProbeName     *string `json:"probe_name"`
	ProbeInterval *int    `json:"probe_interval"`
	ProbeTimeout  *int    `json:"probe_timeout"`
	ProbeMyConfig *string `json:"probe_my_config"` // this config field is specific to this module.
}

func (t *TestProbe) Prepare() error {
	fmt.Println("Do nothing")
	return nil
}

func (t *TestProbe) Name() *string {
	return t.ProbeName
}
func (t *TestProbe) Run(respCh chan<- *ProbeData, errCh chan<- error) {
	isUp := float64(0)
	latency := float64(40)
	startTime := int64(357683)
	endTime := int64(457683)

	respCh <- &ProbeData{
		IsUp:      &isUp,
		Latency:   &latency,
		StartTime: &startTime,
		EndTime:   &endTime,
	}
	return
}
func (t *TestProbe) TimeoutSecs() *int {
	return t.ProbeTimeout
}

func (t *TestProbe) RunIntervalSecs() *int {
	return t.ProbeInterval
}

func (t *TestProbe) RetConfig() string {
	ret, _ := json.Marshal(t)
	return string(ret)
}

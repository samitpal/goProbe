package modules

import (
	"net/http"
)

type HttpFields struct {
	Headers *http.Header
	Status  *string
}

// Probedata is the struct which holds the probe response info.
type ProbeData struct {
	IsUp        *float64    // Indicates the success/failure of the probe.
	PayloadSize *float64    // Optional. Response payload size.
	Latency     *float64    // Latency in milli seconds.
	StartTime   *int64      // Unix epoch in nano seconds.
	EndTime     *int64      // Unix epoch in nano seconds.
	Http        *HttpFields // Optional, for the http module.
	Payload     *[]byte     // Optional.
}

// Prober is the interface that a probe module needs to implement.
type Prober interface {
	// Prepare is used to set up the probe module. Use it to do custom initialization.
	// Prepare is guranteed to be called before executing the Run() method of the module.
	Prepare() error

	// Run runs the probe. It will be called in a loop. This is the most import method of the interface.
	// Implementation should send the probe response data through the ProbeData channel.
	// In case of any error, the same should be send via the error channel. The ProbeData response channel
	// should not be used in error situations.
	Run(chan<- *ProbeData, chan<- error)

	// Name returns the name of the probe.
	Name() *string

	// TimeoutSecs returns the timeout for a given probe.
	TimeoutSecs() *int

	// RunIntervalSecs returns the frequency of the probe.
	RunIntervalSecs() *int

	// RetConfig returns the config values of the probe module. This will be used in the http ui.
	RetConfig() string
}

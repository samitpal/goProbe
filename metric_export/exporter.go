package metric_export

import "github.com/samitpal/goProbe/modules"

// MetricExporter interface is implemented by an exporter which wants to expose the probe metrics in its own format.
type MetricExporter interface {
	// Prepare should be used for initialization. It is guranteed to be called first, before any another methods.
	Prepare()

	// IncProbeCount increments the probe count of a given probe.
	IncProbeCount(string)

	// IncErrorCount increments the error count of a given probe.
	IncProbeErrorCount(string)

	// IncTimeoutCount increments the timeout count of a given probe.
	IncProbeTimeoutCount(string)

	// SetFieldValues function sets the field values during normal times, e.g set the ‘up’ variable to 1 or 0.
	SetFieldValues(string, *modules.ProbeData)

	// SetFieldValuesUnexpected function sets field values during unexpected situations, e.g probe errors/timeouts. For instance
	// one might want to set the ‘up’ variable for a probe which timed out to -1 instead of a 0 or 1.
	SetFieldValuesUnexpected(string)

	//RegisterHttpHandler registers an http handler to expose the metrics.
	RegisterHttpHandler()
}

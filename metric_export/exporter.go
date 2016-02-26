package metric_export

import (
	"errors"
	"github.com/marpaia/graphite-golang"
	"github.com/samitpal/goProbe/modules"
	"net/http"
)

// MetricExporter interface is implemented by an exporter which wants to expose the probe metrics in its own format.
type MetricExporter interface {
	// Prepare should be used for initialization. It is guranteed to be called first, before any another methods.
	Prepare()

	// IncProbeCount increments the probe count of a given probe. It takes the probe name and epoch time (seconds) as args.
	IncProbeCount(string, int64)

	// IncErrorCount increments the error count of a given probe.  It takes the probe name and epoch time (seconds) as args.
	IncProbeErrorCount(string, int64)

	// IncTimeoutCount increments the timeout count of a given probe.  It takes the probe name and epoch time (seconds) as args.
	IncProbeTimeoutCount(string, int64)

	// SetFieldValues function sets the field values during normal times, e.g set the ‘up’ variable to 1 or 0.
	// It takes the probe name, probe response and epoch time (seconds) as args.
	SetFieldValues(string, *modules.ProbeData, int64)

	// SetFieldValuesUnexpected function sets field values during unexpected situations, e.g probe errors/timeouts. For instance
	// one might want to set the ‘up’ variable for a probe which timed out to -1 instead of a 0 or 1.
	// It takes the probe name and epoch time (seconds) as args.
	SetFieldValuesUnexpected(string, int64)

	// MetricHttpHandler returns the http handler to expose the metrics via a given path (e.g /metrics).
	MetricHttpHandler() http.Handler

	// RetGraphiteMetrics is only implemented by the json exporter. This will be used for pushing metrics
	// It returns the metrics for a given probe in graphite metric format
	RetGraphiteMetrics(string) []graphite.Metric
}

func SetupMetricExporter(s string) (MetricExporter, error) {
	var mExp MetricExporter
	if s == "prometheus" {
		mExp = NewPrometheusExport()
	} else if s == "json" {
		mExp = NewJSONExport()
	} else {
		return nil, errors.New("Unknown metric exporter, %s.")
	}
	mExp.Prepare()
	return mExp, nil
}

package metric_export

import (
	"flag"
	"github.com/marpaia/graphite-golang"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/samitpal/goProbe/modules"
	"net/http"
)

type prometheusExport struct {
	ProbeCount        *prometheus.CounterVec
	ProbeErrorCount   *prometheus.CounterVec
	ProbeTimeoutCount *prometheus.CounterVec
	ProbeIsUp         *prometheus.GaugeVec
	ProbeLatency      *prometheus.GaugeVec
	ProbePayloadSize  *prometheus.GaugeVec
}

var (
	labels                   = []string{"probe_name"}
	prometheusProbeNameSpace = flag.String("prometheus_probe_name_space", "probe", "Prometheus name space of the probes. Valid with prometheus exposition type")
)

func NewPrometheusExport() *prometheusExport {
	return new(prometheusExport)
}

// prometheusExport implements MetricExporter

func (p *prometheusExport) Prepare() {
	p.ProbeIsUp = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Namespace: *prometheusProbeNameSpace,
		Name:      "up",
		Help:      "Indicates success/failure of the probe. Value of 1 is a success while 0 is a failure. Value of -1 could be because of probe timeout/error.",
	}, labels)

	p.ProbeLatency = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Namespace: *prometheusProbeNameSpace,
		Name:      "latency",
		Help:      "The probe latency in milliseconds. Value of -1 could be because of probe timeout/error.",
	}, labels)

	p.ProbePayloadSize = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Namespace: *prometheusProbeNameSpace,
		Name:      "payload_size",
		Help:      "The probe response payload size in bytes. Value of -1 could be because of probe timeout/error.",
	}, labels)

	p.ProbeErrorCount = prometheus.NewCounterVec(prometheus.CounterOpts{
		Namespace: *prometheusProbeNameSpace,
		Name:      "failure_count",
		Help:      "The probe error count.",
	}, labels)

	p.ProbeTimeoutCount = prometheus.NewCounterVec(prometheus.CounterOpts{
		Namespace: *prometheusProbeNameSpace,
		Name:      "timeout_count",
		Help:      "The probe timeout count.",
	}, labels)

	p.ProbeCount = prometheus.NewCounterVec(prometheus.CounterOpts{
		Namespace: *prometheusProbeNameSpace,
		Name:      "count",
		Help:      "Total Probe count.",
	}, labels)

	prometheus.MustRegister(p.ProbeCount)
	prometheus.MustRegister(p.ProbeErrorCount)
	prometheus.MustRegister(p.ProbeTimeoutCount)
	prometheus.MustRegister(p.ProbeLatency)
	prometheus.MustRegister(p.ProbeIsUp)
	prometheus.MustRegister(p.ProbePayloadSize)
}

// IncProbeCount increments the probe count of a given probe.
func (p *prometheusExport) IncProbeCount(probeName string, t int64) {
	p.ProbeCount.WithLabelValues(probeName).Inc()
}

// IncErrorCount increments the error count of a given probe.
func (p *prometheusExport) IncProbeErrorCount(probeName string, t int64) {
	p.ProbeErrorCount.WithLabelValues(probeName).Inc()
}

// IncTimeoutCount increments the timeout count of a given probe.
func (p *prometheusExport) IncProbeTimeoutCount(probeName string, t int64) {
	p.ProbeTimeoutCount.WithLabelValues(probeName).Inc()
}

// SetFieldValues function sets the field values during normal times, e.g set the ‘up’ variable to 1 or 0.
func (p *prometheusExport) SetFieldValues(probeName string, pd *modules.ProbeData, t int64) {
	p.ProbeIsUp.WithLabelValues(probeName).Set(*pd.IsUp)
	p.ProbeLatency.WithLabelValues(probeName).Set(*pd.Latency)
	if pd.PayloadSize != nil {
		p.ProbePayloadSize.WithLabelValues(probeName).Set(*pd.PayloadSize)
	}
}

// SetFieldValuesUnexpected function sets field values during unexpected situations, e.g probe errors/timeouts. For instance
// you might want to set the ‘up’ variable for a probe which timed out to -1 instead of a 0 or 1.
func (p *prometheusExport) SetFieldValuesUnexpected(probeName string, t int64) {
	p.ProbeIsUp.WithLabelValues(probeName).Set(-1)
	p.ProbeLatency.WithLabelValues(probeName).Set(-1)
	p.ProbePayloadSize.WithLabelValues(probeName).Set(-1)
}

//MetricHttpHandler registers a http handler to expose the metrics
func (p prometheusExport) MetricHttpHandler() http.Handler {
	return prometheus.Handler()
}

// This is not used, but needs to be defined to satisfy the interface.
func (p prometheusExport) RetGraphiteMetrics(pn string) []graphite.Metric {
	return []graphite.Metric{}
}

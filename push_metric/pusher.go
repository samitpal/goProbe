package push_metric

import (
	"errors"
	"github.com/samitpal/goProbe/metric_export"
	"github.com/samitpal/goProbe/push_metric/provider"
	"os"
)

// Pusher is the interface that needs needs to implement for pushing metric to (e.g graphite. influxdb).
type Pusher interface {
	Setup()
	PushMetric(metric_export.MetricExporter, string)
}

func SetupProviders() (Pusher, error) {
	if os.Getenv("GOPROBE_PUSH_TO") == "graphite" {
		graphite_host := "localhost"
		if os.Getenv("GORPOBE_GRAPHITE_HOST") != "" {
			graphite_host = os.Getenv("GORPOBE_GRAPHITE_HOST")
		}
		graphite_port := 2003
		if os.Getenv("GORPOBE_GRAPHITE_PORT") != "" {
			graphite_host = os.Getenv("GORPOBE_GRAPHITE_PORT")
		}
		return provider.NewGraphitePusher(graphite_host, graphite_port)
	}
	return nil, errors.New("No push provider found")
}

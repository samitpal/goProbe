package provider

import (
	"github.com/golang/glog"
	"github.com/marpaia/graphite-golang"
	"github.com/samitpal/goProbe/metric_export"
)

type graphitePush struct {
	g *graphite.Graphite
}

func NewGraphitePusher(host string, port int) (*graphitePush, error) {
	graphite, err := graphite.NewGraphiteWithMetricPrefix(host, port, "goProbe")
	if err != nil {
		return nil, err
	}
	return &graphitePush{graphite}, nil
}

// Currently not used. Hence doing nothing
func (g *graphitePush) Setup() {

}

func (g *graphitePush) PushMetric(mExp metric_export.MetricExporter, pn string) {
	metrics := mExp.RetGraphiteMetrics(pn)
	err := g.g.SendMetrics(metrics)
	if err != nil {
		glog.Infof("Error pushing metric", err)
	}
}

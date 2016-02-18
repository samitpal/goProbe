package main

import (
	"github.com/hashicorp/consul/api"
	"github.com/samitpal/goProbe/metric_export"
	"github.com/samitpal/goProbe/misc"
	"github.com/samitpal/goProbe/modules"
	"os"
)

func getConsulClient() (*api.Client, error) {
	var consulHost, consulPort string
	if os.Getenv("GOPROBE_CONSUL_HOST") != "" {
		consulHost = os.Getenv("GOPROBE_CONSUL_HOST")
	} else {
		consulHost = "localhost"
	}
	if os.Getenv("GOPROBE_CONSUL_PORT") != "" {
		consulPort = os.Getenv("GOPROBE_CONSUL_PORT")
	} else {
		consulPort = "8500"
	}
	config := api.DefaultConfig()
	config.Address = consulHost + ":" + consulPort
	client, err := api.NewClient(config)
	return client, err
}

type DoJob struct {
	probes []modules.Prober
	mExp   metric_export.MetricExporter
	ps     *misc.ProbesStatus
}

func NewDoJob(probes []modules.Prober, mExp metric_export.MetricExporter, ps *misc.ProbesStatus) *DoJob {
	return &DoJob{probes, mExp, ps}
}

func (j DoJob) DoJobFunc(stopCh chan bool, doneCh chan bool) {
	// we do not use doneCh since this is a continuously function
	runProbes(j.probes, j.mExp, j.ps, stopCh)
}

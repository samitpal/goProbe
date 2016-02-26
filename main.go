package main

import (
	"flag"
	"github.com/golang/glog"
	"github.com/gorilla/handlers"
	leader_election "github.com/samitpal/consul-client-master-election/election_api"
	"github.com/samitpal/goProbe/conf"
	"github.com/samitpal/goProbe/log"
	"github.com/samitpal/goProbe/metric_export"
	"github.com/samitpal/goProbe/misc"
	"github.com/samitpal/goProbe/modules"
	"github.com/samitpal/goProbe/push_metric"
	"io/ioutil"
	"math/rand"
	"net/http"
	_ "net/http/pprof"
	"os"
	"time"
)

var (
	listenAddress     = flag.String("listen-address", ":8080", "Address to listen on for web interface.")
	configFlag        = flag.String("config", "./probe_config.json", "Path to the probe json config.")
	probeSpaceOutTime = flag.Int("probe_space_out_time", 15, "Max sleep time between probes to allow spacing out of the probes at startup.")
	expositionType    = flag.String("exposition_type", "json", "Metric exposition format.")
	dryRun            = flag.Bool("dry_run", false, "Dry run mode where it does everything except running the probes.")
	metricsPath       = flag.String("metric_path", "/metrics", "Metric exposition path.")
	webLogDir         = flag.String("weblog_dir", "", "Directory path of the web log.")
	haMode            = flag.Bool("ha_mode", false, "Whether to use consul for High Availabity mode.")
	pushMetric        = flag.Bool("push_metric", false, "Whether to push metric to a given provier. If set, one needs to set the GOPROBE_PUSH_TO env variable")
)

func checkFlags() {
	if *pushMetric && (*expositionType == "prometheus") {
		glog.Exitln("Exposition type: prometheus is not compatible with push_metric flag")
	}
}

// runProbes actually runs the probes. This is the core.
func runProbes(pusher push_metric.Pusher, probes []modules.Prober, mExp metric_export.MetricExporter, ps *misc.ProbesStatus, stopCh chan bool) {
	for _, p := range probes {
		// Add some randomness to space out the probes a bit at start up.
		r := rand.New(rand.NewSource(time.Now().UnixNano()))
		time.Sleep(time.Duration(r.Intn(*probeSpaceOutTime)) * time.Second)
		if *pushMetric {
			pusher.Setup()
		}
		go func(p modules.Prober) {
			for {
				// Buffered channel so that the read happens even if there is nothing to receive it. Needed to
				// handle the timeout scenario as well as the situaion when the go routine has to return on stop
				// signal.
				respCh := make(chan *modules.ProbeData, 1)
				errCh := make(chan error, 1)

				pn := *p.Name()
				to := *p.TimeoutSecs()
				timer := time.NewTimer(time.Duration(*p.RunIntervalSecs()) * time.Second)

				glog.Infof("Launching new probe:%s", pn)
				startTime := time.Now().UnixNano()
				startTimeSecs := startTime / 1000000000 // used to expose time field in json metric expostion.
				go p.Run(respCh, errCh)

				select {
				case msg := <-respCh:
					err := misc.CheckProbeData(msg)
					mExp.IncProbeCount(pn, startTimeSecs)
					if err != nil {
						glog.Errorf("Error: %v", err)
						mExp.IncProbeErrorCount(pn, startTimeSecs)
						mExp.SetFieldValuesUnexpected(pn, startTimeSecs)
						ps.WriteProbeErrorStatus(pn, startTime, time.Now().UnixNano())
					} else {
						mExp.SetFieldValues(pn, msg, startTimeSecs)
						ps.WriteProbeStatus(pn, msg, startTime, time.Now().UnixNano())
					}
				case err_msg := <-errCh:
					glog.Errorf("Probe %s error'ed out: %v", pn, err_msg)
					mExp.IncProbeCount(pn, startTimeSecs)
					mExp.IncProbeErrorCount(pn, startTimeSecs)
					mExp.SetFieldValuesUnexpected(pn, startTimeSecs)
					ps.WriteProbeErrorStatus(pn, startTime, time.Now().UnixNano())
				case <-time.After(time.Duration(to) * time.Second):
					glog.Errorf("Timed out probe:%v ", pn)
					mExp.IncProbeCount(pn, startTimeSecs)
					mExp.IncProbeTimeoutCount(pn, startTimeSecs)
					mExp.SetFieldValuesUnexpected(pn, startTimeSecs)
					ps.WriteProbeTimeoutStatus(pn, startTime, time.Now().UnixNano())
				case <-stopCh:
					glog.Infof("Goroutine probe named: %s recieved stop signal. Returning.", pn)
					return
				}
				if *pushMetric {
					go pusher.PushMetric(mExp, pn)
				}
				<-timer.C
			}
		}(p)
	}
}

func main() {

	flag.Parse()
	checkFlags()
	var pusher push_metric.Pusher
	var err error
	if *pushMetric {
		pusher, err = push_metric.SetupProviders()
		if err != nil {
			glog.Exitf("Problem while setting up push provider: %v", err)
		}
	}
	config, err := ioutil.ReadFile(*configFlag)
	if err != nil {
		glog.Exitf("Error reading probe config file: %v", err)
	}
	probes, err := conf.SetupConfig(config)
	if err != nil {
		glog.Exitf("Error in probe config setup, exiting: %v", err)
	}
	err = misc.CheckProbeConfig(probes)
	if err != nil {
		glog.Exitf("Error in probe config, exiting: %v", err)
	}

	probeNames := conf.GetProbeNames(probes)
	mExp, err := metric_export.SetupMetricExporter(*expositionType)
	if err != nil {
		glog.Exitf("Error : %v", err)
	}

	var fh *os.File
	if *webLogDir != "" {
		fh, err = log.SetupWebLog(*webLogDir, time.Now())
		if err != nil {
			glog.Exitf("Failed to set up logging", err)
		}
	} else {
		fh = os.Stdout // logs web accesses to stdout. May not be thread safe.
	}

	ps := misc.NewProbesStatus(probeNames)
	http.Handle("/", handlers.CombinedLoggingHandler(fh, http.HandlerFunc(misc.HandleHomePage)))
	http.Handle("/status", handlers.CombinedLoggingHandler(fh, misc.HandleStatus(ps)))
	http.Handle("/config", handlers.CombinedLoggingHandler(fh, http.HandlerFunc(misc.HandleConfig(config))))
	http.Handle(*metricsPath, handlers.CombinedLoggingHandler(fh, mExp.MetricHttpHandler()))

	glog.Info("Starting goProbe server.")
	glog.Infof("Will expose metrics in %s format via %s http path.", *expositionType, *metricsPath)
	glog.Infof("/config shows current config, /status shows current probe status.")

	if !*dryRun {
		// Start probing.
		stopCh := make(chan bool)
		if *haMode {
			glog.Info("Running in HA mode..")
			client, err := getConsulClient()
			if err != nil {
				glog.Fatalf("Fatal error: %v", err)
			}
			job := NewDoJob(pusher, probes, mExp, ps)
			go leader_election.MaybeAcquireLeadership(client, "goProbe/leader", 20, 30, "goProbe", false, job)
		} else {
			go runProbes(pusher, probes, mExp, ps, stopCh)
		}
		if err = http.ListenAndServe(*listenAddress, nil); err != nil {
			panic(err)
		}
	} else {
		glog.Info("Dry run mode.")
	}
}

package main

import (
	"errors"
	"flag"
	"fmt"
	"github.com/golang/glog"
	"github.com/gorilla/handlers"
	"github.com/samitpal/goProbe/metric_export"
	"github.com/samitpal/goProbe/misc"
	"github.com/samitpal/goProbe/modules"
	"html/template"
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
	templates         = template.Must(template.ParseGlob(os.Getenv("GOPROBE_TMPL")))
)

func setupMetricExporter(s string) (metric_export.MetricExporter, error) {
	var mExp metric_export.MetricExporter
	if s == "prometheus" {
		mExp = metric_export.NewPrometheusExport()
	} else if s == "json" {
		mExp = metric_export.NewJSONExport()
	} else {
		return nil, errors.New("Unknown metric exporter, %s.")
	}
	mExp.Prepare()
	return mExp, nil
}

// runProbes actually runs the probes. This is the core.
func runProbes(probes []modules.Prober, mExp metric_export.MetricExporter, ps *misc.ProbesStatus) {
	for _, p := range probes {
		// Add some randomness to space out the probes a bit at start up.
		r := rand.New(rand.NewSource(time.Now().UnixNano()))
		time.Sleep(time.Duration(r.Intn(*probeSpaceOutTime)) * time.Second)
		go func(p modules.Prober) {
			for {
				// Buffered channel so that the read happens even if there is nothing to receive it. Needed to
				// handle the timeout scenario.
				respCh := make(chan *modules.ProbeData, 1)
				errCh := make(chan error, 1)

				pn := *p.Name()
				to := *p.TimeoutSecs()
				timer := time.NewTimer(time.Duration(*p.RunIntervalSecs()) * time.Second)

				glog.Infof("Launching new probe:%s", pn)
				go p.Run(respCh, errCh)

				select {
				case msg := <-respCh:
					err := checkProbeData(msg)
					mExp.IncProbeCount(pn)
					if err != nil {
						glog.Errorf("Error: %v", err)
						mExp.IncProbeErrorCount(pn)
						mExp.SetFieldValuesUnexpected(pn)
						ps.WriteProbeErrorStatus(pn)
					} else {
						mExp.SetFieldValues(pn, msg)
						ps.WriteProbeStatus(pn, msg)
					}
				case err_msg := <-errCh:
					glog.Errorf("Probe %s error'ed out: %v", pn, err_msg)
					mExp.IncProbeCount(pn)
					mExp.IncProbeErrorCount(pn)
					mExp.SetFieldValuesUnexpected(pn)
					ps.WriteProbeErrorStatus(pn)
				case <-time.After(time.Duration(to) * time.Second):
					glog.Errorf("Timed out probe:%v ", pn)
					mExp.IncProbeCount(pn)
					mExp.IncProbeTimeoutCount(pn)
					mExp.SetFieldValuesUnexpected(pn)
					ps.WriteProbeTimeoutStatus(pn)
				}
				<-timer.C
			}
		}(p)
	}
}

// checkProbeConfig function does sanity checks on the probe definition.
func checkProbeConfig(probes []modules.Prober) error {
	if len(probes) == 0 {
		return errors.New("No probe modules defined")
	}
	for _, p := range probes {
		if *p.TimeoutSecs() > *p.RunIntervalSecs() {
			return fmt.Errorf("Timeout can not be more than the Interval %v", p.Name())
		}
	}
	return nil
}

// checkProbedata function verifies the correctness of the probe response.
func checkProbeData(pd *modules.ProbeData) error {
	if pd.IsUp == nil {
		return errors.New("Mandatory field 'IsUp' is missing in probe response")
	}
	if pd.Latency == nil {
		return errors.New("Mandatory field 'Latency' is missing in probe response")
	}
	if pd.StartTime == nil {
		return errors.New("Mandatory field 'StartTime' is missing in probe response")
	}
	if pd.EndTime == nil {
		return errors.New("Mandatory field 'EndTime' is missing in probe response")
	}
	return nil
}

func handleHomePage(w http.ResponseWriter, r *http.Request) {
	err := templates.ExecuteTemplate(w, "indexPage", nil)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

// TODO: Return pure json instead of html.
func handleConfig(w http.ResponseWriter, r *http.Request) {
	config, err := ioutil.ReadFile(*configFlag)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	probes, err := setupConfig(config)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	err = templates.ExecuteTemplate(w, "configPage", probes)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func handleStatus(ps *misc.ProbesStatus) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		showParams := r.URL.Query().Get("showparams")
		if showParams == "single_probe" {
			ps.Tmpl.ProbeSingle = r.URL.Query().Get("probe_name")
		}
		ps.Tmpl.ShowParams = showParams
		err := templates.ExecuteTemplate(w, "statusPage", ps)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}
}

func main() {

	flag.Parse()
	config, err := ioutil.ReadFile(*configFlag)
	if err != nil {
		glog.Exitf("Error reading config file: %v", err)
	}
	probes, err := setupConfig(config)
	if err != nil {
		glog.Exitf("Error in config setup, exiting: %v", err)
	}
	err = checkProbeConfig(probes)
	if err != nil {
		glog.Exitf("Error in probe config, exiting: %v", err)
	}

	probeNames := getProbeNames(probes)
	mExp, err := setupMetricExporter(*expositionType)
	if err != nil {
		glog.Exitf("Error : %v", err)
	}

	var fh *os.File
	if *webLogDir != "" {
		fh, err = setupWebLog(*webLogDir, time.Now())
		if err != nil {
			glog.Exitf("Failed to set up logging", err)
		}
	} else {
		fh = os.Stdout // logs web accesses to stdout. May not be thread safe.
	}

	ps := misc.NewProbesStatus(probeNames)
	http.Handle("/", handlers.CombinedLoggingHandler(fh, http.HandlerFunc(handleHomePage)))
	http.Handle("/status", handlers.CombinedLoggingHandler(fh, handleStatus(ps)))
	http.Handle("/config", handlers.CombinedLoggingHandler(fh, http.HandlerFunc(handleConfig)))
	http.Handle(*metricsPath, handlers.CombinedLoggingHandler(fh, mExp.MetricHttpHandler()))

	glog.Info("Starting goProbe server.")
	glog.Infof("Will expose metrics in %s format via %s http path.", *expositionType, *metricsPath)
	glog.Infof("/config shows current config, /status shows current probe status.")

	if !*dryRun {
		// Start probing.
		go runProbes(probes, mExp, ps)
		if err = http.ListenAndServe(*listenAddress, nil); err != nil {
			panic(err)
		}
	} else {
		glog.Info("Dry run mode.")
	}
}

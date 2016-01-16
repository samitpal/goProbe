package misc

import (
	"errors"
	"fmt"
	"github.com/samitpal/goProbe/conf"
	"github.com/samitpal/goProbe/modules"
	"html/template"
	"net/http"
	"os"
)

var templates = template.Must(template.ParseGlob(os.Getenv("GOPROBE_TMPL")))

// CheckProbeConfig function does sanity checks on the probe definition.
func CheckProbeConfig(probes []modules.Prober) error {
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

// CheckProbedata function verifies the correctness of the probe response.
func CheckProbeData(pd *modules.ProbeData) error {
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

func HandleHomePage(w http.ResponseWriter, r *http.Request) {
	err := templates.ExecuteTemplate(w, "indexPage", nil)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

// TODO: Return pure json instead of html.
func HandleConfig(config []byte) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		probes, err := conf.SetupConfig(config)
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
}

func HandleStatus(ps *ProbesStatus) http.HandlerFunc {
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

package misc

import (
	"github.com/samitpal/goProbe/modules"
	"html/template"
	"net/http"
	"strings"
	"sync"
	"time"
)

type ProbeStatus struct {
	ProbeResp    *modules.ProbeData
	ProbeError   bool
	ProbeTimeout bool
}

type TemplateParams struct {
	ShowParams  string // value can be all, single_probe
	ProbeSingle string // name of the probe to show in the templates when ShowParams==single_probe
}

type ProbesStatus struct {
	Tmpl           TemplateParams
	Probes         []string
	ProbeStatusMap map[string]*ProbeStatus
	lock           sync.RWMutex
}

func NewProbesStatus(p []string) ProbesStatus {
	return ProbesStatus{
		Tmpl:           TemplateParams{ShowParams: "all", ProbeSingle: ""},
		Probes:         p,
		ProbeStatusMap: make(map[string]*ProbeStatus),
	}
}

func (ps ProbesStatus) WriteProbeStatus(pn string, pd *modules.ProbeData) {
	ps.lock.Lock()
	defer ps.lock.Unlock()

	ps.ProbeStatusMap[pn] = &ProbeStatus{
		ProbeResp:    pd,
		ProbeError:   false,
		ProbeTimeout: false,
	}

}

func (ps ProbesStatus) WriteProbeErrorStatus(pn string) {
	ps.lock.Lock()
	defer ps.lock.Unlock()

	ps.ProbeStatusMap[pn] = &ProbeStatus{
		ProbeResp:    nil,
		ProbeError:   true,
		ProbeTimeout: false,
	}

}

func (ps ProbesStatus) WriteProbeTimeoutStatus(pn string) {
	ps.lock.Lock()
	defer ps.lock.Unlock()

	ps.ProbeStatusMap[pn] = &ProbeStatus{
		ProbeResp:    nil,
		ProbeError:   false,
		ProbeTimeout: true,
	}
}

func (ps ProbesStatus) ReadProbeStatus(pn string) *ProbeStatus {
	ps.lock.RLock()
	defer ps.lock.RUnlock()

	val, ok := ps.ProbeStatusMap[pn]
	if ok {
		return val
	} else {
		return nil
	}
}

func (ps ProbesStatus) ConvertToInt(i *float64) int {
	return int(*i)
}

func (ps ProbesStatus) FormattedTime(t *int64) time.Time {
	return time.Unix(0, *t)
}

func (ps ProbesStatus) FormattedHttpHeaders(h http.Header) template.HTML {

	var headers string
	for key, value := range h {
		headers = headers + key + ": " + strings.Join(value, " ") + "<br />"
	}
	return template.HTML(headers) // we are assuming the headers are html safe.
}

func (ps ProbesStatus) FormattedHttpBody(body *[]byte) string {
	return string(*body)
}

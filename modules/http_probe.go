package modules

import (
	"encoding/json"
	"errors"
	"github.com/golang/glog"
	"io/ioutil"
	"net/http"
	"regexp"
	"time"
)

type httpProbe struct {
	ProbeName        *string `json:"probe_name"`
	ProbeURL         *string `json:"probe_url"`
	ProbeHttpMethod  *string `json:"probe_http_method"`
	ProbeAction      *string `json:"probe_action"`
	ProbeMatchString *string `json:"probe_match_string"` // a regulat expression.
	ProbeInterval    *int    `json:"probe_interval"`
	ProbeTimeout     *int    `json:"probe_timeout"`
}

func NewHttpProbe() *httpProbe {
	return new(httpProbe)
}

func (p httpProbe) checkConfig() error {
	if p.ProbeName == nil {
		return errors.New("Required field ProbeName is not set")
	}
	if p.ProbeURL == nil {
		return errors.New("Required field ProbeURL is not set")
	}

	if p.ProbeAction != nil && p.ProbeMatchString == nil {
		return errors.New("ProbeMatchString is required")
	}
	return nil
}

func (p *httpProbe) setDefaults() {
	if p.ProbeHttpMethod == nil {
		str := "GET"
		p.ProbeHttpMethod = &str
	}
	if p.ProbeAction == nil {
		str := "check_ret_200"
		p.ProbeAction = &str
	}
	if p.ProbeTimeout == nil {
		i := 40
		p.ProbeTimeout = &i
	}
	if p.ProbeInterval == nil {
		i := 60
		p.ProbeInterval = &i
	}

}

// httpProbe implements the Prober interface.

func (p *httpProbe) Prepare() error {
	err := p.checkConfig()
	if err != nil {
		return err
	}
	p.setDefaults()
	return nil
}

func (p httpProbe) Run(respCh chan<- *ProbeData, errCh chan<- error) {
	// Run the http probe
	startTime := time.Now().UnixNano()
	client := &http.Client{Timeout: time.Duration(*p.ProbeTimeout) * time.Second}

	req, err := http.NewRequest(*p.ProbeHttpMethod, *p.ProbeURL, nil)
	if err != nil {
		glog.Errorf("Error: %v", err)
		errCh <- err
		return
	}

	resp, err := client.Do(req)
	if err != nil {
		glog.Errorf("Error: %v", err)
		errCh <- err
		return
	}

	respPayloadSize := float64(resp.ContentLength)
	respHeader := resp.Header
	respPayload, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		glog.Errorf("Error: %v", err)
		errCh <- err
		return
	}

	defer resp.Body.Close()
	var isUp float64
	if *p.ProbeAction == "check_ret_200" {
		if resp.StatusCode == 200 {
			isUp = 1
		} else {
			isUp = 0
		}
	} else if *p.ProbeAction == "check_match_payload" {
		// Match the response body with the given regexp.
		r := regexp.MustCompile(*p.ProbeMatchString)
		if r.Match(respPayload) {
			isUp = 1
		} else {
			isUp = 0
		}
	}

	endTime := time.Now().UnixNano()
	latency := (float64(endTime - startTime)) / 1000000

	respCh <- &ProbeData{
		IsUp:        &isUp,
		PayloadSize: &respPayloadSize,
		Latency:     &latency,
		StartTime:   &startTime,
		EndTime:     &endTime,
		Headers:     &respHeader,
		Payload:     &respPayload,
	}
	return
}

func (p httpProbe) Name() *string {
	return p.ProbeName
}

func (p httpProbe) RunIntervalSecs() *int {
	return p.ProbeInterval
}

func (p httpProbe) TimeoutSecs() *int {
	return p.ProbeTimeout
}

func (p *httpProbe) RetConfig() string {
	ret, _ := json.MarshalIndent(p, "", " ")
	return string(ret)
}

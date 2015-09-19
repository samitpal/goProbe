// Package to ping a port using tcp/udp protocol. Note that this does not send/receive any data.
// It just checks if the given host:port is reachable using the specified network.
package ping_port

import (
	"encoding/json"
	"errors"
	"github.com/golang/glog"
	"github.com/samitpal/goProbe/modules"
	"net"
	"strconv"
	"time"
)

type pingPortProbe struct {
	ProbeName     *string `json:"probe_name"`
	ProbeInterval *int    `json:"probe_interval"`
	ProbeTimeout  *int    `json:"probe_timeout"`
	ProbeHostName *string `json:"probe_host_name"`
	ProbeHostPort *int    `json:"probe_host_port"`
	ProbeNetwork  *string `json:"probe_network"` //tcp or udp.
}

func NewPingPortProbe() *pingPortProbe {
	return new(pingPortProbe)
}

func (p pingPortProbe) checkConfig() error {
	if p.ProbeName == nil {
		return errors.New("Required field Probe Name is not set")
	}
	if p.ProbeHostName == nil {
		return errors.New("Required field ProbeH ostName is not set")
	}
	if p.ProbeHostPort == nil {
		return errors.New("Required field Probe Host Port is not set")
	}
	if p.ProbeNetwork != nil {
		if *p.ProbeNetwork == "tcp" || *p.ProbeNetwork == "udp" {
			// we are good.
		} else {
			return errors.New("Probe method can only be either of 'tcp' or 'udp'")

		}
	}
	return nil
}

func (p *pingPortProbe) setDefaults() {
	if p.ProbeNetwork == nil {
		network := "tcp"
		p.ProbeNetwork = &network
	}
	if p.ProbeTimeout == nil {
		timeout := 10
		p.ProbeTimeout = &timeout // since we don't send/receive any data, setting it to low value.
	}
	if p.ProbeInterval == nil {
		interval := 60
		p.ProbeInterval = &interval
	}
}

func (p *pingPortProbe) Prepare() error {
	if err := p.checkConfig(); err != nil {
		glog.Errorf("Error in config %v", err)
		return err
	}
	p.setDefaults()
	return nil
}

func (p *pingPortProbe) Name() *string {
	return p.ProbeName
}

func (p *pingPortProbe) Run(respCh chan<- *modules.ProbeData, errCh chan<- error) {
	startTime := time.Now().UnixNano()
	var isUp float64
	// we set timeout less by 1 sec since we want a slightly higher timeout for the caller (core).
	timeout := *p.ProbeTimeout - 1
	conn, err := net.DialTimeout(*p.ProbeNetwork, *p.ProbeHostName+":"+strconv.Itoa(*p.ProbeHostPort), time.Duration(timeout)*time.Second)
	
	if conn != nil {
		defer conn.Close()
	}

	if err != nil {
		isUp = 0
	} else {
		isUp = 1
	}
	endTime := time.Now().UnixNano()
	latency := (float64(endTime - startTime)) / 1000000

	respCh <- &modules.ProbeData{
		IsUp:        &isUp,
		Latency:     &latency,
		StartTime:   &startTime,
		EndTime:     &endTime,
		PayloadSize: nil,
	}
	return
}
func (p *pingPortProbe) TimeoutSecs() *int {
	return p.ProbeTimeout
}

func (p *pingPortProbe) RunIntervalSecs() *int {
	return p.ProbeInterval
}

func (p *pingPortProbe) RetConfig() string {
	ret, _ := json.MarshalIndent(p, "", " ")
	return string(ret)
}

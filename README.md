[![Build Status](https://travis-ci.org/samitpal/goProbe.svg?branch=master)](https://travis-ci.org/samipal/goProbe)

[google group](https://groups.google.com/forum/#!forum/goprobe)

Summary
------------------
goProbe is a probe service written in Go programming language. It has three parts to it, a core, probe modules and a metric exposition component. 

It currently supports probe metric exposition in json (default) as well as [prometheus](http://prometheus.io) compatible format. 

Probes are modules and goProbe can potentially be confgured with arbitrary number of modules. Currently it has a built in http probe module and a ping_port module. The latter helps probe a host:port using tcp/udp.

goProbe takes a json file as a config input which typically supplies the probe name (each should be unique), run intervals, timeouts etc. Each module can have its own json config fields. Below is an example config snippet of the built-in http, ping_port probe modules. It configures two http probes and 1 ping_port probe.

    [
        {
        "probe_type": "http",
        "probe_config": {
            "probe_name": "example_explicit_settings",
            "probe_url": "http://example.com",
            "probe_timeout": 20,
            "probe_interval": 30
            }
        },
        {
        "probe_type": "http",
        "probe_config": {
            "probe_name": "example_default_settings",
            "probe_url": "http://example.com"
            }
        },
        {
        "probe_type": "ping_port",
        "probe_config": {
            "probe_name": "example_ping_port_80",
            "probe_host_name": "example.com",
            "probe_host_port": 80
            }
        }
    ]  


Installation
-------------------
##### Docker image is available at
 https://hub.docker.com/r/samitpal/goprobe/

##### Precompiled binaries are available from the releases link.
https://github.com/samitpal/goProbe/releases/

##### To build from source follow the steps below: 

* Install mercuruial. On ubuntu,
$ sudo apt-get install mercurial

* $ go get -u github.com/samitpal/goProbe

* $ cd $GOPATH/src/github.com/samitpal/goProbe

* $ go install

Running the binary
-------------------

First set the environmental variable which points to the html templates, e.g

$ export GOPROBE_TMPL="templates/*" 

$ $GOPATH/bin/goProbe -config <*path to config file*>

By default goProbe displays the probe metrics via the **/metrics** http handler in json format. It also displays the current configs via its **/config** http handler. The /status handler displays the lastest probe status.

To expose the metrics in prometheus format, run it as follows,

$ $GOPATH/bin/goProbe -config <*path to config file*> -exposition_type prometheus

The json format is time series friendly in that the metrics contain a time field. It just needs a simple script to parse the data from the /metrics end point and push that to a time series database like graphite, influxdb etc. Example push scripts are available at https://github.com/samitpal/goProbe-metric-push.

Http probe json configs
-------------------

### Mandatory fields 
* probe_name: The name of the probe. This should be unique globally.
* probe_url: The complete url.

### Other fields

* probe\_http_method : The http method. Default value is http GET.
* probe\_action : This configures the mechanism used to determine the  success/failure of a probe. Currently it can be one of the following strings
    * check\_ret\_200 : It checks if the response status code is 200. This is the default.
    * check\_match\_string : If this is set then probe\_match\_string needs to be set as well. Essentially the module will match the http body with the probe\_match\_string string.
    * check\_sslcert\_expiry : It checks if the ssl cert is expiring within probe\_sslcert\_expire\_in\_days days.
* probe\_match\_string : This needs to be set to a regexp if probe_action is set to "check\_match\_string". 
* probe\_sslcert\_expire\_in_days: Goes with "check\_sslcert\_expiry" action. Default value is 30 days.
* probe\_interval : The frequency (in seconds) with which to run the probe. Default value is 60.
* probe\_timeout : Time out in seconds for a given probe. Default value is 40. This value needs to be less than the probe\_interval.
* probe\_http\_headers: This config object is used to set http request headers. Following is an example usage.

        "probe_http_headers": {
                "host": "test",
                "user_agent": "test-agent"
          }

Ping port probe json configs
-------------------

### Mandatory fields 
* probe\_name: The name of the probe. This should be unique globally.
* probe\_host_name: The target host.
* probe\_host\_port: The port of the target host.

### Other fields

* probe\_network: The network protocol to use. It can be either tcp (default) or udp.

Developing a probe module
------------------
There is a sample module called sample_probe.go which can give you a start. Essentially all you need to do is implement the **Prober** interface defined in **prober.go**. In addition to the module source it also needs a minor change in **config.go** to register the new module type.